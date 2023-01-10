resource "whisparr_notification_join" "example" {
  on_grab                          = false
  on_download                      = false
  on_upgrade                       = false
  on_movie_added                   = false
  on_movie_delete                  = false
  on_movie_file_delete             = false
  on_movie_file_delete_for_upgrade = true
  on_health_issue                  = false
  on_application_update            = false

  include_health_warnings = false
  name                    = "Example"

  api_key  = "Key"
  priority = 2
}