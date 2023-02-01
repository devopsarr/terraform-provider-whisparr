data "whisparr_custom_format_condition_edition" "example" {
  name     = "Example"
  negate   = false
  required = false
  value    = ".*Extended.*"
}

resource "whisparr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.whisparr_custom_format_condition_edition.example]
}