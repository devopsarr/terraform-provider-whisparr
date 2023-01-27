---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_import_list_trakt_popular Resource - terraform-provider-whisparr"
subcategory: "Import Lists"
description: |-
  Import List Trakt Popular resource.
  For more information refer to Import List https://wiki.servarr.com/whisparr/settings#import-lists and Trakt Popular https://wiki.servarr.com/whisparr/supported#traktpopularimport.
---

# whisparr_import_list_trakt_popular (Resource)

<!-- subcategory:Import Lists -->Import List Trakt Popular resource.
For more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#import-lists) and [Trakt Popular](https://wiki.servarr.com/whisparr/supported#traktpopularimport).

## Example Usage

```terraform
resource "whisparr_import_list_trakt_popular" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  should_monitor       = true
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  auth_user            = "User1"
  access_token         = "Token"
  trakt_list_type      = 0
  limit                = 100
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `access_token` (String, Sensitive) Access token.
- `auth_user` (String) Auth user.
- `limit` (Number) limit.
- `minimum_availability` (String) Minimum availability.
- `name` (String) Import List name.
- `quality_profile_id` (Number) Quality profile ID.
- `root_folder_path` (String) Root folder path.
- `should_monitor` (Boolean) Should monitor.

### Optional

- `certification` (String) Certification.
- `enable_auto` (Boolean) Enable automatic add flag.
- `enabled` (Boolean) Enabled flag.
- `expires` (String) Expires.
- `genres` (String) Genres.
- `list_order` (Number) List order.
- `rating` (String) Rating.
- `refresh_token` (String, Sensitive) Refresh token.
- `search_on_add` (Boolean) Search on add flag.
- `tags` (Set of Number) List of associated tags.
- `trakt_additional_parameters` (String) Trakt additional parameters.
- `trakt_list_type` (Number) Trakt list type.`0` Trending, `1` Popular, `2` Anticipated, `3` BoxOffice, `4` TopWatchedByWeek, `5` TopWatchedByMonth, `6` TopWatchedByYear, `7` TopWatchedByAllTime, `8` RecommendedByWeek, `9` RecommendedByMonth, `10` RecommendedByYear, `10` RecommendedByAllTime.
- `years` (String) Years.

### Read-Only

- `id` (Number) Import List ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import whisparr_import_list_trakt_popular.example 1
```