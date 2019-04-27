package grader

import (
	"errors"
	"os/exec"
	"strings"
	"time"
)

// Test stores information about each test case
type Test struct {
	Name           string `json:"name" binding:"required"`
	ExpectedOutput string `json:"expectedOutput" binding:"required"`
	StudentFacing  bool   `json:"studentFacing" binding:"exists"`
	TestCMD        string `son:"testCMD" binding:"required"`
}

// TestData stores all the information about the current test
type TestData struct {
	AssignmentID string `json:"assignmentID" binding:"required"`
	SubmissionID string `json:"submissionID" binding:"required"`
	TestBuildCMD string `json:"testBuildCMD" binding:"required"`
	Tests        []Test `json:"tests" binding:"required"`
}

// Grade runs Student/AdminFacing tests in goroutines
func (data TestData) Grade() ([]WorkerResult, error) {
	results := make(chan WorkerResult, len(data.Tests))

	for id, test := range data.Tests {
		go Worker(id, test, results)
	}

	testResults := make([]WorkerResult, len(data.Tests))
	for range data.Tests {
		select {
		case res := <-results:
			testResults[res.ID] = res
		case <-time.After(2 * time.Minute):
			return nil, errors.New("Timed out while running tests. (Two Minutes)")
		}
	}

	return testResults, nil
}

// Build will build the specified
func (data TestData) Build() (map[string]string, error) {
	cmds := strings.Split(data.TestBuildCMD, "\n")
	totalOutput := make(map[string]string)
	for _, cmd := range cmds {
		args := strings.Fields(cmd)
		if args[0][0] != '#' {
			buildCmd := exec.Command(args[0], args[1:]...)
			buildOut, err := buildCmd.Output()
			if err != nil {
				return nil, err
			}
			totalOutput[cmd] = string(buildOut)
		}
	}

	return totalOutput, nil
}
