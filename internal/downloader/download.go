package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadURL(url string, filePathArgs ...string) (file *os.File, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close response body: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var filePath string
	if len(filePathArgs) == 0 || len(filePathArgs[0]) == 0 {
		filePath = filepath.Base(resp.Request.URL.String())
	} else {
		filePath = filePathArgs[0]
	}
	file, err = os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return file, nil
}
