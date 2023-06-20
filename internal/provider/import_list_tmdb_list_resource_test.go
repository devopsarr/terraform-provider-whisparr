package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListTMDBListResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListTMDBListResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTMDBListResourceConfig("resourceTMDListTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_tmdb_list.test", "should_monitor", "true"),
					resource.TestCheckResourceAttrSet("whisparr_import_list_tmdb_list.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListTMDBListResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListTMDBListResourceConfig("resourceTMDListTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_tmdb_list.test", "should_monitor", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_import_list_tmdb_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTMDBListResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "whisparr_import_list_tmdb_list" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		should_monitor = %s
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		list_id = "11842"
	}`, monitor, name)
}
