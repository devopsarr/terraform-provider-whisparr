resource "whisparr_notification_simplepush" "example" {
  on_grab                          = false
  on_download                      = true
  on_upgrade                       = true
  on_movie_delete                  = false
  on_movie_file_delete             = false
  on_movie_file_delete_for_upgrade = true
  on_health_issue                  = false
  on_application_update            = false

  include_health_warnings = false
  name                    = "Example"

  key   = "Token"
  event = "ringtone:default"
}