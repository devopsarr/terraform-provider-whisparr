---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_import_list_stevenlu Resource - terraform-provider-whisparr"
subcategory: "Import Lists"
description: |-
  Import List Stevenlu resource.
  For more information refer to Import List https://wiki.servarr.com/whisparr/settings#import-lists and Stevenlu https://wiki.servarr.com/whisparr/supported#stevenluimport.
---

# whisparr_import_list_stevenlu (Resource)

<!-- subcategory:Import Lists -->Import List Stevenlu resource.
For more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#import-lists) and [Stevenlu](https://wiki.servarr.com/whisparr/supported#stevenluimport).

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `link` (String) Link.
- `minimum_availability` (String) Minimum availability.
- `name` (String) Import List name.
- `quality_profile_id` (Number) Quality profile ID.
- `root_folder_path` (String) Root folder path.
- `should_monitor` (Boolean) Should monitor.

### Optional

- `enable_auto` (Boolean) Enable automatic add flag.
- `enabled` (Boolean) Enabled flag.
- `list_order` (Number) List order.
- `search_on_add` (Boolean) Search on add flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Import List ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import whisparr_import_list_stevenlu.example 1
```
