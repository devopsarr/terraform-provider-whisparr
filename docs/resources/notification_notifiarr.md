---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_notification_notifiarr Resource - terraform-provider-whisparr"
subcategory: "Notifications"
description: |-
  Notification Notifiarr resource.
  For more information refer to Notification https://wiki.servarr.com/whisparr/settings#connect and Notifiarr https://wiki.servarr.com/whisparr/supported#notifiarr.
---

# whisparr_notification_notifiarr (Resource)

<!-- subcategory:Notifications -->Notification Notifiarr resource.
For more information refer to [Notification](https://wiki.servarr.com/whisparr/settings#connect) and [Notifiarr](https://wiki.servarr.com/whisparr/supported#notifiarr).

## Example Usage

```terraform
resource "whisparr_notification_notifiarr" "example" {
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

  api_key = "Token"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) API key.
- `name` (String) NotificationNotifiarr name.
- `on_movie_delete` (Boolean) On movie delete flag.

### Optional

- `include_health_warnings` (Boolean) Include health warnings.
- `on_application_update` (Boolean) On application update flag.
- `on_download` (Boolean) On download flag.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_movie_file_delete` (Boolean) On movie file delete flag.
- `on_movie_file_delete_for_upgrade` (Boolean) On movie file delete for upgrade flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import whisparr_notification_notifiarr.example 1
```
