package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMovieDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMovieDataSourceConfig("999") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccMovieDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find movie"),
			},
			// Read testing
			{
				Config: testAccMovieResourceConfig("Blue Movie", "Blue_Movie_1969", 242423) + testAccMovieDataSourceConfig("whisparr_movie.test.tmdb_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_movie.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_movie.test", "title", "Blue Movie"),
				),
			},
		},
	})
}

func testAccMovieDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	data "whisparr_movie" "test" {
		tmdb_id = %s
	}
	`, id)
}
