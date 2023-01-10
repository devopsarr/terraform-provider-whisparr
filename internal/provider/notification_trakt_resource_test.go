package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationTraktResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationTraktResourceConfig("resourceTraktTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_trakt.test", "access_token", "token123"),
					resource.TestCheckResourceAttrSet("whisparr_notification_trakt.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationTraktResourceConfig("resourceTraktTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_trakt.test", "access_token", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_notification_trakt.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationTraktResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "whisparr_notification_trakt" "test" {
		on_download                        = false
		on_upgrade                         = false
		on_movie_delete                    = false
		on_movie_file_delete               = false
		on_movie_file_delete_for_upgrade   = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		auth_user = "User"
		access_token = "%s"
	}`, name, token)
}
