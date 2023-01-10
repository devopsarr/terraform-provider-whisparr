package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerHdbitsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerHdbitsResourceConfig("hdbitsResourceTest", "user"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_hdbits.test", "username", "user"),
					resource.TestCheckResourceAttrSet("whisparr_indexer_hdbits.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerHdbitsResourceConfig("hdbitsResourceTest", "Username"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_hdbits.test", "username", "Username"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_indexer_hdbits.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerHdbitsResourceConfig(name, username string) string {
	return fmt.Sprintf(`
	resource "whisparr_indexer_hdbits" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://hdbits.org"
		username = "%s"
		api_key = "Key"
		minimum_seeders = 1
		categories = [1]
		codecs = [1,5]
	}`, name, username)
}
