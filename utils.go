package ffmpeg_downloader

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v74/github"
)

func getLatestVersion() string {
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "vegidio", "ffmpeg-downloader")
	if err != nil {
		return ""
	}

	return release.GetTagName()
}

func download(url, fileName string) error {
	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	directory := filepath.Dir(fileName)
	if err = os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %s", err)
	}

	// Create the file on disk
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dest string) error {
	// Open the zip file specified by src
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Ensure the destination directory exists
	if err = os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// Iterate through each file in the zip archive
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip vulnerability: ensure that the file path is within the destination directory
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create a directory if it doesn't exist
			if err = os.MkdirAll(fpath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		// Ensure the parent directory exists
		if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		// Create the destination file
		outFile, fErr := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if fErr != nil {
			return fErr
		}

		// Make the file executable
		if err = os.Chmod(fpath, 0755); err != nil {
			return err
		}

		// Open the file inside the zip archive
		rc, fErr := f.Open()
		if fErr != nil {
			outFile.Close()
			return fErr
		}

		// Copy file contents from the zip archive to a destination file
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
