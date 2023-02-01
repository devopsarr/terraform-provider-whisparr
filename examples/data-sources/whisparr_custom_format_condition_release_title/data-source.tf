data "whisparr_custom_format_condition_release_title" "example" {
  name     = "x265"
  negate   = false
  required = false
  value    = "(((x|h)\\.?265)|(HEVC))"
}

resource "whisparr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.whisparr_custom_format_condition_release_title.example]
}