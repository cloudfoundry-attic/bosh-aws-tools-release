package bosh

import (
	. "github.com/onsi/gomega"
	"github.com/vito/cmdtest"
	"github.com/vito/cmdtest/matchers"

	"io/ioutil"
	"os"
	"text/template"
)

func BoshDeleteDeployment(name string) {
	result := BoshCombinedOutput("-n", "delete", "deployment", name)
	Expect(result).To(cmdtest_matchers.SayBranches(
		cmdtest.ExpectBranch{
			Pattern:  "Deleted deployment",
			Callback: func() {},
		},
		cmdtest.ExpectBranch{
			Pattern:  "doesn't exist", // already clean
			Callback: func() {},
		},
	))
}

func BoshDeployDeployment(tplString string, vals interface{}) *os.File {
	manifestTpl, err := template.New("bosh_manifest").Parse(tplString)
	Expect(err).ShouldNot(HaveOccurred())

	file, err := ioutil.TempFile("", "")
	Expect(err).ShouldNot(HaveOccurred())

	err = manifestTpl.Execute(file, vals)
	Expect(err).ShouldNot(HaveOccurred())

	result := Bosh("-n", "deployment", file.Name())
	Expect(result).To(cmdtest_matchers.Say("Deployment set to"))

	result = Bosh("-n", "deploy")
	Expect(result).To(cmdtest_matchers.Say("Deployed"))

	return file
}
