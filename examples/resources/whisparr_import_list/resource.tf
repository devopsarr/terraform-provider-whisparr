resource "whisparr_import_list" "example" {
  enabled              = false
  enable_auto          = true
  search_on_add        = false
  should_monitor       = true
  minimum_availability = "tba"
  list_type            = "program"
  root_folder_path     = whisparr_root_folder.example.path
  quality_profile_id   = whisparr_quality_profile.example.id
  name                 = "Example"
  implementation       = "WhisparrImport"
  config_contract      = "WhisparrSettings"
  tags                 = [1, 2]

  tag_ids     = [1, 2]
  profile_ids = [1]
  base_url    = "http://127.0.0.1:8686"
  api_key     = "APIKey"
}