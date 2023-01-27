package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListStevenlu2Resource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListStevenlu2ResourceConfig("resourceStevenlu2Test", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_stevenlu2.test", "should_monitor", "true"),
					resource.TestCheckResourceAttrSet("whisparr_import_list_stevenlu2.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListStevenlu2ResourceConfig("resourceStevenlu2Test", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_stevenlu2.test", "should_monitor", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_import_list_stevenlu2.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListStevenlu2ResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "whisparr_import_list_stevenlu2" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		should_monitor = %s
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		source = 0
		min_score = 5
	}`, monitor, name)
}
