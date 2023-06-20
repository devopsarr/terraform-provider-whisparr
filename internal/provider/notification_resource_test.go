package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification.test", "on_upgrade", "false"),
					resource.TestCheckResourceAttrSet("whisparr_notification.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification.test", "on_upgrade", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_notification.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "whisparr_notification" "test" {
		on_grab                            = false
		on_download                        = true
		on_upgrade                         = %s
		on_rename                          = false
		on_movie_delete                    = false
		on_movie_file_delete               = false
		on_movie_file_delete_for_upgrade   = true
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		implementation  = "CustomScript"
		config_contract = "CustomScriptSettings"
	  
		path = "/scripts/test.sh"
	}`, upgrade, name)
}
