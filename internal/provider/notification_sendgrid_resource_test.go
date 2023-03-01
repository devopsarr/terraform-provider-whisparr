package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSendgridResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSendgridResourceConfig("error", "test@sendgrid.com") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSendgridResourceConfig("resourceSendgridTest", "test@sendgrid.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_sendgrid.test", "from", "test@sendgrid.com"),
					resource.TestCheckResourceAttrSet("whisparr_notification_sendgrid.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSimplepushResourceConfig("error", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSendgridResourceConfig("resourceSendgridTest", "test123@sendgrid.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_notification_sendgrid.test", "from", "test123@sendgrid.com"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_notification_sendgrid.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSendgridResourceConfig(name, from string) string {
	return fmt.Sprintf(`
	resource "whisparr_notification_sendgrid" "test" {
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
		
		api_key = "APIkey"
		from = "%s"
		recipients = ["test@test.com", "test1@test.com"]
	}`, name, from)
}
