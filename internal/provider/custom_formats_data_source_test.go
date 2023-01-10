package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccCustomFormatResourceConfig("datasourceTest", "true"),
			},
			// Read testing
			{
				Config: testAccCustomFormatsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.whisparr_custom_formats.test", "custom_formats.*", map[string]string{"include_custom_format_when_renaming": "true"}),
				),
			},
		},
	})
}

const testAccCustomFormatsDataSourceConfig = `
data "whisparr_custom_formats" "test" {
}
`
