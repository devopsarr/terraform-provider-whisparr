package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationGotifyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationGotifyResourceConfig("error", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationGotifyResourceConfig("resourceGotifyTest", 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_gotify.test", "priority", "0"),
					resource.TestCheckResourceAttrSet("whisparr_notification_gotify.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationGotifyResourceConfig("error", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationGotifyResourceConfig("resourceGotifyTest", 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_gotify.test", "priority", "5"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_notification_gotify.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationGotifyResourceConfig(name string, priority int) string {
	return fmt.Sprintf(`
	resource "whisparr_notification_gotify" "test" {
		on_grab                            = false
		on_download                        = false
		on_upgrade                         = false
		on_movie_delete                    = false
		on_movie_file_delete               = false
		on_movie_file_delete_for_upgrade   = false
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		server = "http://gotify-server.net"
		app_token = "Token"
		priority = %d
	}`, name, priority)
}
