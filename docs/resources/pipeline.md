# Pipeline Resource

The central component of the Codefresh Platform. Pipelines are workflows that contain individual steps. Each step is responsible for a specific action in the process.
See the [documentation](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/introduction-to-codefresh-pipelines/) for the details.

## Example Usage

```hcl
resource "codefresh_project" "test" {
  name = "myproject"
}

resource "codefresh_pipeline" "test" {

  lifecycle {
    ignore_changes = [
      revision
    ]
  }

  name    = "${codefresh_project.test.name}/react-sample-app"

  tags = [
    "production",
    "docker",
  ]

  spec {
    concurrency         = 1
    branch_concurrency  = 1
    trigger_concurrency = 1

    priority    = 5

    spec_template {
      repo        = "codefresh-contrib/react-sample-app"
      path        = "./codefresh.yml"
      revision    = "master"
      context     = "git"
    }

    contexts = [
      "context1-name",
      "context2-name",
    ]

    trigger {
      branch_regex  = "/.*/gi"
      context       = "git"
      description   = "Trigger for commits"
      disabled      = false
      events        = [
        "push.heads"
      ]
      modified_files_glob = ""
      name                = "commits"
      provider            = "github"
      repo                = "codefresh-contrib/react-sample-app"
      type                = "git"
    }

    trigger {
      branch_regex  = "/.*/gi"
      context       = "git"
      description   = "Trigger for tags"
      disabled      = false
      events        = [
        "push.tags"
      ]
      modified_files_glob = ""
      commit_status_title = "tags-trigger"
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
```

## Argument Reference

- `name` - (Required) The display name for the pipeline.
- `revision` - (Optional) The pipeline's revision. Should be added to the **lifecycle/ignore_changes** or incremented mannually each update.
- `is_public` - (Optional) Boolean that specifies if the build logs are publicly accessible. Default: false
- `tags` - (Optional) A list of tags to mark a project for easy management and access control.
- `spec` - (Required) A collection of `spec` blocks as documented below.
- `original_yaml_string` - (Optional) A string with original yaml pipeline.
  - `original_yaml_string = "version: \"1.0\"\nsteps:\n  test:\n    image: alpine:latest\n    commands:\n      - echo \"ACC tests\""`
  - or `original_yaml_string = file("/path/to/my/codefresh.yml")`

---

`spec` supports the following:

- `concurrency` - (Optional) The maximum amount of concurrent builds.
- `branch_concurrency` - (Optional) The maximum amount of concurrent builds that may run for each branch
- `trigger_concurrency` - (Optional) The maximum amount of concurrent builds that may run for each trigger.
- `priority` - (optional) Helps to organize the order of builds execution in case of reaching the concurrency limit.
- `variables` - (Optional) Pipeline variables.
- `trigger` - (Optional) A collection of `trigger` blocks as documented below. Triggers [documentation](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/triggers/git-triggers/).
- `spec_template` - (Optional) A collection of `spec_template` blocks as documented below.
- `runtime_environment` - (Optional) A collection of `runtime_environment` blocks as documented below.
- `contexts` - (Optional) A list of strings representing the contexts ([shared_configuration](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/shared-configuration/)) to be configured for the pipeline
- `termination_policy` - (Optional) A `termination_policy` block as documented below.
- `options` - (Optional) A `options` block as documented below.

---

`spec_template` supports the following:

