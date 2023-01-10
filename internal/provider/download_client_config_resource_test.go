package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientConfigResourceConfig("true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_config.test", "auto_redownload_failed", "true"),
					resource.TestCheckResourceAttrSet("whisparr_download_client_config.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientConfigResourceConfig("false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_download_client_config.test", "auto_redownload_failed", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_download_client_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientConfigResourceConfig(redownload string) string {
	return fmt.Sprintf(`
	resource "whisparr_download_client_config" "test" {
		check_for_finished_download_interval = 1
		enable_completed_download_handling = true
		auto_redownload_failed = %s
	}`, redownload)
}
