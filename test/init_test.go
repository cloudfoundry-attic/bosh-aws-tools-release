package route53_backup_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/cloudfoundry/bosh-aws-tools/test/config"
)

var Config = config.LoadAndValidate()

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}
