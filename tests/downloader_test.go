package tests

import (
	"downloader-converter-pricelists/internal/config"
	"downloader-converter-pricelists/internal/source"
	"downloader-converter-pricelists/internal/utils"
	"os"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {

	cfg, err := config.Load("./../configs/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	remotePath := "/data/unload_all.xlsx"
	localPath := "test_input/unload_all.xlsx"

	utils.CleanDir(localPath)

	sftpDl := &source.SFTPDownloader{
		Addr:     cfg.SFTP.Addr,
		User:     cfg.SFTP.User,
		Password: cfg.SFTP.Password,
		Timeout:  time.Duration(cfg.SFTP.TimeoutSeconds),
	}

	if err := sftpDl.Download(remotePath, localPath); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(localPath)
	if err != nil {
		t.Fatal(err)
	}

	if info.Size() == 0 {
		t.Fatal("Downloaded files is empty")
	}

	t.Log("File downloaded success!")
	t.Log(info)
}
