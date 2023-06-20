package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatConditionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_custom_format_condition.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_custom_format_condition.test", "name", "Surround Sound"),
					resource.TestCheckResourceAttr("whisparr_custom_format.test", "specifications.0.implementation", "ReleaseTitleSpecification")),
			},
		},
	})
}

const testAccCustomFormatConditionDataSourceConfig = `
data  "whisparr_custom_format_condition" "test" {
	name = "Surround Sound"
	implementation = "ReleaseTitleSpecification"
	negate = false
	required = false
	value = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
}

data  "whisparr_custom_format_condition" "test1" {
	name = "Size"
	implementation = "SizeSpecification"
	negate = false
	required = false
	min = 0
	max = 100
}

resource "whisparr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDS"
	
	specifications = [data.whisparr_custom_format_condition.test,data.whisparr_custom_format_condition.test1]	
}`
