data "codefresh_account" "acc" {
  name = var.account_name != ""? var.account_name : ""
  _id = var.account_id != ""? var.account_id : ""

}

resource "random_string" "random" {
  length = 16
  special = false
}

resource "codefresh_api_key" "new" {
  account_id = data.codefresh_account.acc.id
  user_id = data.codefresh_account.acc.admins[0]
  name = "tfkey_${random_string.random.result}"

  scopes = [
    "agent",
    "agents",
    "audit",
    "build",
    "cluster",
    "clusters",
    "environments-v2",
    "github-action",
    "helm",
    "kubernetes",
    "pipeline",
    "project",
    "repos",
    "runner-installation",
    "step-type",
    "step-types",
    "view",
    "workflow",
  ]
}

output "token" {
  value = codefresh_api_key.new.token
}

output "account_name" {
  value = data.codefresh_account.acc.name
}

output "account_id" {
  value = data.codefresh_account.acc.id
}