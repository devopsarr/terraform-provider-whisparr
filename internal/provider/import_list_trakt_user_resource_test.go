package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListTraktUserResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListTraktUserResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTraktUserResourceConfig("resourceTraktUserTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_trakt_user.test", "should_monitor", "true"),
					resource.TestCheckResourceAttrSet("whisparr_import_list_trakt_user.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListTraktUserResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListTraktUserResourceConfig("resourceTraktUserTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_trakt_user.test", "should_monitor", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_import_list_trakt_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTraktUserResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "whisparr_import_list_trakt_user" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		should_monitor = %s
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		auth_user = "User1"
		access_token = "Token"
		trakt_list_type = 0
		limit = 100
	}`, monitor, name)
}
