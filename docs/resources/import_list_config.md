---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_import_list_config Resource - terraform-provider-whisparr"
subcategory: "Import Lists"
description: |-
  Import List Config resource.
  For more information refer to Import List https://wiki.servarr.com/whisparr/settings#completed-download-handling documentation.
---

# whisparr_import_list_config (Resource)

<!-- subcategory:Import Lists -->Import List Config resource.
For more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#completed-download-handling) documentation.

## Example Usage

```terraform
resource "whisparr_import_list_config" "example" {
  sync_interval = 24
  sync_level    = "logOnly"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `sync_interval` (Number) List Update Interval.
- `sync_level` (String) Clean library level.

### Read-Only

- `id` (Number) Import List Config ID.

## Import

Import is supported using the following syntax:

```shell
# import does not need parameters
terraform import whisparr_import_list_config.example
```