package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListExclusionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImportListExclusionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_import_list_exclusion.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_import_list_exclusion.test", "title", "testDS"),
				),
			},
		},
	})
}

const testAccImportListExclusionDataSourceConfig = `
resource "whisparr_import_list_exclusion" "test" {
	title = "testDS"
	tmdb_id = 987
	year = 1990
}

data "whisparr_import_list_exclusion" "test" {
	tmdb_id = whisparr_import_list_exclusion.test.tmdb_id
}
`
