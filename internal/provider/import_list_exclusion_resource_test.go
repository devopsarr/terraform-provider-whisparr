package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListExclusionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", 123),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_exclusion.test", "tmdb_id", "123"),
					resource.TestCheckResourceAttrSet("whisparr_import_list_exclusion.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", 1234),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_import_list_exclusion.test", "tmdb_id", "1234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_import_list_exclusion.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListExclusionResourceConfig(name string, tvID int) string {
	return fmt.Sprintf(`
		resource "whisparr_import_list_exclusion" "%s" {
  			title = "Test"
			tmdb_id = %d
			year = 1900
		}
	`, name, tvID)
}
