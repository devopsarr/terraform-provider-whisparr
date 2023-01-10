package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationNotifiarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationNotifiarrResourceConfig("resourceNotifiarrTest", "key1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_notifiarr.test", "api_key", "key1"),
					resource.TestCheckResourceAttrSet("whisparr_notification_notifiarr.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationNotifiarrResourceConfig("resourceNotifiarrTest", "key2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_notifiarr.test", "api_key", "key2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_notification_notifiarr.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationNotifiarrResourceConfig(name, key string) string {
	return fmt.Sprintf(`
	resource "whisparr_notification_notifiarr" "test" {
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
	  
		api_key = "%s"
	}`, name, key)
}
