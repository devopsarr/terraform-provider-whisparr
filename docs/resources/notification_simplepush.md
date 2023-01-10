---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_notification_simplepush Resource - terraform-provider-whisparr"
subcategory: "Notifications"
description: |-
  Notification Simplepush resource.
  For more information refer to Notification https://wiki.servarr.com/whisparr/settings#connect and Simplepush https://wiki.servarr.com/whisparr/supported#simplepush.
---

# whisparr_notification_simplepush (Resource)

<!-- subcategory:Notifications -->Notification Simplepush resource.
For more information refer to [Notification](https://wiki.servarr.com/whisparr/settings#connect) and [Simplepush](https://wiki.servarr.com/whisparr/supported#simplepush).



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String, Sensitive) Key.
- `name` (String) NotificationSimplepush name.
- `on_movie_delete` (Boolean) On movie delete flag.

### Optional

- `event` (String) Event.
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

