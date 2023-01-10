---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_indexer Resource - terraform-provider-whisparr"
subcategory: "Indexers"
description: |-
  Generic Indexer resource. When possible use a specific resource instead.
  For more information refer to Indexer https://wiki.servarr.com/whisparr/settings#indexers documentation.
---

# whisparr_indexer (Resource)

<!-- subcategory:Indexers -->Generic Indexer resource. When possible use a specific resource instead.
For more information refer to [Indexer](https://wiki.servarr.com/whisparr/settings#indexers) documentation.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config_contract` (String) Indexer configuration template.
- `implementation` (String) Indexer implementation name.
- `name` (String) Indexer name.
- `protocol` (String) Protocol. Valid values are 'usenet' and 'torrent'.

### Optional

- `additional_parameters` (String) Additional parameters.
- `allow_zero_size` (Boolean) Allow zero size files.
- `api_key` (String) API key.
- `api_path` (String) API path.
- `api_user` (String) API User.
- `base_url` (String) Base URL.
- `captcha_token` (String) Captcha token.
- `categories` (Set of Number) Series list.
- `codecs` (Set of Number) Codecs.
- `cookie` (String) Cookie.
- `delay` (Number) Delay before grabbing.
- `download_client_id` (Number) Download client ID.
- `enable_automatic_search` (Boolean) Enable automatic search flag.
- `enable_interactive_search` (Boolean) Enable interactive search flag.
- `enable_rss` (Boolean) Enable RSS flag.
- `mediums` (Set of Number) Mediumd.
- `minimum_seeders` (Number) Minimum seeders.
- `multi_languages` (Set of Number) Language list.
- `passkey` (String) Passkey.
- `priority` (Number) Priority.
- `ranked_only` (Boolean) Allow ranked only.
- `required_flags` (Set of Number) Required flags.
- `seed_ratio` (Number) Seed ratio.
- `seed_time` (Number) Seed time.
- `tags` (Set of Number) List of associated tags.
- `user` (String) Username.
- `username` (String) Username.

### Read-Only

- `id` (Number) Indexer ID.

