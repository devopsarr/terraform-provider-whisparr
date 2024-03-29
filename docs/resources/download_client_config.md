---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_download_client_config Resource - terraform-provider-whisparr"
subcategory: "Download Clients"
description: |-
  Download Client Config resource.
  For more information refer to Download Client https://wiki.servarr.com/whisparr/settings#completed-download-handling documentation.
---

# whisparr_download_client_config (Resource)

<!-- subcategory:Download Clients -->Download Client Config resource.
For more information refer to [Download Client](https://wiki.servarr.com/whisparr/settings#completed-download-handling) documentation.

## Example Usage

```terraform
resource "whisparr_download_client_config" "example" {
  check_for_finished_download_interval = 1
  enable_completed_download_handling   = true
  auto_redownload_failed               = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auto_redownload_failed` (Boolean) Auto Redownload Failed flag.
- `check_for_finished_download_interval` (Number) Check for finished download interval.
- `enable_completed_download_handling` (Boolean) Enable Completed Download Handling flag.

### Read-Only

- `download_client_working_folders` (String) Download Client Working Folders.
- `id` (Number) Download Client Config ID.

## Import

Import is supported using the following syntax:

```shell
# import does not need parameters
terraform import whisparr_download_client_config.example ""
```