- `location` - (Optional) Default value - **git**.
- `repo` - (Required) The GitHub `account/repo_name`.
- `path` - (Required) The relative path to the Codefresh pipeline file.
- `revison` - (Required) The git revision.
- `context` - (Optional) The Codefresh Git [context](https://codefresh.io/docs/docs/integrations/git-providers/).

---

`trigger` supports the following:

- `name` - (Optional) The display name for the pipeline.
- `description` - (Optional) The trigger description.
- `type` - (Optional) The trigger type. Default value - **git**.
- `repo` - (Optional) The GitHub `account/repo_name`.
- `branch_regex` - (Optional) A regular expression and will only trigger for branches that match this naming pattern.
- `branch_regex_input` - (Optional) Flag to manage how the `branch_regex` field is interpreted. Possible values: "multiselect-exclude", "multiselect", "regex". Default: "regex"
- `pull_request_target_branch_regex` - (Optional) A regular expression and will only trigger for pull requests to branches that match this naming pattern.
- `comment_regex` - (Optional) A regular expression and will only trigger for pull requests where a comment matches this naming pattern.
- `modified_files_glob` - (Optional) Allows to constrain the build and trigger it only if the modified files from the commit match this glob expression.
- `events` - (Optional) A list of GitHub events for which a Pipeline is triggered. Default value - **push.heads**.
- `provider` - (Optional) Default value - **github**.
- `context` - (Optional) Codefresh Git context.
- `commit_status_title` - (Optional) The commit status title pushed to the GIT version control system.
- `variables` - (Optional) Trigger variables.
- `disabled` - (Optional) Boolean. If false, trigger will never be activated.
- `pull_request_allow_fork_events` - (Optional) Boolean. If this trigger is also applicable to Git forks.
- `contexts` - (Optional) A list of strings representing the contexts ([shared_configuration](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/shared-configuration/)) to be loaded when the trigger is executed
- `runtime_environment` - (Optional) A collection of `runtime_environment` blocks as documented below.
- `options`: (Optional) A collection `option` blocks as documented below.
---

`runtime_environment` supports the following:

- `name` - (Required) A name of runtime.
- `cpu` - (Optional) A required amount of CPU.
- `memory` - (Optional) A required amount of memory.
- `dind_storage` - (Optional) A pipeline shared storage.

---

`options` supports the following:

- `no_cache` - (Required) Boolean. If true, docker layer cache is disabled. Default false
- `no_cf_cache` - (Optional) Boolean. If true, extra Codefresh caching is disabled. Default false
- `reset_volume` - (Optional) Boolean. If true, all files on volume will be deleted before each execution. Default false
- `enable_notifications` - (Optional) Boolean. If false the pipeline will not send notifications to Slack and status updates back to the Git provider. Default false

---

`termination_policy` supports the following:

- `on_create_branch` - (Optional) A `on_create_branch` block as documented below.
- `on_terminate_annotation` - (Optional) Boolean. Enables the policy `Once a build is terminated, terminate all child builds initiated from it`. Default false.

---

`on_create_branch` supports the following:

- `branch_name` - (Optional) A regular expression to filter the branches on with the termination policy applies.
- `ignore_trigger` - (Optional) Boolean. See table below for usage.
- `ignore_branch` - (Optional) Boolean. See table below for usage.

The following table presents how to configure this block based on the options available in the UI:
| Option Description                                                            | Value Selected           | on_create_branch | branch_name | ignore_trigger | ignore_branch |
| ----------------------------------------------------------------------------- |:------------------------:|:----------------:|:-----------:|---------------:| -------------:|
| Once a build is created terminate previous builds from the same branch        | Disabled                 |        Omit      |     N/A     |       N/A      |      N/A      |
| Once a build is created terminate previous builds from the same branch        | From the SAME trigger    |       Defined    |     N/A     |      false     |      N/A      |
| Once a build is created terminate previous builds from the same branch        | From ANY trigger         |       Defined    |     N/A     |      true      |      N/A      |
| Once a build is created terminate previous builds only from a specific branch | Disabled                 |        Omit      |     N/A     |       N/A      |      N/A      |
| Once a build is created terminate previous builds only from a specific branch | From the SAME trigger    |       Defined    |    Regex    |      false     |      N/A      |
| Once a build is created terminate previous builds only from a specific branch | From ANY trigger         |       Defined    |    Regex    |      true      |      N/A      |
| Once a build is created, terminate all other running builds                   | Disabled                 |        Omit      |     N/A     |       N/A      |      N/A      |
| Once a build is created, terminate all other running builds                   | From the SAME trigger    |       Defined    |     N/A     |      false     |      true     |
| Once a build is created, terminate all other running builds                   | From ANY trigger         |       Defined    |     N/A     |      true      |      true     |

---

`options` supports the following:

- `keep_pvcs_for_pending_approval` - (Optional) Boolean for the Settings under pending approval: `When build enters "Pending Approval" state, volume should`:
    * Default (attribute not specified): "Use Setting accounts"
    * true: "Remain (build remains active)"
    * false: "Be removed"
- `pending_approval_concurrency_applied` - (Optional) Boolean for the Settings under pending approval: `Pipeline concurrency policy: Builds on "Pending Approval" state should be`:
    * Default (attribute not specified): "Use Setting accounts"
    * true: "Included in concurrency"
    * false: "Not included in concurrency"

## Attributes Reference

- `id` - The Pipeline ID.

## Import

```sh
terraform import codefresh_pipeline.test xxxxxxxxxxxxxxxxxxx
```
