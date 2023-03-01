package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSynologyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSynologyResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSynologyResourceConfig("resourceSynologyTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_synology_indexer.test", "update_library", "false"),
					resource.TestCheckResourceAttrSet("whisparr_notification_synology_indexer.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSynologyResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSynologyResourceConfig("resourceSynologyTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_synology_indexer.test", "update_library", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_notification_synology_indexer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSynologyResourceConfig(name, update string) string {
	return fmt.Sprintf(`
	resource "whisparr_notification_synology_indexer" "test" {
		on_download                        = false
		on_upgrade                         = false
		on_rename                          = false
		on_movie_delete                    = false
		on_movie_file_delete               = false
		on_movie_file_delete_for_upgrade   = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		update_library = %s
	}`, name, update)
}
