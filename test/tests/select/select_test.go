package select_test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
  util "github.com/rancher/terraform-aws-server/test/tests"
)

func TestSelectServer(t *testing.T) {
	t.Parallel()
	uniqueID := os.Getenv("IDENTIFIER") + "-" + random.UniqueId()

	category := "select"
	directory := "server"
	region := "us-west-1"
	owner := "terraform-ci@suse.com"
	terraformOptions, keyPair := util.Setup(t, category, directory, region, owner, uniqueID)

	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
	terraformOptions.SshAgent = sshAgent
	defer util.Teardown(t, category, directory, keyPair, sshAgent, uniqueID, terraformOptions)
	delete(terraformOptions.Vars, "key_name")
	delete(terraformOptions.Vars, "key")
	terraform.InitAndApply(t, terraformOptions)
}
func TestSelectImage(t *testing.T) {
	t.Parallel()
	uniqueID := os.Getenv("IDENTIFIER") + "-" + random.UniqueId()

	category := "select"
	directory := "image"
	region := "us-west-1"
	owner := "terraform-ci@suse.com"
	terraformOptions, keyPair := util.Setup(t, category, directory, region, owner, uniqueID)

	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
	terraformOptions.SshAgent = sshAgent
	defer util.Teardown(t, category, directory, keyPair, sshAgent, uniqueID, terraformOptions)
	delete(terraformOptions.Vars, "key_name")
	delete(terraformOptions.Vars, "key")
	terraform.InitAndApply(t, terraformOptions)
}
func TestSelectAll(t *testing.T) {
	t.Parallel()
	uniqueID := os.Getenv("IDENTIFIER") + "-" + random.UniqueId()

	category := "select"
	directory := "all"
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-west-1"
	}
	owner := "terraform-ci@suse.com"
	terraformOptions, keyPair := util.Setup(t, category, directory, region, owner, uniqueID)
	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
	terraformOptions.SshAgent = sshAgent
	defer util.Teardown(t, category, directory, keyPair, sshAgent, uniqueID, terraformOptions)
	delete(terraformOptions.Vars, "key_name")
	delete(terraformOptions.Vars, "key")
	terraform.InitAndApply(t, terraformOptions)
}

// leaving this as an example of a multi-phase test
// func TestAssociation(t *testing.T) {
// 	// in this test we are going to select everything in the server module, but force the association of a new security group onto the selected server
// 	t.Parallel()
// 	//domain := os.Getenv("DOMAIN")
// 	uniqueID := os.Getenv("IDENTIFIER")
// 	if uniqueID == "" {
// 		uniqueID = random.UniqueId()
// 	}
// 	category := "overrides"
// 	directory := "association"
// 	region := "us-west-1"
// 	owner := "terraform-ci@suse.com"

// 	// multi part terraform, setup is an independent module in the test directory that sets up this test
// 	setupDirectory := fmt.Sprintf("%s/setup", directory)
// 	setupTerraformOptions, setupKeyPair := setup(t, category, setupDirectory, region, owner, uniqueID)
// 	setupSshAgent := ssh.SshAgentWithKeyPair(t, setupKeyPair.KeyPair)
// 	setupTerraformOptions.SshAgent = setupSshAgent
// 	defer setupSshAgent.Stop()
// 	defer util.Teardown(t, category, setupDirectory, setupKeyPair)
// 	defer terraform.Destroy(t, setupTerraformOptions)
// 	terraform.InitAndApply(t, setupTerraformOptions)
// 	serverId := terraform.Output(t, setupTerraformOptions, "id")
// 	uniqueId := terraform.Output(t, setupTerraformOptions, "identifier")
// 	keyName := terraform.Output(t, setupTerraformOptions, "key_name")

// 	// after setup completes we can run the actual test, passing in the server id from setup
// 	terraformOptions, keyPair := util.Setup(t, category, directory, region, owner, uniqueId)
// 	defer util.Teardown(t, category, directory, keyPair)
// 	defer terraform.Destroy(t, terraformOptions)
// 	terraformOptions.Vars["identifier"] = uniqueId
// 	terraformOptions.Vars["server"] = serverId
// 	terraformOptions.Vars["key_name"] = keyName
// 	delete(terraformOptions.Vars, "key")
// 	terraform.InitAndApply(t, terraformOptions)
// }
