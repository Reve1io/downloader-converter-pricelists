package source

type Downloader interface {
	Download(url, dest string) error
}
