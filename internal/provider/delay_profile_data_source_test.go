package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDelayProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDelayProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_delay_profile.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_delay_profile.test", "enable_usenet", "true")),
			},
		},
	})
}

const testAccDelayProfileDataSourceConfig = `
data "whisparr_delay_profile" "test" {
	id = 1
}
`
