package route53_backup_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/cmdtest"
	. "github.com/vito/cmdtest/matchers"

	"fmt"
	. "github.com/cloudfoundry/bosh-aws-tools/test/bosh"
)

var _ = Describe("route53_backup job", func() {
	type deploymentManifestVals struct {
		DirectorUUID string
		DirectorHost string

		DeploymentName string
		JobName        string

		AwsAccessId         string
		AwsSecretAcccessKey string

		Route53ZoneNames []string
	}

	deploymentConfig := deploymentManifestVals{
		DirectorUUID: Config.DirectorUUID,
		DirectorHost: Config.DirectorHost,

		DeploymentName: "route53_backup",
		JobName:        "route53_backup",

		AwsAccessId:         Config.AwsAccessId,
		AwsSecretAcccessKey: Config.AwsSecretAcccessKey,

		Route53ZoneNames: Config.Route53ZoneNames,
	}

	boshSshToReturnBackupCounts := func() *cmdtest.Session {
		return Bosh(
			"-n", "ssh",
			deploymentConfig.JobName,
			"--gateway_host", deploymentConfig.DirectorHost,
			"--gateway_user", "vcap",
			"ls /var/vcap/store/route53_backup/* | xargs head -1 -q | sort -n | uniq -c",
		)
	}

	cleanUp := func() {
		BoshDeleteDeployment(deploymentConfig.DeploymentName)
		BoshDeleteRelease()
	}

	BeforeEach(cleanUp)
	AfterEach(cleanUp)

	BeforeEach(func() {
		BoshCreateRelease()
		BoshUploadRelease()
	})

	It("starts route53_backup job which periodically backs up route53 zones", func() {
		manifestFile := BoshDeployDeployment(deploymentManifestTplStr, deploymentConfig)

		defer manifestFile.Close()

		for _, zone := range deploymentConfig.Route53ZoneNames {
			// 1+ count for this zone
			countMatch := fmt.Sprintf("\\s+([2-9]|\\d{2,})\\s+\\$ORIGIN %s", zone)

			Eventually(
				boshSshToReturnBackupCounts,
				float64(3*60), // timeout in secs
				float64(10),   // interval in secs
			).Should(Say(countMatch))
		}
	})
})

// Keep template localized to this file
// to stop sharing templates between unrelated tests!
const deploymentManifestTplStr = `
---
name: {{ .DeploymentName }}
director_uuid: {{ .DirectorUUID }}

releases:
- name: bosh-aws-tools
  version: latest

networks:
- name: default
  type: manual
  subnets:
  - range: 10.10.16.0/24
    gateway: 10.10.16.1
    reserved:
    - 10.10.16.2 - 10.10.16.10 # full bosh is .7
    dns:
    - 10.10.16.6
    cloud_properties:
      subnet: subnet-f8744a8c

resource_pools:
- name: default
  stemcell:
    name: bosh-aws-xen-ubuntu
    version: latest
  network: default
  size: 1
  cloud_properties:
    instance_type: m1.small
    availability_zone: us-east-1b

compilation:
  reuse_compilation_vms: true
  workers: 1
  network: default
  cloud_properties:
    instance_type: c1.medium
    availability_zone: us-east-1b

update:
  canaries: 1
  canary_watch_time: 1000 - 90000
  update_watch_time: 1000 - 90000
  max_in_flight: 1
  max_errors: 1

jobs:
- name: {{ .JobName }}
  template: route53_backup
  resource_pool: default
  instances: 1
  networks:
  - name: default

properties:
  route53_backup:
    aws_access_key_id: {{ .AwsAccessId }}
    aws_secret_access_key: {{ .AwsSecretAcccessKey }}
    schedule: "*/1 * * * *"
`
