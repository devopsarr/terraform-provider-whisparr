resource "whisparr_import_list_stevenlu" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  should_monitor       = true
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  link                 = "https://s3.amazonaws.com/popular-movies/movies.json"
}