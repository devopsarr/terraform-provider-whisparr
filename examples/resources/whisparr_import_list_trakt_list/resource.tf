resource "whisparr_import_list_trakt_list" "example" {
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
  username             = "User2"
  listname             = "test"
  limit                = 100
}