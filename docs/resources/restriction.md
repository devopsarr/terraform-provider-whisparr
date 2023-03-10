---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_restriction Resource - terraform-provider-whisparr"
subcategory: "Indexers"
description: |-
  Restriction resource.
  For more information refer to Restriction https://wiki.servarr.com/whisparr/settings#remote-path-restrictions documentation.
---

# whisparr_restriction (Resource)

<!-- subcategory:Indexers -->Restriction resource.
For more information refer to [Restriction](https://wiki.servarr.com/whisparr/settings#remote-path-restrictions) documentation.

## Example Usage

```terraform
resource "whisparr_restriction" "example" {
  ignored  = "string1"
  required = "string2"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ignored` (String) Ignored. Either one of 'required' or 'ignored' must be set.
- `required` (String) Required. Either one of 'required' or 'ignored' must be set.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Restriction ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import whisparr_restriction.example 10
```
