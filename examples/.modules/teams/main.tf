data "codefresh_current_account" "acc" {

}

locals {
  user_ids = tomap({
    for u in data.codefresh_current_account.acc.users:
      u.name => u.id
  })

}

resource "codefresh_team" "teams" {
  for_each = var.teams
  name = each.key

  users = [for u in each.value: lookup(local.user_ids, u)]
}
