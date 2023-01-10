package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRestrictionsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a resource to check
			{
				Config: testAccRestrictionResourceConfig("testDataSource", "testDataSource2"),
			},
			// Read testing
			{
				Config: testAccRestrictionsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.whisparr_restrictions.test", "restrictions.*", map[string]string{"ignored": "testDataSource"}),
				),
			},
		},
	})
}

const testAccRestrictionsDataSourceConfig = `
data "whisparr_restrictions" "test" {
}
`
