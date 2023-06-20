package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerConfigResourceConfig(0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerConfigResourceConfig(20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_config.test", "rss_sync_interval", "20"),
					resource.TestCheckResourceAttrSet("whisparr_indexer_config.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerConfigResourceConfig(0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerConfigResourceConfig(30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_indexer_config.test", "rss_sync_interval", "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_indexer_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerConfigResourceConfig(rss int) string {
	return fmt.Sprintf(`
	resource "whisparr_indexer_config" "test" {
		maximum_size = 0
		minimum_age = 0
		retention = 0
		rss_sync_interval = %d
		availability_delay = 0
		whitelisted_hardcoded_subs = ""
		prefer_indexer_flags = false
		allow_hardcoded_subs = false
	}`, rss)
}
