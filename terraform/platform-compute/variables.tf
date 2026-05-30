variable "cluster_name" {
  default = "doclane"
}

variable "region" {
  default = "eu-west-1"
}

variable "az1" {
  default = "eu-west-1a"
}

variable "az2" {
  default = "eu-west-1b"
}

variable "db_name" {
  default = "doclane"
}

variable "db_username" {
  default = "doclane"
}

variable "db_password" {
  type      = string
  sensitive = true
  default   = "DoclanePass2026!"
}
