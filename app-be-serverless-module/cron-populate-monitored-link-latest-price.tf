#Put resources here.
#If this files become to long, you can move related resources to their own files.

#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
  cron_populate_monitored_link_latest_price_schedule = "cron(0 8 * * ? *)" #15:00 UTC+7 (8:00 UTC)
  cron_populate_monitored_link_latest_price_is_enabled = false
}

resource "null_resource" "go_build_functions_output_cron_populate_monitored_link_latest_price" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = "go run build_functions/main.go cron-populate-monitored-link-latest-price"
  }
}

data "archive_file" "lambda_cron_populate_monitored_link_latest_price" {
  type       = "zip"
  depends_on = [null_resource.go_build_functions_output_cron_populate_monitored_link_latest_price]

  source_file = "${path.module}/dist/functions/cron-populate-monitored-link-latest-price"
  output_path = "${path.module}/dist/functions/cron-populate-monitored-link-latest-price.zip"
}

resource "aws_lambda_function" "cron_populate_monitored_link_latest_price" {
  function_name = "CronPopulateMonitoredLinkLatestPrice"
  filename      = data.archive_file.lambda_cron_populate_monitored_link_latest_price.output_path
  runtime       = "go1.x"
  handler       = "cron-populate-monitored-link-latest-price"
  publish       = true

  environment {
    variables = {
      aws_deployment_account_id      = var.aws_deployment_account_id
      aws_deployment_region_short    = var.aws_deployment_region_short
      deployment_environment_name    = var.deployment_environment_name
      deployment_environment_purpose = var.deployment_environment_purpose
      project_name                   = var.project_name
      project_name_short             = var.project_name_short
      module_name                    = var.module_name
    }
  }

  #increasing memory size of a lambda will also increase it's cpu
  memory_size = 2048 #in MB. from 128 (default) up to 10240.
  timeout     = 15   #in seconds. Amount of time your Lambda Function has to run. Defaults to 3

  source_code_hash = data.archive_file.lambda_cron_populate_monitored_link_latest_price.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "cron_populate_monitored_link_latest_price" {
  name = "/aws/lambda/${aws_lambda_function.cron_populate_monitored_link_latest_price.function_name}"

  retention_in_days = 30
}

resource "aws_lambda_permission" "cron_populate_monitored_link_latest_price" {
  statement_id = "AllowExecutionFromCloudWatch"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cron_populate_monitored_link_latest_price.function_name
  principal = "events.amazonaws.com"
  
  source_arn = aws_cloudwatch_event_rule.cron_populate_monitored_link_latest_price.arn
}

resource "aws_cloudwatch_event_rule" "cron_populate_monitored_link_latest_price" {
  name = "schedule_for_cron_populate_monitored_link_latest_price"
  description = "Schedule for Lambda Function"
  schedule_expression = local.cron_populate_monitored_link_latest_price_schedule
  is_enabled = local.cron_populate_monitored_link_latest_price_is_enabled
}

resource "aws_cloudwatch_event_target" "cron_populate_monitored_link_latest_price" {
  rule = aws_cloudwatch_event_rule.cron_populate_monitored_link_latest_price.name
  target_id = "cron_populate_monitored_link_latest_price"
  arn = aws_lambda_function.cron_populate_monitored_link_latest_price.arn
}