package grader

import (
	"os/exec"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// WorkerResult is used in the channel to hold worker data
type WorkerResult struct {
	Panicked bool   `json:"panicked" binding:"required"`
	Passed   bool   `json:"passed" binding:"required"`
	Output   string `json:"output" binding:"required"`
	HTML     string `json:"html" binding:"required"`
	TestCMD  string `json:"testCMD" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// Worker is used in a goroutine to exec a cmd in the job
func Worker(t Test, results chan WorkerResult) {
	args := strings.Fields(t.TestCMD)
	testCmd := exec.Command(args[0], args[1:]...)
	testCmd.Dir = "/tmp/job"
	testOut, err := testCmd.Output()
	if err != nil {
		results <- WorkerResult{Panicked: true, Passed: false, Output: "", HTML: "", TestCMD: t.TestCMD, Name: t.Name}
	} else {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(t.ExpectedOutput, string(testOut), false)
		passed := true
		if len(diffs) > 1 {
			passed = false
		}
		results <- WorkerResult{Panicked: false, Passed: passed, Output: string(testOut), HTML: dmp.DiffPrettyHtml(diffs), TestCMD: t.TestCMD, Name: t.Name}
	}
}
