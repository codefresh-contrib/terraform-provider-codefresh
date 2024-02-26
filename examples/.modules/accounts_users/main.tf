data "codefresh_idps" "idps" {
  for_each = var.default_idps
  _id = lookup(each.value, "_id", "")
  display_name = lookup(each.value, "display_name", "")
  client_name = lookup(each.value, "client_name", "")
  client_type = lookup(each.value, "client_type", "")
}

resource "codefresh_account" "acc" {
  for_each = var.accounts
  name = each.key

  features = var.default_account_features

  limits {
    collaborators = lookup(var.default_acccount_limits, "collaborators", 10)
  }

  build {
    parallel = lookup(var.default_acccount_limits, "parallel_builds", 1)
  }

}

resource "codefresh_idp_accounts" "acc_idp" {
  for_each = var.default_idps
  idp_id = data.codefresh_idps.idps[each.key].id
  account_ids = values(codefresh_account.acc)[*].id 
}

resource "codefresh_user" "users" {
  for_each = var.users
  user_name = each.key
  email = each.value.email
  
  accounts = [
    for acc_name in each.value.accounts: codefresh_account.acc[acc_name].id
  ]

  activate = true

  roles = each.value.global_admin ? ["Admin","User"] : ["User"]

  dynamic "login" {
    for_each = var.default_idps
    content {
      idp_id = data.codefresh_idps.idps[login.key].id
      sso = login.value.sso
    }      
  }

  personal {
    first_name = each.value.personal.first_name
    last_name = each.value.personal.last_name
  }
}

resource "codefresh_account_admins" "acc_admins" {
  for_each = toset(flatten([
    for u in var.users:
      u.admin_of_accounts if length(u.admin_of_accounts) > 0
  ]))

  account_id = codefresh_account.acc[each.value].id
  users = [
      for k, u in var.users:
        codefresh_user.users[k].id if contains(u.admin_of_accounts, each.key) 
  ]
}