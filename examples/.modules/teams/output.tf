output "users" {
  value = local.user_ids 
}
output "teams" {
  value = codefresh_team.teams
}
