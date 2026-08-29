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

func TestAccResourceDynamicGroup_aliases(t *testing.T) {
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
				Config: testAccResourceDynamicGroup_aliases(testDynamicGroupVals),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.#", "2"),
					resource.TestCheckTypeSetElemAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.*",
						Nprintf("%{email}-alias-1@%{domainName}", testDynamicGroupVals)),
					resource.TestCheckTypeSetElemAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.*",
						Nprintf("%{email}-alias-2@%{domainName}", testDynamicGroupVals)),
				),
			},
			{
				Config: testAccResourceDynamicGroup_aliasesUpdate(testDynamicGroupVals),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.#", "2"),
					resource.TestCheckTypeSetElemAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.*",
						Nprintf("%{email}-alias-2@%{domainName}", testDynamicGroupVals)),
					resource.TestCheckTypeSetElemAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.*",
						Nprintf("%{email}-new-alias@%{domainName}", testDynamicGroupVals)),
				),
			},
		},
	})
}

func TestAccResourceDynamicGroup_renameKeepAlias(t *testing.T) {
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
				Config: testAccResourceDynamicGroup_renameKeepAlias(testDynamicGroupVals),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_dynamic_group.my-dynamic-group", "email",
						Nprintf("%{email}-renamed@%{domainName}", testDynamicGroupVals)),
					resource.TestCheckResourceAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.#", "1"),
					resource.TestCheckTypeSetElemAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.*",
						Nprintf("%{email}@%{domainName}", testDynamicGroupVals)),
				),
			},
		},
	})
}

func TestAccResourceDynamicGroup_renameDropAlias(t *testing.T) {
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
				Config: testAccResourceDynamicGroup_renameDropAlias(testDynamicGroupVals),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_dynamic_group.my-dynamic-group", "email",
						Nprintf("%{email}-renamed@%{domainName}", testDynamicGroupVals)),
					resource.TestCheckResourceAttr("googleworkspace_dynamic_group.my-dynamic-group", "aliases.#", "0"),
				),
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

func testAccResourceDynamicGroup_aliases(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-dynamic-group" {
  email = "%{email}@%{domainName}"
  query = "%{query}"

  aliases = ["%{email}-alias-1@%{domainName}", "%{email}-alias-2@%{domainName}"]
}
`, testDynamicGroupVals)
}

func testAccResourceDynamicGroup_aliasesUpdate(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-dynamic-group" {
  email = "%{email}@%{domainName}"
  query = "%{query}"

  aliases = ["%{email}-alias-2@%{domainName}", "%{email}-new-alias@%{domainName}"]
}
`, testDynamicGroupVals)
}

func testAccResourceDynamicGroup_renameKeepAlias(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-dynamic-group" {
  email = "%{email}-renamed@%{domainName}"
  query = "%{query}"

  aliases = ["%{email}@%{domainName}"]
}
`, testDynamicGroupVals)
}

func testAccResourceDynamicGroup_renameDropAlias(testDynamicGroupVals map[string]interface{}) string {
	return Nprintf(`
resource "googleworkspace_dynamic_group" "my-dynamic-group" {
  email = "%{email}-renamed@%{domainName}"
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
