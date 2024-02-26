## Set of account names
variable "accounts" {
  type = set(string)
}

module "account_tokens" {
  source   = "../.modules/account_tokens"
  accounts = var.accounts
}

output "account_tokens" {
  value = module.account_tokens.tokens
}