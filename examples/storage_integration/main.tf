resource "codefresh_context" "gcs" {
  for_each = toset(["create"])
  name     = "gcs"
  spec {
    storagegc {
      data {
        auth {
          type = "basic"
          json_config = tomap({
            "config" : "cf"
          })
        }
      }
    }
  }
}

resource "codefresh_context" "s3" {
  for_each = toset(["create"])
  name     = "s3"
  spec {
    storages3 {
      data {
        auth {
          type = "basic"
          json_config = tomap({
            "config" : "cf"
          })
        }
      }
    }
  }
}

resource "codefresh_context" "azure" {
  for_each = toset(["create"])
  name     = "azure"
  spec {
    storageazuref {
      data {
        auth {
          type         = "basic"
          account_name = "accName"
          account_key  = "accKey"
        }
      }
    }
  }
}
