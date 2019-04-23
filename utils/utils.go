package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// DownloadFile downloads a file to a specific filepath based on url
func DownloadFile(url string, filepath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// ExtractFile extracts a certain filetype based on the filepath and filetype
func ExtractFile(filepath string, filetype string) error {
	switch filetype {
	case "tarball":
		untar := exec.Command("tar", "-xvf", filepath)

		if _, err := untar.Output(); err != nil {
			return err
		}
		return nil
	case "zip":
		return errors.New("zip not implemented")
	default:
		return errors.New("invalid format specified")
	}
}

// DownloadAndExtract calls download and extract :)
func DownloadAndExtract(url string, filepath string, filetype string) error {
	if err := DownloadFile(url, filepath); err != nil {
		return err
	}
	if err := ExtractFile(filepath, filetype); err != nil {
		return err
	}
	return nil
}
