package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatConditionReleaseTitleDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionReleaseTitleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_custom_format_condition_release_title.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_custom_format_condition_release_title.test", "name", "x265"),
					resource.TestCheckResourceAttr("whisparr_custom_format.test", "specifications.0.value", "(((x|h)\\.?265)|(HEVC))")),
			},
		},
	})
}

const testAccCustomFormatConditionReleaseTitleDataSourceConfig = `
data  "whisparr_custom_format_condition_release_title" "test" {
	name = "x265"
	negate = false
	required = false
	value = "(((x|h)\\.?265)|(HEVC))"
}

resource "whisparr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSReleaseTitle"
	
	specifications = [data.whisparr_custom_format_condition_release_title.test]	
}`
