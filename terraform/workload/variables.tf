variable "db_username" {
  default = "doclane"
}

variable "db_password" {
  type      = string
  sensitive = true
  default   = "DoclanePass2026!"
}

variable "seed_secret" {
  type      = string
  sensitive = true
  default   = "G+/3d4cnZDUcQhI3EB5lQTmNJbLpGiOdF5DqRdcXz5k="
}
