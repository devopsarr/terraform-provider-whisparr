package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccIndexerDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccIndexerDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find indexer"),
			},
			// Read testing
			{
				Config: testAccIndexerResourceConfig("indexerdata", "30") + testAccIndexerDataSourceConfig("whisparr_indexer.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_indexer.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_indexer.test", "protocol", "usenet")),
			},
		},
	})
}

func testAccIndexerDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "whisparr_indexer" "test" {
		name = %s
	}
	`, name)
}
