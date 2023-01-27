package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMoviesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMoviesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.whisparr_movies.test", "movies.*", map[string]string{"title": "Kim Kardashian, Superstar"}),
				),
			},
		},
	})
}

const testAccMoviesDataSourceConfig = `
resource "whisparr_movie" "test" {
	monitored = false
	title = "Kim Kardashian, Superstar"
	path = "/config/Kim_Kardashian_Superstar_2007"
	quality_profile_id = 1
	tmdb_id = 45323

	minimum_availability = "inCinemas"
}

data "whisparr_movies" "test" {
	depends_on = [whisparr_movie.test]
}
`
