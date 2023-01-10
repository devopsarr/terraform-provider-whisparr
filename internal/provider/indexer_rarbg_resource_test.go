package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerRarbgResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerRarbgResourceConfig("rarbgResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_rarbg.test", "ranked_only", "false"),
					resource.TestCheckResourceAttr("whisparr_indexer_rarbg.test", "base_url", "https://torrentapi.org"),
					resource.TestCheckResourceAttrSet("whisparr_indexer_rarbg.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerRarbgResourceConfig("rarbgResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_rarbg.test", "ranked_only", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_indexer_rarbg.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerRarbgResourceConfig(name, ranked string) string {
	return fmt.Sprintf(`
	resource "whisparr_indexer_rarbg" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://torrentapi.org"
		ranked_only = "%s"
		minimum_seeders = 1
		categories = [46, 42]
	}`, name, ranked)
}
