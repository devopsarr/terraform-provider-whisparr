---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_metadata_roksbox Resource - terraform-provider-whisparr"
subcategory: "Metadata"
description: |-
  Metadata Roksbox resource.
  For more information refer to Metadata https://wiki.servarr.com/whisparr/settings#metadata and ROKSBOX https://wiki.servarr.com/whisparr/supported#roksboxmetadata.
---

# whisparr_metadata_roksbox (Resource)

<!-- subcategory:Metadata -->Metadata Roksbox resource.
For more information refer to [Metadata](https://wiki.servarr.com/whisparr/settings#metadata) and [ROKSBOX](https://wiki.servarr.com/whisparr/supported#roksboxmetadata).

## Example Usage

```terraform
resource "whisparr_metadata_roksbox" "example" {
  enable         = true
  name           = "Example"
  movie_metadata = true
  movie_images   = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `movie_images` (Boolean) Movie images flag.
- `movie_metadata` (Boolean) Movie metadata flag.
- `name` (String) Metadata name.

### Optional

- `enable` (Boolean) Enable flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Metadata ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import whisparr_metadata_roksbox.example 1
```