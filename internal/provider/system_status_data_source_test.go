package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSystemStatusDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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
