variable "default_acccount_limits" {
  type = map(any)
  default = {
    collaborators   = 100
    parallel_builds = 10
  }
}

variable "default_idps" {
  type = map(any)
}

variable "accounts" {
  type = map(any)
}

variable "users" {
  //type = map(any)
}

