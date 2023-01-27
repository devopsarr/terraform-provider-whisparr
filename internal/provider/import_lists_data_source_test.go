package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListResourceConfig("importListsDataTest", "false"),
			},
			// Read testing
			{
				Config: testAccImportListsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.whisparr_import_lists.test", "import_lists.*", map[string]string{"base_url": "http://127.0.0.1:6969"}),
				),
			},
		},
	})
}

const testAccImportListsDataSourceConfig = `
data "whisparr_import_lists" "test" {
}
`
