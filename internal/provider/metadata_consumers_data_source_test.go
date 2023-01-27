package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataConsumersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccMetadataResourceConfig("datasourceTest", "false"),
			},
			// Read testing
			{
				Config: testAccMetadataConsumersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.whisparr_metadata_consumers.test", "metadata_consumers.*", map[string]string{"movie_metadata": "false"}),
				),
			},
		},
	})
}

const testAccMetadataConsumersDataSourceConfig = `
data "whisparr_metadata_consumers" "test" {
}
`
