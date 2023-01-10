package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRestrictionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRestrictionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_restriction.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_restriction.test", "ignored", "datatest1")),
			},
		},
	})
}

const testAccRestrictionDataSourceConfig = `
resource "whisparr_restriction" "test" {
	ignored = "datatest1"
    required = "datatest2"
}

data "whisparr_restriction" "test" {
	id = whisparr_restriction.test.id
}
`
