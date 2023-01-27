---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_import_list_whisparr Resource - terraform-provider-whisparr"
subcategory: "Import Lists"
description: |-
  Import List Whisparr resource.
  For more information refer to Import List https://wiki.servarr.com/whisparr/settings#import-lists and Whisparr https://wiki.servarr.com/whisparr/supported#whisparrimport.
---

# whisparr_import_list_whisparr (Resource)

<!-- subcategory:Import Lists -->Import List Whisparr resource.
For more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#import-lists) and [Whisparr](https://wiki.servarr.com/whisparr/supported#whisparrimport).

## Example Usage

```terraform
resource "whisparr_import_list_whisparr" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  should_monitor       = true
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  api_key              = "ExampleAPIKey"
  tag_ids              = [1, 2]
  profile_ids          = [1]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) API key.
- `base_url` (String) Base URL.
- `minimum_availability` (String) Minimum availability.
- `name` (String) Import List name.
- `quality_profile_id` (Number) Quality profile ID.
- `root_folder_path` (String) Root folder path.
- `should_monitor` (Boolean) Should monitor.

### Optional

- `enable_auto` (Boolean) Enable automatic add flag.
- `enabled` (Boolean) Enabled flag.
- `list_order` (Number) List order.
- `profile_ids` (Set of Number) Profile IDs.
- `search_on_add` (Boolean) Search on add flag.
- `tag_ids` (Set of Number) Tag IDs.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Import List ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import whisparr_import_list_whisparr.example 1
```