---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "whisparr_indexer_torrent_rss Resource - terraform-provider-whisparr"
subcategory: "Indexers"
description: |-
  Indexer Torrent RSS resource.
  For more information refer to Indexer https://wiki.servarr.com/whisparr/settings#indexers and Torrent RSS https://wiki.servarr.com/whisparr/supported#torrentrssindexer.
---

# whisparr_indexer_torrent_rss (Resource)

<!-- subcategory:Indexers -->Indexer Torrent RSS resource.
For more information refer to [Indexer](https://wiki.servarr.com/whisparr/settings#indexers) and [Torrent RSS](https://wiki.servarr.com/whisparr/supported#torrentrssindexer).



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `base_url` (String) Base URL.
- `name` (String) IndexerTorrentRss name.

### Optional

- `allow_zero_size` (Boolean) Allow zero size files.
- `cookie` (String) Cookie.
- `download_client_id` (Number) Download client ID.
- `enable_rss` (Boolean) Enable RSS flag.
- `minimum_seeders` (Number) Minimum seeders.
- `multi_languages` (Set of Number) Languages list.
- `priority` (Number) Priority.
- `required_flags` (Set of Number) Flag list.
- `seed_ratio` (Number) Seed ratio.
- `seed_time` (Number) Seed time.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) IndexerTorrentRss ID.

