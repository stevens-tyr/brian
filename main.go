package main

import (
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

	backendURL := os.Getenv("BACKEND_URL")
	subID := os.Getenv("SUB_ID")
	assignID := os.Getenv("ASSIGN_ID")
	testBuildCMD := os.Getenv("BUILD_CMD")
	testsString := os.Getenv("TESTS")
	secret := os.Getenv("JOB_SECRET")

	bs := []byte(testsString)
	var tests []grader.Test
	json.Unmarshal(bs, &tests)

	testData := grader.TestData{
		AssignmentID: assignID,
		SubmissionID: subID,
		TestBuildCMD: testBuildCMD,
		Tests:        tests,
	}

	fmt.Println("\n[Brian]: ENV:")
	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	// Make Folder
	mkdirCmd := exec.Command("mkdir", "-p", "/tmp/job")
	if _, err := mkdirCmd.Output(); err != nil {
		utils.EmitError(err)
	}

	// Change CWD
	if err := os.Chdir("/tmp/job"); err != nil {
		utils.EmitError(err)
	}

	// Download and extract supporting files
	path := "/tmp/job/sup.tar.gz"
	url := fmt.Sprintf("%s/job/%s/assignment/%s/supportingfiles/download", backendURL, secret, testData.AssignmentID)
	if err := utils.DownloadAndExtract(url, path, "tarball"); err != nil {
		utils.EmitError(err)
	}

	// Download and extract submission
	path = "/tmp/job/sub.tar.gz"
	url = fmt.Sprintf("%s/job/%s/submission/%s/download", backendURL, secret, testData.SubmissionID)
	if err := utils.DownloadAndExtract(url, path, "tarball"); err != nil {
		utils.EmitError(err)
	}

	// Build
	fmt.Println("\n[Brian]: Starting Build Script")
	buildOut, err := testData.Build()
	if err != nil {
		utils.EmitError(err)
	}

	for k, v := range buildOut {
		fmt.Printf("[Brian]: Build - Exec '%s' Output:\n%s", k, v)
	}

	// Grade the assignment
	results, err := testData.Grade()
	if err != nil {
		utils.EmitError(err)
	}

	// Send grade back to court-herald
	fmt.Printf("\n[Brian]: Sending Test Results to Backend\n")
	url = fmt.Sprintf("%s/job/%s/submission/%s/update", backendURL, secret, testData.SubmissionID)
	err = utils.SendResults(url, results)
	if err != nil {
		utils.EmitError(err)
	}

	fmt.Printf("\n[Brian]: Finished Successfully.\n")
}
