package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"k8s-grader/grader"
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

	fmt.Printf("API: %s\nData: %s\nSecret: %s\n", api, testDataStr, secret)

	var testData grader.TestData
	if err := json.Unmarshal([]byte(testDataStr), &testData); err != nil {
		panic(err)
	}

	// Change working directory
	mkdirCmd := exec.Command("mkdir", "-p", "/tmp/job")
	if _, err := mkdirCmd.Output(); err != nil {
		panic(err)
	}

	results, err := testData.Grade()
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(" ", "  ")
	_ = enc.Encode(results)
	fmt.Printf("\nTest Results: \n%+v", string(buf.String()))

	return
	// Download assignment
	// job/:secret/submission/:sid/download
	// path := "/tmp/submission.tar.gz"
	// url := fmt.Sprintf("%s/job/%s/submission/%s/download", api, secret, subID)
	// if err := utils.DownloadFile(path, url); err != nil {
	// 	panic(err)
	// }

	// Download supporting files
	// job/:secret/assignment/:aid/supportingfiles/download
	// path = "/tmp/support.tar.gz"
	// url = fmt.Sprintf("%s/job/%s/assignment/%s/supportingfiles/download", api, secret, assID)
	// if err := utils.DownloadFile(path, url); err != nil {
	// 	panic(err)
	// }

	// Grade Submission

	// Send grade back to court-herald

	// fmt.Printf("finished!")
}
