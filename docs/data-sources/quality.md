---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_quality Data Source - terraform-provider-whisparr"
subcategory: "Profiles"
description: |-
  Single Quality.
---

# whisparr_quality (Data Source)

<!-- subcategory:Profiles -->Single Quality.

## Example Usage

```terraform
data "whisparr_quality" "bluray" {
  name = "Bluray-2160p"
}

data "whisparr_quality" "webdl" {
  name = "WEBDL-2160p"
}

data "whisparr_quality" "webrip" {
  name = "WEBRip-2160p"
}

resource "whisparr_quality_profile" "Example" {
  name            = "Example"
  upgrade_allowed = true
  cutoff          = 2000

  language = data.whisparr_language.test

  quality_groups = [
    {
      id   = 2000
      name = "WEB 2160p"
      qualities = [
        data.whisparr_quality.webdl,
        data.whisparr_quality.webrip,
      ]
    },
    {
      qualities = [data.whisparr_quality.bluray]
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Quality Name.

### Read-Only

- `id` (Number) Quality  ID.
- `resolution` (Number) Quality Resolution.
- `source` (String) Quality source.


