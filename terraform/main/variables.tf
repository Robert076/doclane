variable "public_subnet_list" {
  type        = list(string)
  description = "The list of public subnets"
  default     = ["10.0.0.0/24", "10.0.1.0/24"]
}

variable "private_subnet_list" {
  type        = list(string)
  description = "The list of private subnets"
  default     = ["10.0.2.0/24", "10.0.3.0/24"]
}

variable "region" {
  type        = string
  description = "Region used in the deployment"
  default     = "eu-west-1"
}

variable "project_name" {
  type        = string
  description = "The name of the project, used as a tag and name in various places"
  default     = "doclane"
}
