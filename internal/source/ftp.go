package source

import (
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPDownloader struct {
	Addr     string // host:port
	User     string
	Password string
	Timeout  time.Duration
}

func (d *FTPDownloader) Download(remotePath, localPath string) error {
	c, err := ftp.Dial(d.Addr, ftp.DialWithTimeout(d.Timeout))
	if err != nil {
		return err
	}
	defer c.Quit()

	if err := c.Login(d.User, d.Password); err != nil {
		return err
	}

	r, err := c.Retr(remotePath)
	if err != nil {
		return err
	}
	defer r.Close()

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return err
	}

	return nil
}
