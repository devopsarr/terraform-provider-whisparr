package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccImportListConfigResourceConfig("removeAndDelete"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_config.test", "sync_level", "removeAndDelete"),
					resource.TestCheckResourceAttrSet("whisparr_import_list_config.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListConfigResourceConfig("logOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_config.test", "sync_level", "logOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_import_list_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListConfigResourceConfig(level string) string {
	return fmt.Sprintf(`
	resource "whisparr_import_list_config" "test" {
		sync_interval = 24
		sync_level = "%s"
	}`, level)
}
