package googleworkspace

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDynamicGroup_basic(t *testing.T) {
	t.Parallel()

	domainName := os.Getenv("GOOGLEWORKSPACE_DOMAIN")

	if domainName == "" {
		t.Skip("GOOGLEWORKSPACE_DOMAIN needs to be set to run this test")
	}

	testDynamicGroupVals := map[string]interface{}{
		"domainName": domainName,
		"email":      fmt.Sprintf("tf-test-%s", acctest.RandString(10)),
		"query":      "user.organizations.exists(org, org.department=='engineering')",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDynamicGroup_basic(testDynamicGroupVals),
			},
			{
				// TestStep imports by `id` by default
				ResourceName:      "googleworkspace_dynamic_group.my-dynamic-group",
				ImportState:       true,
				ImportStateCheck:  checkDynamicGroupImportState(),
				ImportStateVerify: true,
			},
		},
	})
}

func checkDynamicGroupImportState() resource.ImportStateCheckFunc {
	return resource.ImportStateCheckFunc(
		func(state []*terraform.InstanceState) error {
			if len(state) > 1 {
				return fmt.Errorf("state should only contain one dynamic group resource, got: %d", len(state))
			}

			return nil
		},
	)
}

func TestAccResourceDynamicGroup_full(t *testing.T) {
	t.Parallel()

	domainName := os.Getenv("GOOGLEWORKSPACE_DOMAIN")

	if domainName == "" {
		t.Skip("GOOGLEWORKSPACE_DOMAIN needs to be set to run this test")
	}

	testDynamicGroupVals := map[string]interface{}{
		"domainName": domainName,
		"email":      fmt.Sprintf("tf-test-%s", acctest.RandString(10)),
		"query":      "user.organizations.exists(org, org.department=='engineering')",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDynamicGroup_full(testDynamicGroupVals),
			},
			{
				ResourceName:      "googleworkspace_dynamic_group.my-dynamic-group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceDynamicGroup_fullUpdate(testDynamicGroupVals),
			},
			{
				ResourceName:      "googleworkspace_dynamic_group.my-dynamic-group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceDynamicGroup_basic(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-dynamic-group" {
  email = "%{email}@%{domainName}"
  query = "%{query}"
}
`, testDynamicGroupVals)
}

func testAccResourceDynamicGroup_full(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-dynamic-group" {
  email = "%{email}@%{domainName}"
  name  = "tf-test-name"
  description = "my test description"
  query = "%{query}"
}
`, testDynamicGroupVals)
}

func testAccResourceDynamicGroup_fullUpdate(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-dynamic-group" {
  email = "%{email}@%{domainName}"
  name  = "tf-new-name"
  description = "my new description"
  query = "%{query}"
}
`, testDynamicGroupVals)
}
