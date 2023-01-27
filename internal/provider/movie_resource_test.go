package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMovieResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccMovieResourceConfig("test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_movie.test", "path", "/config/test"),
					resource.TestCheckResourceAttrSet("whisparr_movie.test", "id"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "original_title", "Deep Throat"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "status", "released"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "monitored", "false"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "year", "1972"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "minimum_availability", "inCinemas"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "imdb_id", "tt0068468"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "is_available", "true"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "original_language.id", "1"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "original_language.name", "English"),
					resource.TestCheckResourceAttr("whisparr_movie.test", "genres.0", "Comedy"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMovieResourceConfig("test123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_movie.test", "path", "/config/test123"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_movie.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMovieResourceConfig(path string) string {
	return fmt.Sprintf(`
		resource "whisparr_movie" "test" {
			monitored = false
			title = "Deep Throat"
			path = "/config/%s"
			quality_profile_id = 1
			tmdb_id = 5853

			minimum_availability = "inCinemas"
		}
	`, path)
}
