package codefresh

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var stepTypesNamePrefix = "TerraformAccTest_"

// Unit Testing
func TestCleanUpStepFromTransientValues(t *testing.T) {
	stepTypes := cfclient.StepTypes{
		Metadata: map[string]interface{}{
			"accountId":  "test",
			"created_at": "test",
			"id":         "test",
			"updated_at": "test",
			"latest":     true,
			"name":       "originalTestName",
			"version":    "oldVersion",
		},
	}
	newName := "newTestName"
	newVersion := "newVersion"
	cleanUpStepFromTransientValues(&stepTypes, newName, newVersion)
	for _, attributeName := range []string{"created_at", "accountId", "id", "updated_at"} {
		if _, ok := stepTypes.Metadata[attributeName]; ok {
			t.Errorf("Attribute %s wasn't removed from the Metadata %v.", attributeName, stepTypes)
		}
	}
	if stepTypes.Metadata["name"] != newName {
		t.Errorf("Name wasn't updated in Metadata. Expected %s found %s.", newName, stepTypes.Metadata["name"])
	}
	if stepTypes.Metadata["version"] != newVersion {
		t.Errorf("Version wasn't updated in Metadata. Expected %s found %s.", newVersion, stepTypes.Metadata["version"])
	}
}

func TestNormalizeYamlStringStepTypes(t *testing.T) {
	testFile := "../testing/fixtures/testStepWithRuntimeData.yaml"
	yamlString, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Unable to open test file  %s. Err: #%v ", testFile, err)
	}
	if normalizedYAML, error := normalizeYamlStringStepTypes(string(yamlString)); error == nil {
		if strings.Contains(normalizedYAML, "latest: true") {
			t.Errorf("Latest attribute wasn't removed from Metadata. %s.", normalizedYAML)
		}
		if strings.Contains(normalizedYAML, "name: test/step") {
			t.Errorf("Name attribute wasn't removed from Metadata. %s.", normalizedYAML)
		}
		if strings.Contains(normalizedYAML, "version: 0.0.0") {
			t.Errorf("Version attribute wasn't removed from Metadata. %s.", normalizedYAML)
		}
	} else {
		t.Errorf("Error while normalising Yaml string for StepTypes%s.", error)
	}
}

func TestSortVersions(t *testing.T) {
	versions := []string{"2.13.0", "1.0.0", "1.0.23", "0.12.1", "1.0.8", "2.8.0"}
	sortedVersions := []string{"0.12.1", "1.0.0", "1.0.8", "1.0.23", "2.8.0", "2.13.0"}
	sortedCollection := sortVersions(versions)
	for index, item := range sortedCollection {
		checkVersion, err := semver.NewVersion(sortedVersions[index])
		if err != nil {
			t.Errorf("Error parsing checkVersion: %s. Err = %v", checkVersion, err)
		}
		if item.Compare(checkVersion) != 0 {
			t.Errorf("Expected version: %s at index %d found %s.", checkVersion, index, item.String())
		}
	}
}

func TestExtractSteps(t *testing.T) {
	testFile := "../testing/fixtures/testStepTypesOrder.yaml"
	yamlString, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Unable to read file %s", testFile)
	}
	orderedSteps := extractSteps(string(yamlString))
	expectedOrdeer := []string{"first_message", "check_second_message_order_maintained"}
	for index, stepName := range orderedSteps.Keys() {
		if stepName != expectedOrdeer[index] {
			t.Errorf("Expected step %s in position %d but got %s", expectedOrdeer[index], index, stepName)
		}
	}
}

// Acceptance testing
func TestAccCodefreshStepTypes(t *testing.T) {
	// Adding check if we are executing Acceptance test
	// This is needed to ensure we have cfclient initialised so that we can retrieve the accountName dynamically
	if os.Getenv("TF_ACC") == "1" {
		apiClient := testAccProvider.Meta().(*cfclient.Client)
		var accountName string
		if account, err := apiClient.GetCurrentAccount(); err == nil {
			accountName = account.Name
		} else {
			log.Fatalf("Error, unable to retrieve current account name: %s", err)
		}
		name := accountName + "/" + stepTypesNamePrefix + acctest.RandString(10)
		resourceName := "codefresh_step_types.test"
		contentStepsV1, err := os.ReadFile("../testing/fixtures/testSteps.yaml")
		if err != nil {
			log.Fatal(err)
		}
		contentStepsV2, err := os.ReadFile("../testing/fixtures/testStepsTemplate.yaml")
		if err != nil {
			log.Fatal(err)
		}
		stepTypesV1 := string(contentStepsV1)
		stepTypesV2 := string(contentStepsV2)

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckCodefreshStepTypesDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCodefreshStepTypesConfig(name, "0.0.1", stepTypesV1, "0.0.2", stepTypesV2),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckCodefreshStepTypesExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "name", name),
						resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "0.0.1"),
						resource.TestCheckResourceAttr(resourceName, "version.0.step_types_yaml", stepTypesV1),
						resource.TestCheckResourceAttr(resourceName, "version.1.version_number", "0.0.2"),
						resource.TestCheckResourceAttr(resourceName, "version.1.step_types_yaml", stepTypesV2),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
					// `codefresH_step_types` cannot retrieve `version` on Read only using d.GetId(),
					// hence ImportStateVerify will always retrieve an unequilvalent value
					// for `version`. We ignore this field for the purpose of the test.
					// See: https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/import#importstateverifyignore-1
					ImportStateVerifyIgnore: []string{"version"},
				},
			},
		})
	}

}

func testAccCheckCodefreshStepTypesExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		stepTypeID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetStepTypes(stepTypeID)

		if err != nil {
			return fmt.Errorf("error fetching step types with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCodefreshStepTypesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*cfclient.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "codefresh_step_types" {
			continue
		}

		_, err := apiClient.GetStepTypes(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

// CONFIGS
func testAccCodefreshStepTypesConfig(rName, version1, stepTypesYaml1, version2, stepTypesYaml2 string) string {
	return fmt.Sprintf(`
resource "codefresh_step_types" "test" {
  name = "%s"
  version {
	version_number  = "%s"
	step_types_yaml = %#v
  }
  version {
	version_number  = "%s"
	step_types_yaml = %#v
  }
}
`, rName, version1, stepTypesYaml1, version2, stepTypesYaml2)
}
