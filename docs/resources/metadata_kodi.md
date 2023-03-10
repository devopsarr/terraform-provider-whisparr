---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_metadata_kodi Resource - terraform-provider-whisparr"
subcategory: "Metadata"
description: |-
  Metadata Kodi resource.
  For more information refer to Metadata https://wiki.servarr.com/whisparr/settings#metadata and KODI https://wiki.servarr.com/whisparr/supported#xbmcmetadata.
---

# whisparr_metadata_kodi (Resource)

<!-- subcategory:Metadata -->Metadata Kodi resource.
For more information refer to [Metadata](https://wiki.servarr.com/whisparr/settings#metadata) and [KODI](https://wiki.servarr.com/whisparr/supported#xbmcmetadata).

## Example Usage

```terraform
resource "whisparr_metadata_kodi" "example" {
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
- `movie_metadata_language` (Number) Movie metadata language.
- `movie_metadata_url` (Boolean) Movie metadata URL flag.
- `name` (String) Metadata name.
- `use_movie_nfo` (Boolean) Use movie nfo flag.

### Optional

- `enable` (Boolean) Enable flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Metadata ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import whisparr_metadata_kodi.example 1
```
