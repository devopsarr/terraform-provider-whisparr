resource "whisparr_import_list_trakt_popular" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  should_monitor       = true
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  auth_user            = "User1"
  access_token         = "Token"
  trakt_list_type      = 0
  limit                = 100
}