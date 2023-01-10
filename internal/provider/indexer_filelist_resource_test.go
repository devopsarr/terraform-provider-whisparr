package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerFilelistResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerFilelistResourceConfig("filelistResourceTest", "user"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_filelist.test", "username", "user"),
					resource.TestCheckResourceAttrSet("whisparr_indexer_filelist.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerFilelistResourceConfig("filelistResourceTest", "Username"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_filelist.test", "username", "Username"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_indexer_filelist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerFilelistResourceConfig(name, username string) string {
	return fmt.Sprintf(`
	resource "whisparr_indexer_filelist" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://filelist.io"
		username = "%s"
		passkey = "Pass"
		categories = [4,6,1]
		minimum_seeders = 1
		required_flags = [1,4]
	}`, name, username)
}
