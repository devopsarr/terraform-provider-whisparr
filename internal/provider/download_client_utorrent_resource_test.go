package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientUtorrentResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientUtorrentResourceConfig("resourceUtorrentTest", "utorrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_utorrent.test", "host", "utorrent"),
					resource.TestCheckResourceAttr("whisparr_download_client_utorrent.test", "url_base", "/utorrent/"),
					resource.TestCheckResourceAttrSet("whisparr_download_client_utorrent.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientUtorrentResourceConfig("resourceUtorrentTest", "utorrent-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_utorrent.test", "host", "utorrent-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_download_client_utorrent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientUtorrentResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "whisparr_download_client_utorrent" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/utorrent/"
		port = 9091
		movie_category = "tv-whisparr"
	}`, name, host)
}
