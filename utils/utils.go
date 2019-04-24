package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"k8s-grader/grader"
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("unable to download file from backend")
	}

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

// SendResults sends the results via a PATCH request to the specified URL
func SendResults(url string, results []grader.WorkerResult) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(" ", "  ")
	_ = enc.Encode(results)

	req, err := http.NewRequest("PATCH", url, buf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("invalid response from backend")
	}

	return nil
}
