module "pipeline" {
  source  = "../pipelines"
}

resource "codefresh_pipeline_cron_trigger" "default" {
  pipeline_id = module.pipeline.id
  expression  = "0 0 2 ? * MON-FRI"
  message     = "test"
}
