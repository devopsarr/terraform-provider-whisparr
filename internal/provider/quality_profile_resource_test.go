package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-4k"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_quality_profile.test", "name", "example-4k"),
					resource.TestCheckResourceAttrSet("whisparr_quality_profile.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-HD"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_quality_profile.test", "name", "example-HD"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_quality_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccQualityProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "whisparr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 1003

		language = {
			id   = 1
			name = "English"
		}

		quality_groups = [
			{
				id   = 1003
				name = "WEB 2160p"
				qualities = [
					{
						id         = 18
						name       = "WEBDL-2160p"
						source     = "webdl"
						resolution = 2160
					},
					{
						id         = 17
						name       = "WEBRip-2160p"
						source     = "webrip"
						resolution = 2160
					}
				]
			}
		]
	}
	`, name)
}