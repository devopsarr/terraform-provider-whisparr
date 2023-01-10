package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientQbittorrentResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientQbittorrentResourceConfig("resourceQbittorrentTest", "qbittorrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_qbittorrent.test", "host", "qbittorrent"),
					resource.TestCheckResourceAttr("whisparr_download_client_qbittorrent.test", "url_base", "/qbittorrent/"),
					resource.TestCheckResourceAttrSet("whisparr_download_client_qbittorrent.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientQbittorrentResourceConfig("resourceQbittorrentTest", "qbittorrent-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_qbittorrent.test", "host", "qbittorrent-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_download_client_qbittorrent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientQbittorrentResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "whisparr_download_client_qbittorrent" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/qbittorrent/"
		port = 9091
		movie_category = "tv-whisparr"
		first_and_last = true
	}`, name, host)
}
