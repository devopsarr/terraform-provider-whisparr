package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRootFolderResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRootFolderResourceConfig("/config"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_root_folder.test", "path", "/config"),
					resource.TestCheckResourceAttrSet("whisparr_root_folder.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRootFolderResourceConfig("/config/logs"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_root_folder.test", "path", "/config/logs"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_root_folder.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRootFolderResourceConfig(path string) string {
	return fmt.Sprintf(`
		resource "whisparr_root_folder" "test" {
  			path = "%s"
		}
	`, path)
}