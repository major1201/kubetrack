package goutils

import (
	"io"
	"net/http"
	"os"
)

// Download is used to conveniently download a file to the certain path
func Download(url string, dest string) (err error) {
	out, fileError := os.Create(dest)
	defer func() {
		_ = out.Close()
	}()
	if fileError != nil {
		err = fileError
		return
	}
	resp, httpError := http.Get(url)
	if httpError != nil {
		err = httpError
		return
	}
	defer resp.Body.Close()
	_, copyError := io.Copy(out, resp.Body)
	if copyError != nil {
		err = copyError
		return
	}
	return err
}
