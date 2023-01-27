package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_import_list.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_import_list.test", "should_monitor", "true")),
			},
		},
	})
}

const testAccImportListDataSourceConfig = `
resource "whisparr_import_list" "test" {
	enabled = false
	enable_auto = false
	search_on_add = false
	list_type = "program"
	root_folder_path = "/config"
	should_monitor = true
	minimum_availability = "tba"
	quality_profile_id = 1
	name = "importListDataTest"
	implementation = "WhisparrImport"
	config_contract = "WhisparrSettings"
	base_url = "http://127.0.0.1:6969"
	api_key = "testAPIKey"
}

data "whisparr_import_list" "test" {
	name = whisparr_import_list.test.name
}
`
