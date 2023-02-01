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