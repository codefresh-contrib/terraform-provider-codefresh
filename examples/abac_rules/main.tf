data "codefresh_team" "users" {
  name = "users"
}

resource "codefresh_abac_rules" "app_rule" {
  entity_type = "gitopsApplications"
  teams       = [data.codefresh_team.users.id]
  actions     = ["REFRESH", "SYNC", "TERMINATE_SYNC", "VIEW_POD_LOGS", "APP_ROLLBACK"]

  attribute {
    name  = "LABEL"
    key   = "KEY"
    value = "VALUE"
  }

  tags = ["dev", "untagged"]
}
