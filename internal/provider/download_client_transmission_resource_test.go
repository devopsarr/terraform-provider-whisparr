package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientTransmissionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientTransmissionResourceConfig("resourceTransmissionTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_transmission.test", "enable", "false"),
					resource.TestCheckResourceAttr("whisparr_download_client_transmission.test", "url_base", "/transmission/"),
					resource.TestCheckResourceAttrSet("whisparr_download_client_transmission.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientTransmissionResourceConfig("resourceTransmissionTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_transmission.test", "enable", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_download_client_transmission.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientTransmissionResourceConfig(name, enable string) string {
	return fmt.Sprintf(`
	resource "whisparr_download_client_transmission" "test" {
		enable = %s
		priority = 1
		name = "%s"
		host = "transmission"
		url_base = "/transmission/"
		port = 9091
	}`, enable, name)
}
