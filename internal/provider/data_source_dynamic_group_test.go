package googleworkspace

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDynamicGroup_withId(t *testing.T) {
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
				Config: testAccDataSourceDynamicGroup_withId(testDynamicGroupVals),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.googleworkspace_dynamic_group.my-new-dynamic-group", "email", Nprintf("%{email}@%{domainName}", testDynamicGroupVals)),
					resource.TestCheckResourceAttr(
						"data.googleworkspace_dynamic_group.my-new-dynamic-group", "name", testDynamicGroupVals["email"].(string)),
				),
			},
		},
	})
}

func TestAccDataSourceDynamicGroup_withEmail(t *testing.T) {
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
				Config: testAccDataSourceDynamicGroup_withEmail(testDynamicGroupVals),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.googleworkspace_dynamic_group.my-new-dynamic-group", "email", Nprintf("%{email}@%{domainName}", testDynamicGroupVals)),
					resource.TestCheckResourceAttr(
						"data.googleworkspace_dynamic_group.my-new-dynamic-group", "name", testDynamicGroupVals["email"].(string)),
				),
			},
		},
	})
}

func testAccDataSourceDynamicGroup_withId(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-new-dynamic-group" {
  email = "%{email}@%{domainName}"
  query = "%{query}"
}

data "googleworkspace_dynamic_group" "my-new-dynamic-group" {
  id = googleworkspace_dynamic_group.my-new-dynamic-group.id
}
`, testDynamicGroupVals)
}

func testAccDataSourceDynamicGroup_withEmail(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-new-dynamic-group" {
  email = "%{email}@%{domainName}"
  query = "%{query}"
}

data "googleworkspace_dynamic_group" "my-new-dynamic-group" {
  email = googleworkspace_dynamic_group.my-new-dynamic-group.email
}
`, testDynamicGroupVals)
}
