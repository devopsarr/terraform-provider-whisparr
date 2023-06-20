package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRemotePathMappingDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccRemotePathMappingDataSourceConfig("999") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccRemotePathMappingDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find remote_path_mapping"),
			},
			// Read testing
			{
				Config: testAccRemotePathMappingResourceConfig("dataTest", "/test2/") + testAccRemotePathMappingDataSourceConfig("whisparr_remote_path_mapping.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_remote_path_mapping.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_remote_path_mapping.test", "host", "dataTest")),
			},
		},
	})
}

func testAccRemotePathMappingDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	data "whisparr_remote_path_mapping" "test" {
		id = %s
	}
	`, id)
}
