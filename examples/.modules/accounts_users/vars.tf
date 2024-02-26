variable api_url {
  type = string
}

# 
variable token {
  type = string
  default = ""
}

variable default_account_features {
  type = map(any)
  default = {
    OfflineLogging = true,
    ssoManagement = true,
    teamsManagement = true,
    abac = true,
    customKubernetesCluster = true,
    launchDarklyManagement = false,
  }
}

variable default_acccount_limits {
    type = map(any)
    default = {
        collaborators = 100
        parallel_builds = 10
    }
}

variable default_idps {
    type = map(any)   
    default = {
        local  = { 
            display_name = "local"
            sso = false
        }
    } 
}

variable accounts {
  type = map(any)
}

variable users {
  //type = map(any)
}

