package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListWhisparrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListWhisparrResourceConfig("resourceWhisparrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_whisparr.test", "should_monitor", "true"),
					resource.TestCheckResourceAttrSet("whisparr_import_list_whisparr.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListWhisparrResourceConfig("resourceWhisparrTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_whisparr.test", "should_monitor", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_import_list_whisparr.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListWhisparrResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "whisparr_import_list_whisparr" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		should_monitor = %s
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		base_url = "http://127.0.0.1:6969"
		api_key = "testAPIKey"
		tag_ids = [1,2]
		profile_ids = [1]
	}`, monitor, name)
}
