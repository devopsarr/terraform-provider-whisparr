package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityProfileResource(t *testing.T) {
	// no parallel to avoid conflict with custom formats
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-4k"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_quality_profile.test", "name", "example-4k"),
					resource.TestCheckResourceAttrSet("whisparr_quality_profile.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-HD"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("whisparr_quality_profile.test", "name", "example-HD"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "whisparr_quality_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccQualityProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "whisparr_custom_format" "test" {
		include_custom_format_when_renaming = false
		name = "QualityFormatTest"
		
		specifications = [
			{
				name = "Surround Sound"
				implementation = "ReleaseTitleSpecification"
				negate = false
				required = false
				value = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
			},
			{
				name = "Arabic"
				implementation = "LanguageSpecification"
				negate = false
				required = false
				value = "31"
			},
			{
				name = "Size"
				implementation = "SizeSpecification"
				negate = false
				required = false
				min = 0
				max = 100
			}
		]	
	}

	data "whisparr_custom_formats" "test" {
		depends_on = [whisparr_custom_format.test]
	}

	data "whisparr_language" "test" {
		name = "English"
	}

	data "whisparr_quality" "bluray" {
		name = "Bluray-2160p"
	}

	data "whisparr_quality" "webdl" {
		name = "WEBDL-2160p"
	}

	data "whisparr_quality" "webrip" {
		name = "WEBRip-2160p"
	}

	resource "whisparr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 2000

		language = data.whisparr_language.test

		quality_groups = [
			{
				id   = 2000
				name = "WEB 2160p"
				qualities = [
					data.whisparr_quality.webdl,
					data.whisparr_quality.webrip,
				]
			},
			{
				qualities = [data.whisparr_quality.bluray]
			}
		]

		format_items = [
			for format in data.whisparr_custom_formats.test.custom_formats :
			{
				name   = format.name
				format = format.id
				score  = 0
			}
		]
	}`, name)
}
