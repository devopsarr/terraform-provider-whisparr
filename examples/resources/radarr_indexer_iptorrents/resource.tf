resource "whisparr_indexer_iptorrents" "example" {
  enable_rss      = true
  name            = "Example"
  base_url        = "https://iptorrent.io"
  minimum_seeders = 1
}
