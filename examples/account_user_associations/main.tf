resource "codefresh_account_user_association" "user" {
  email = "terraform-test-user+user@codefresh.io"
}

resource "codefresh_account_user_association" "admin" {
  email = "terraform-test-user+admin@codefresh.io"
  admin = true
}