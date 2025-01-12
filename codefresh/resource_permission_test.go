package codefresh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	funk "github.com/thoas/go-funk"
)

func TestAccCodefreshPermissionConfig(t *testing.T) {
	resourceName := "codefresh_permission.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshPermissionConfig("create", "pipeline", "null", []string{"production", "test"}, "all"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action", "create"),
					resource.TestCheckResourceAttr(resourceName, "resource", "pipeline"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "production"),
					resource.TestCheckResourceAttr(resourceName, "related_resource", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "test"),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "all"),
				),
			},
			{
				Config: testAccCodefreshPermissionConfig("create", "pipeline", "project", []string{"production", "*"}, "any"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action", "create"),
					resource.TestCheckResourceAttr(resourceName, "resource", "pipeline"),
					resource.TestCheckResourceAttr(resourceName, "related_resource", "project"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "production"),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "any"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckCodefreshPermissionExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		permissionID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetPermissionByID(permissionID)

		if err != nil {
			return fmt.Errorf("error fetching permission with ID %s. %s", permissionID, err)
		}
		return nil
	}
}

// CONFIGS
func testAccCodefreshPermissionConfig(action, resource, relatedResource string, tags []string, ruleType string) string {
	escapeString := func(str string) string {
		if str == "null" {
			return str // null means Terraform should ignore this field
		}
		return fmt.Sprintf(`"%s"`, str)
	}
	tagsEscaped := funk.Map(tags, escapeString).([]string)

	return fmt.Sprintf(`
	data "codefresh_team" "users" {
		name = "users"
	}

	resource "codefresh_permission" "test" {
		team             = data.codefresh_team.users.id
		action           = %s
		resource         = %s
		related_resource = %s
		tags             = [%s]
		rule_type        = %s
	}
`, escapeString(action), escapeString(resource), escapeString(relatedResource), strings.Join(tagsEscaped[:], ","), escapeString(ruleType))
}
