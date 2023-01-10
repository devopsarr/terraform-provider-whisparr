package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDelayProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDelayProfileResourceConfig("usenet"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_delay_profile.test", "preferred_protocol", "usenet"),
					resource.TestCheckResourceAttrSet("whisparr_delay_profile.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDelayProfileResourceConfig("torrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_delay_profile.test", "preferred_protocol", "torrent"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_delay_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDelayProfileResourceConfig(protocol string) string {
	return fmt.Sprintf(`
	resource "whisparr_tag" "test" {
		label = "delay_profile_resource"
	}

	resource "whisparr_delay_profile" "test" {
		enable_usenet = true
		enable_torrent = true
		bypass_if_highest_quality = true
		order = 100
		usenet_delay = 0
		torrent_delay = 0
		preferred_protocol= "%s"
		tags = [whisparr_tag.test.id]
	}`, protocol)
}
