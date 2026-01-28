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

func unzipSingleFile(zipPath, destDir string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		// берём первый DBF
		if strings.HasSuffix(strings.ToLower(f.Name), ".dbf") {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			outPath := filepath.Join(destDir, filepath.Base(f.Name))
			out, err := os.Create(outPath)
			if err != nil {
				return "", err
			}
			defer out.Close()

			if _, err := io.Copy(out, rc); err != nil {
				return "", err
			}

			return outPath, nil
		}
	}

	return "", fmt.Errorf("no .dbf found in zip")
}

func (d *HTTPDownloader) Download(url, dest string) error {
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

		contentType := strings.ToLower(resp.Header.Get("Content-Type"))
		isZip := strings.Contains(contentType, "zip")

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

		// ZIP → распаковываем
		if isZip {
			log.Println("ZIP detected, extracting DBF")

			dbfPath, err := unzipSingleFile(tmp, filepath.Dir(dest))
			os.Remove(tmp)

			if err != nil {
				lastErr = err
				continue
			}

			log.Println("Extracted DBF:", dbfPath)
			return nil
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
