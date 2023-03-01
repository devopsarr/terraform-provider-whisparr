package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNamingResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNamingResourceConfig("spaceDash") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNamingResourceConfig("spaceDash"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_naming.test", "colon_replacement_format", "spaceDash"),
					resource.TestCheckResourceAttrSet("whisparr_naming.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNamingResourceConfig("spaceDash") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNamingResourceConfig("dash"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_naming.test", "colon_replacement_format", "dash"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_naming.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNamingResourceConfig(replace string) string {
	return fmt.Sprintf(`
	resource "whisparr_naming" "test" {
		include_quality = false
		rename_movies = true
		replace_illegal_characters = false
		replace_spaces = false
		colon_replacement_format =  "%s"
		standard_movie_format = "{Movie Title} ({Release Year}) {Quality Full}"
		movie_folder_format = "{Movie Title} ({Release Year})"
	}`, replace)
}
