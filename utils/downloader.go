package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadImage downloads an image from a URL and saves it to a local file.
func DownloadImage(url, destDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response is OK
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	// Create the destination file
	fileName := filepath.Join(destDir, filepath.Base(url))
	out, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the data to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
