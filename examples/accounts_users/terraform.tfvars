default_idps = {
  local = {
    display_name = "local"
    sso          = false
  }
  azure_sso = {
    display_name = "azure-sso-1"
    sso          = true
  }
}

accounts = {
  acc1 = {}
  acc2 = {
    limits = {
      collaborators   = 50
      parallel_builds = 5
    }
  }
}

users = {
  user1 = {
    email = "user1@example.com"
    personal = {
      first_name = "John"
      last_name  = "Smith"
    }
    accounts          = ["acc1", "acc2"]
    admin_of_accounts = ["acc1"]
    global_admin      = true
  }
  user2 = {
    email = "live.com#user2@gmail.com"
    personal = {
      first_name = "Q"
      last_name  = "D"
    }
    accounts          = ["acc2"]
    admin_of_accounts = []
    global_admin      = false
  }
  user3 = {
    email = "user3@example.com"
    personal = {
      first_name = "Sam"
      last_name  = "Johnson"
    }
    accounts          = ["acc1", "acc2"]
    admin_of_accounts = ["acc1", "acc2"]
    global_admin      = true
  }
}