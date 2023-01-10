package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRootFolderDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_root_folder.test", "path", "/tmp")),
			},
		},
	})
}

const testAccRootFolderDataSourceConfig = `
resource "whisparr_root_folder" "test" {
	path = "/tmp"
}

data "whisparr_root_folder" "test" {
	path = whisparr_root_folder.test.path
}
`
