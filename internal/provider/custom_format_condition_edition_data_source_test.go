package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatConditionEditionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionEditionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_custom_format_condition_edition.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_custom_format_condition_edition.test", "name", "Extended"),
					resource.TestCheckResourceAttr("whisparr_custom_format.test", "specifications.0.value", ".*Extended.*")),
			},
		},
	})
}

const testAccCustomFormatConditionEditionDataSourceConfig = `
data  "whisparr_custom_format_condition_edition" "test" {
	name = "Extended"
	negate = false
	required = false
	value = ".*Extended.*"
}

resource "whisparr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSEdition"
	
	specifications = [data.whisparr_custom_format_condition_edition.test]	
}`
