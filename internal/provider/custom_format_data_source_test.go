package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccCustomFormatDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccCustomFormatDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find custom_format"),
			},
			// Read testing
			{
				Config: testAccCustomFormatResourceConfig("dataTest", "false") + testAccCustomFormatDataSourceConfig("whisparr_custom_format.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_custom_format.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_custom_format.test", "include_custom_format_when_renaming", "false")),
			},
		},
	})
}

func testAccCustomFormatDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "whisparr_custom_format" "test" {
		name = %s
	}
	`, name)
}
