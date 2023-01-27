resource "whisparr_movie" "example" {
  monitored            = false
  title                = "Blue Movie"
  path                 = "/movies/Blue_Movie_1969"
  quality_profile_id   = 1
  tmdb_id              = 242423
  minimum_availability = "inCinemas"
}