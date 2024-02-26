variable "test_password" {
  type    = string
  default = ""
}

resource "codefresh_registry" "acr" {
  name    = "acr"
  default = true
  spec {
    acr {
      domain            = "acr.io"
      client_id         = "test"
      client_secret     = "test"
      repository_prefix = "test"
    }
  }
}

resource "codefresh_registry" "gcr" {
  name = "gcr"
  #  all registries SHOULD be dependent on each other to be created/updated sequentially
  depends_on = [codefresh_registry.acr]
  spec {
    gcr {
      domain            = "gcr.io"
      keyfile           = "test"
      repository_prefix = "codefresh-inc"
    }
  }
}

resource "codefresh_registry" "gar" {
  name       = "gar"
  depends_on = [codefresh_registry.gcr]
  spec {
    gar {
      location          = "asia"
      keyfile           = "test"
      repository_prefix = "codefresh-inc"
    }
  }
}

data "codefresh_registry" "dockerhub" {
  name = "dockerhub"
}

# example with using data reference to existing registry, not managed by terraform
resource "codefresh_registry" "dockerhub1" {
  name       = "dockerhub1"
  primary    = !data.codefresh_registry.dockerhub.primary
  depends_on = [codefresh_registry.gar]
  spec {
    dockerhub {
      username = "test"
      password = "test"
    }
  }
  fallback_registry = data.codefresh_registry.dockerhub.id
}

resource "codefresh_registry" "bintray" {
  name       = "bintray"
  depends_on = [codefresh_registry.dockerhub1]
  spec {
    bintray {
      domain   = "bintray.io"
      username = "test"
      token    = "test"
    }
  }
}

resource "codefresh_registry" "other" {
  name       = "other"
  depends_on = [codefresh_registry.bintray]
  primary    = false
  spec {
    other {
      domain   = "other.io"
      username = "test"
      password = "test"
    }
  }
}

# when you have multiple registries under the same domain
# they MUST be dependant on each other and `primary`
# MUST be specified at least and only for one registry
# as `true`
resource "codefresh_registry" "other1" {
  name       = "other1"
  primary    = true
  depends_on = [codefresh_registry.other]
  spec {
    other {
      domain   = "other.io"
      username = "test"
      password = "test"
    }
  }
}

resource "codefresh_registry" "other2" {
  name       = "other2"
  primary    = false
  depends_on = [codefresh_registry.other1, codefresh_registry.bintray]
  spec {
    other {
      domain   = "other.io"
      username = "test"
      password = "test"
    }
  }
  fallback_registry = codefresh_registry.bintray.id
}
