resource "whisparr_import_list_whisparr" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  should_monitor       = true
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  api_key              = "ExampleAPIKey"
  tag_ids              = [1, 2]
  profile_ids          = [1]
}