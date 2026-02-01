package source

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HTTPDownloader struct {
	Timeout time.Duration
	Retries int
}

func unzipSingleFile(zipPath, destFile string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		name := strings.ToLower(filepath.Base(f.Name))

		if strings.HasSuffix(name, ".dbf") {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			log.Println("Extracted DBF directly to:", destFile)

			out, err := os.Create(destFile)
			if err != nil {
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, rc)
			return err
		}
	}

	return fmt.Errorf("no .dbf found in zip")
}

func (d *HTTPDownloader) Download(url, dest string, fileType string) error {
	transport := &http.Transport{
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   d.Timeout,
		Transport: transport,
	}

	var lastErr error

	for attempt := 1; attempt <= d.Retries; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			continue
		}

		log.Println("HTTP status:", resp.Status)
		log.Println("Content-Length:", resp.ContentLength)
		log.Println("Content-Encoding:", resp.Header.Get("Content-Encoding"))
		log.Println("Content-Type:", resp.Header.Get("Content-Type"))

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("http status %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		tmp := dest + ".tmp"
		out, err := os.Create(tmp)
		if err != nil {
			resp.Body.Close()
			return err
		}

		written, err := io.Copy(out, resp.Body)
		out.Close()
		resp.Body.Close()

		if err != nil {
			lastErr = err
			os.Remove(tmp)
			continue
		}

		// 🔍 Проверка размера
		if resp.ContentLength > 0 && written != resp.ContentLength {
			lastErr = fmt.Errorf(
				"downloaded size mismatch: got %d, expected %d",
				written, resp.ContentLength,
			)
			os.Remove(tmp)
			continue
		}

		log.Printf("Downloader: type=%s url=%s dest=%s", fileType, url, dest)

		switch fileType {
		case "dbf_zip":
			log.Println("DBF ZIP detected, extracting")
			err := unzipSingleFile(tmp, dest)
			if err != nil {
				return err
			}
			os.Remove(tmp)
			return nil

		case "xlsx":
			// просто сохраняем
			return os.Rename(tmp, dest)
		}

		// обычный файл → атомарная замена
		if err := os.Rename(tmp, dest); err != nil {
			return err
		}

		log.Println("Downloaded OK:", dest, "bytes:", written)
		return nil
	}

	return lastErr
}
