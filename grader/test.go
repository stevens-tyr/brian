package grader

import (
	"errors"
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

	for _, test := range data.Tests {
		go Worker(test, results)
	}

	var testResults []WorkerResult
	for range data.Tests {
		select {
		case res := <-results:
			testResults = append(testResults, res)
		case <-time.After(2 * time.Minute):
			return nil, errors.New("Timed out while running tests. (One Minute)")
		}
	}

	return testResults, nil
}
