resource "whisparr_notification_trakt" "example" {
  on_download                      = true
  on_upgrade                       = true
  on_movie_added                   = false
  on_movie_delete                  = false
  on_movie_file_delete             = false
  on_movie_file_delete_for_upgrade = true

  include_health_warnings = false
  name                    = "Example"

  auth_user    = "User"
  access_token = "AuthTOKEN"
}