package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatConditionIndexerFlagDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionIndexerFlagDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_custom_format_condition_indexer_flag.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_custom_format_condition_indexer_flag.test", "name", "PTPGolden"),
					resource.TestCheckResourceAttr("whisparr_custom_format.test", "specifications.0.value", "8")),
			},
		},
	})
}

const testAccCustomFormatConditionIndexerFlagDataSourceConfig = `
data  "whisparr_custom_format_condition_indexer_flag" "test" {
	name = "PTPGolden"
	negate = false
	required = false
	value = "8"
}

resource "whisparr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSIndexerFlag"
	
	specifications = [data.whisparr_custom_format_condition_indexer_flag.test]	
}`
