package os_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
  util "github.com/rancher/terraform-aws-server/test/tests"
)

func TestOs(t *testing.T) {
	id := os.Getenv("IDENTIFIER")
	if id == "" {
		id = random.UniqueId()
	}
	uniqueID := id + "-" + random.UniqueId()
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-west-2"
	}
	owner := "terraform-ci@suse.com"
  imageType := os.Getenv("IMAGE")

	// get the image list from the imagetype example
	category := "imagetype"
	directory := "basic"
	imageTypesTerraformOptions, keyPair := util.Setup(t, category, directory, region, owner, uniqueID)
	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
	// don't pass key or key_name to the image module
	delete(imageTypesTerraformOptions.Vars, "key")
	delete(imageTypesTerraformOptions.Vars, "key_name")
	_, err := terraform.InitAndApplyE(t, imageTypesTerraformOptions)
	if err != nil {
		util.Teardown(t, category, directory, keyPair, sshAgent, uniqueID, imageTypesTerraformOptions)
		t.Error(err)
		t.Fail()
	}
	info := terraform.OutputMap(t, imageTypesTerraformOptions, "image_names")
	images := keys(info)
	util.Teardown(t, category, directory, keyPair, sshAgent, uniqueID, imageTypesTerraformOptions)
	for k := range images {
		image := images[k].String()
    if imageType != "" && imageType != image {
      continue
    }

    t.Run(image, func(t *testing.T) {
			t.Parallel()
			t.Logf("Running test for %s", image)
			uniqueID := id + "-" + random.UniqueId()
			category := "os"
			directory := "all"
			terraformOptions, keyPair := util.Setup(t, category, directory, region, owner, uniqueID)
			sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
			terraformOptions.SshAgent = sshAgent
			terraformOptions.Vars["image"] = image
			_, err := terraform.InitAndApplyE(t, terraformOptions)
			if err != nil {
				util.Teardown(t, category, directory, keyPair, sshAgent, uniqueID, terraformOptions)
				t.Error(err)
				t.Fail()
			}
			out := terraform.OutputAll(t, terraformOptions)
			t.Logf("out: %v", out)
			outputServer, ok := out["server"].(map[string]interface{})
			assert.True(t, ok, fmt.Sprintf("Wrong data type for 'server', expected map[string], got %T", out["server"]))
			outputImage, ok := out["image"].(map[string]interface{})
			assert.True(t, ok, fmt.Sprintf("Wrong data type for 'image', expected map[string], got %T", out["image"]))
			assert.NotEmpty(t, outputServer["public_ip"], "The 'server.public_ip' is empty")
			assert.NotEmpty(t, outputImage["id"], "The 'image.id' is empty")
			util.Teardown(t, category, directory, keyPair, sshAgent, uniqueID, terraformOptions)
		})
	}
}

func keys(m map[string]string) []reflect.Value {
	return reflect.ValueOf(m).MapKeys()
}
