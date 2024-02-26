resource "codefresh_project" "test" {
  name = "myproject"

  tags = [
    "docker",
  ]
}

resource "codefresh_pipeline" "test" {
  name = "${codefresh_project.test.name}/react-sample-app"

  tags = [
    "production",
    "docker",
  ]

  original_yaml_string = <<EOT
version: "1.0"
hooks: 
 on_finish:
   steps:
     b:
       image: alpine:3.9
       commands:
         - echo "echo cleanup step"
     a:
       image: cloudposse/slack-notifier
       commands:
         - echo "Notify slack"
steps:
    freestyle:
        image: alpine
        commands:
        - sleep 10
    a_freestyle:
        image: alpine
        commands:
        - sleep 10
        - echo Hey!
        arguments:
            c: 3
            a: 1
            b: 2
  EOT

  spec {
    concurrency = 1
    priority    = 5

    # spec_template {
    #   repo        = "codefresh-contrib/react-sample-app"
    #   path        = "./codefresh.yml"
    #   revision    = "master"
    #   context     = "git"
    # }

    contexts = [
      "context1-name",
      "context2-name",
    ]

    trigger {
      branch_regex = "/.*/gi"
      context      = "git"
      description  = "Trigger for commits"
      disabled     = false
      events = [
        "push.heads"
      ]
      modified_files_glob = ""
      name                = "commits"
      provider            = "github"
      repo                = "codefresh-contrib/react-sample-app"
      type                = "git"
    }

    trigger {
      branch_regex = "/.*/gi"
      context      = "git"
      description  = "Trigger for tags"
      disabled     = false
      events = [
        "push.tags"
      ]
      modified_files_glob = ""
      name                = "tags"
      provider            = "github"
      repo                = "codefresh-contrib/react-sample-app"
      type                = "git"
    }

    variables = {
      MY_PIP_VAR      = "value"
      ANOTHER_PIP_VAR = "another_value"
    }
  }
}