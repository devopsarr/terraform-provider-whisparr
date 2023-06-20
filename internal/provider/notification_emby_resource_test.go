package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationEmbyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationEmbyResourceConfig("error", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationEmbyResourceConfig("resourceEmbyTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_emby.test", "api_key", "token123"),
					resource.TestCheckResourceAttrSet("whisparr_notification_emby.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationEmbyResourceConfig("error", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationEmbyResourceConfig("resourceEmbyTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_emby.test", "api_key", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_notification_emby.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationEmbyResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "whisparr_notification_emby" "test" {
		on_grab                            = false
		on_download                        = false
		on_upgrade                         = false
		on_rename                          = false
		on_movie_delete                    = false
		on_movie_file_delete               = false
		on_movie_file_delete_for_upgrade   = false
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		host = "emby.lcl"
		port = 8096
		api_key = "%s"
	}`, name, token)
}
