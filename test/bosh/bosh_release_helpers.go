package bosh

import (
	. "github.com/onsi/gomega"
	"github.com/vito/cmdtest"
	"github.com/vito/cmdtest/matchers"

	"os"
	"strings"
)

func BoshDeleteRelease() {
	result := BoshCombinedOutput("-n", "delete", "release", "bosh-aws-tools")
	Expect(result).To(cmdtest_matchers.SayBranches(
		cmdtest.ExpectBranch{
			Pattern:  "Deleted",
			Callback: func() {},
		},
		cmdtest.ExpectBranch{
			Pattern:  "doesn't exist", // already clean
			Callback: func() {},
		},
	))
}

func BoshCreateRelease() {
	result := BoshInDir(projectDir(), "create", "release", "--force")
	Expect(result).To(cmdtest_matchers.Say("Release manifest:"))
}

func BoshUploadRelease() {
	// Make sure it uploads new release since we cannot rely
	// on already uploaded release having the same bits
	// as the newly created release (even if their versions match)
	result := BoshInDir(projectDir(), "-n", "upload", "release")
	Expect(result).To(cmdtest_matchers.Say("Release uploaded"))
}

func projectDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return strings.Split(cwd, "/test")[0]
}
