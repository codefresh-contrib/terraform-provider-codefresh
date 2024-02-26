data "codefresh_account" "acc" {
  for_each = var.accounts
  name = each.value
}

resource "random_string" "random" {
  for_each = var.accounts
  length = 5
  special = false
}

resource "codefresh_api_key" "new" {
  for_each = var.accounts
  account_id = data.codefresh_account.acc[each.value].id
  user_id = data.codefresh_account.acc[each.value].admins[0]
  name = "tfkey_${random_string.random[each.value].result}"

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

output "tokens" {
  value = {
    for acc, token in codefresh_api_key.new:
      acc => token.token
  }  
}