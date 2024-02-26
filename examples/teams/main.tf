variable "teams" {
  type = map(any)
  default = {
    developers = ["user1", "user3"]
    managers   = ["user3", "user2"]
  }
}

module "teams" {
  source = "../.modules/teams"
  teams  = var.teams
}
