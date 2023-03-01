package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccQualityDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccQualityDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find quality"),
			},
			// Read testing
			{
				Config: testAccQualityDataSourceConfig("Remux-2160p"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_quality.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_quality.test", "resolution", "2160")),
			},
		},
	})
}

func testAccQualityDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "whisparr_quality" "test" {
		name = "%s"
	}
	`, name)
}
