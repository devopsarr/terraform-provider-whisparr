package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMetadataDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_metadata.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_metadata.test", "movie_metadata", "false")),
			},
		},
	})
}

const testAccMetadataDataSourceConfig = `
resource "whisparr_metadata" "test" {
	enable = true
	name = "metadataData"
	implementation = "MediaBrowserMetadata"
	config_contract = "MediaBrowserMetadataSettings"
	movie_metadata = false
}

data "whisparr_metadata" "test" {
	name = whisparr_metadata.test.name
}
`
