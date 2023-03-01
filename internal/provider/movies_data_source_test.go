package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMoviesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMovieResourceConfig("Error", "error", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccMovieResourceConfig("Kim Kardashian, Superstar", "Kim_Kardashian_Superstar_2007", 45323) + testAccMoviesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.whisparr_movies.test", "movies.*", map[string]string{"title": "Kim Kardashian, Superstar"}),
				),
			},
		},
	})
}

const testAccMoviesDataSourceConfig = `
data "whisparr_movies" "test" {
	depends_on = [whisparr_movie.test]
}
`
