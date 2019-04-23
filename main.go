package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"k8s-grader/grader"
	"k8s-grader/utils"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	api := os.Getenv("API_URI")
	testDataStr := os.Getenv("TEST_DATA")
	secret := os.Getenv("JOB_SECRET")

	fmt.Println("\n[Brian]: ENV Variables:")
	fmt.Printf("API: %s\nSecret: %s\nData: %s\n", api, secret, testDataStr)

	var testData grader.TestData
	if err := json.Unmarshal([]byte(testDataStr), &testData); err != nil {
		panic(err)
	}

	// Make Folder
	mkdirCmd := exec.Command("mkdir", "-p", "/tmp/job")
	if _, err := mkdirCmd.Output(); err != nil {
		panic(err)
	}

	// Change CWD
	if err := os.Chdir("/tmp/job"); err != nil {
		panic(err)
	}

	// Download and extract supporting files
	path := "/tmp/job/sup.tar.gz"
	url := fmt.Sprintf("%s/job/%s/assignment/%s/supportingfiles/download", api, secret, testData.AssignmentID)
	if err := utils.DownloadAndExtract(url, path, "tarball"); err != nil {
		panic(err)
	}

	// Download and extract submission
	path = "/tmp/job/sub.tar.gz"
	url = fmt.Sprintf("%s/job/%s/submission/%s/download", api, secret, testData.SubmissionID)
	if err := utils.DownloadAndExtract(url, path, "tarball"); err != nil {
		panic(err)
	}

	// Build
	fmt.Println("\n[Brian]: Starting Build Script")
	buildOut, err := testData.Build()
	if err != nil {
		panic(err)
	}

	for k, v := range buildOut {
		fmt.Printf("[Brian]: Build - Exec '%s' Output:\n%s", k, v)
	}

	// Grade the assignment
	results, err := testData.Grade()
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(" ", "  ")
	_ = enc.Encode(results)
	fmt.Printf("\n[Brian]: Test Results: \n%+v", string(buf.String()))

	// Send grade back to court-herald

	fmt.Printf("\n[Brian]: Finished Successfully.\n")
	return
}
