package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemStatusDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			// {
			// 	Config:      testAccSystemStatusDataSourceConfig + testUnauthorizedProvider,
			// 	ExpectError: regexp.MustCompile("Client Error"),
			// },
			// Read testing
			{
				Config: testAccSystemStatusDataSourceConfig,
				// Check: resource.ComposeAggregateTestCheckFunc(
				// 	resource.TestCheckResourceAttrSet("data.whisparr_system_status.test", "id"),
				// 	resource.TestCheckResourceAttr("data.whisparr_system_status.test", "is_production", "true")),
			},
		},
	})
}

const testAccSystemStatusDataSourceConfig = `
data "whisparr_system_status" "test" {
}
`
