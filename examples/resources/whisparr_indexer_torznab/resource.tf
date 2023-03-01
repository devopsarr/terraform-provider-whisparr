resource "whisparr_indexer_torznab" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://feed.animetosho.org"
  api_path                = "/nabapi"
  categories              = [2000, 2010]
  minimum_seeders         = 1
}
