package grader

import (
	"os/exec"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// WorkerResult is used in the channel to hold worker data
type WorkerResult struct {
	ID            int    `bson:"id" json:"id" binding:"required"`
	Panicked      bool   `bson:"panicked" json:"panicked" binding:"required"`
	Passed        bool   `bson:"passed" json:"passed" binding:"required"`
	StudentFacing bool   `bson:"studentFacing" json:"studentFacing" binding:"required"`
	Output        string `bson:"output" json:"output" binding:"required"`
	HTML          string `bson:"html" json:"html" binding:"required"`
	TestCMD       string `bson:"testCMD" json:"testCMD" binding:"required"`
	Name          string `bson:"name" json:"name" binding:"required"`
}

// Worker is used in a goroutine to exec a cmd in the job
func Worker(id int, t Test, results chan WorkerResult) {
	args := strings.Fields(t.TestCMD)
	testCmd := exec.Command(args[0], args[1:]...)
	testCmd.Dir = "/tmp/job"
	testOut, err := testCmd.Output()
	if err != nil {
		results <- WorkerResult{ID: id, Panicked: true, Passed: false, StudentFacing: t.StudentFacing, Output: "", HTML: "", TestCMD: t.TestCMD, Name: t.Name}
	} else {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(testOut), t.ExpectedOutput, false)
		passed := true
		if len(diffs) > 1 {
			passed = false
		}
		results <- WorkerResult{ID: id, Panicked: false, Passed: passed, StudentFacing: t.StudentFacing, Output: string(testOut), HTML: dmp.DiffPrettyHtml(diffs), TestCMD: t.TestCMD, Name: t.Name}
	}
}
