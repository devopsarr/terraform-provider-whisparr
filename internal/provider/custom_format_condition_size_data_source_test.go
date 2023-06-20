package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatConditionSizeDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionSizeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_custom_format_condition_size.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_custom_format_condition_size.test", "name", "Test"),
					resource.TestCheckResourceAttr("whisparr_custom_format.test", "specifications.0.min", "5"),
					resource.TestCheckResourceAttr("whisparr_custom_format.test", "specifications.0.max", "50")),
			},
		},
	})
}

const testAccCustomFormatConditionSizeDataSourceConfig = `
data  "whisparr_custom_format_condition_size" "test" {
	name = "Test"
	negate = false
	required = false
	min = 5
	max = 50
}

resource "whisparr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSSize"
	
	specifications = [data.whisparr_custom_format_condition_size.test]	
}`
