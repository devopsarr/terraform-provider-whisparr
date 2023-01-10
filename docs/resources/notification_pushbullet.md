---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_notification_pushbullet Resource - terraform-provider-whisparr"
subcategory: "Notifications"
description: |-
  Notification Pushbullet resource.
  For more information refer to Notification https://wiki.servarr.com/whisparr/settings#connect and Pushbullet https://wiki.servarr.com/whisparr/supported#pushbullet.
---

# whisparr_notification_pushbullet (Resource)

<!-- subcategory:Notifications -->Notification Pushbullet resource.
For more information refer to [Notification](https://wiki.servarr.com/whisparr/settings#connect) and [Pushbullet](https://wiki.servarr.com/whisparr/supported#pushbullet).



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) API key.
- `name` (String) NotificationPushbullet name.
- `on_movie_delete` (Boolean) On movie delete flag.

### Optional

- `channel_tags` (Set of String) List of channel tags.
- `device_ids` (Set of String) List of devices IDs.
- `include_health_warnings` (Boolean) Include health warnings.
- `on_application_update` (Boolean) On application update flag.
- `on_download` (Boolean) On download flag.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_movie_file_delete` (Boolean) On movie file delete flag.
- `on_movie_file_delete_for_upgrade` (Boolean) On movie file delete for upgrade flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `sender_id` (String) Sender ID.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

