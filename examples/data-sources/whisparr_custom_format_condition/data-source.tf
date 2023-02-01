data "whisparr_custom_format_condition" "example" {
  name           = "Example"
  implementation = "SizeSpecification"
  negate         = false
  required       = false
  min            = 0
  max            = 100
}

resource "whisparr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.whisparr_custom_format_condition.example]
}