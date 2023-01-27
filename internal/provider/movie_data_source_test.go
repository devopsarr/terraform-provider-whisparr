package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMovieDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMovieDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_movie.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_movie.test", "title", "Blue Movie"),
				),
			},
		},
	})
}

const testAccMovieDataSourceConfig = `
resource "whisparr_movie" "test" {
	monitored = false
	title = "Blue Movie"
	path = "/config/Blue_Movie_1969"
	quality_profile_id = 1
	tmdb_id = 242423
}

data "whisparr_movie" "test" {
	tmdb_id = whisparr_movie.test.tmdb_id
}
`
