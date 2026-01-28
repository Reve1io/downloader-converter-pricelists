package source

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPDownloader struct {
	Host    string
	Port    int
	User    string
	Pass    string
	Timeout time.Duration
	Retries int
}

func (d *FTPDownloader) Download(remoteFile, dest string) error {
	var lastErr error

	for attempt := 1; attempt <= d.Retries; attempt++ {
		conn, err := ftp.Dial(
			d.Host+":"+fmt.Sprint(d.Port),
			ftp.DialWithTimeout(d.Timeout),
		)
		if err != nil {
			lastErr = err
			continue
		}

		if err := conn.Login(d.User, d.Pass); err != nil {
			conn.Quit()
			lastErr = err
			continue
		}

		r, err := conn.Retr(remoteFile)
		if err != nil {
			conn.Quit()
			lastErr = err
			continue
		}

		tmp := dest + ".tmp"
		f, err := os.Create(tmp)
		if err != nil {
			r.Close()
			conn.Quit()
			return err
		}

		_, err = io.Copy(f, r)
		f.Close()
		r.Close()
		conn.Quit()

		if err != nil {
			lastErr = err
			continue
		}

		return os.Rename(tmp, dest)
	}

	return lastErr
}
