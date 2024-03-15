#Put resources here.
#If this files become to long, you can move related resources to their own files.

#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
  cron_send_monitored_link_alert_schedule = "cron(0 9 * * ? *)" #16:00 UTC+7 (9:00 UTC)
  cron_send_monitored_link_alert_is_enabled = false
}

resource "null_resource" "go_build_functions_output_cron_send_monitored_link_alert" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = "go run build_functions/main.go cron-send-monitored-link-alert"
  }
}

data "archive_file" "lambda_cron_send_monitored_link_alert" {
  type       = "zip"
  depends_on = [null_resource.go_build_functions_output_cron_send_monitored_link_alert]

  source_file = "${path.module}/dist/functions/cron-send-monitored-link-alert"
  output_path = "${path.module}/dist/functions/cron-send-monitored-link-alert.zip"
}

resource "aws_lambda_function" "cron_send_monitored_link_alert" {
  function_name = "CronSendMonitoredLinkAlert"
  filename      = data.archive_file.lambda_cron_send_monitored_link_alert.output_path
  runtime       = "go1.x"
  handler       = "cron-send-monitored-link-alert"
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

  source_code_hash = data.archive_file.lambda_cron_send_monitored_link_alert.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "cron_send_monitored_link_alert" {
  name = "/aws/lambda/${aws_lambda_function.cron_send_monitored_link_alert.function_name}"

  retention_in_days = 30
}

resource "aws_lambda_permission" "cron_send_monitored_link_alert" {
  statement_id = "AllowExecutionFromCloudWatch"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cron_send_monitored_link_alert.function_name
  principal = "events.amazonaws.com"
  
  source_arn = aws_cloudwatch_event_rule.cron_send_monitored_link_alert.arn
}

resource "aws_cloudwatch_event_rule" "cron_send_monitored_link_alert" {
  name = "schedule_for_cron_send_monitored_link_alert"
  description = "Schedule for Lambda Function"
  schedule_expression = local.cron_send_monitored_link_alert_schedule
  is_enabled = local.cron_send_monitored_link_alert_is_enabled
}

resource "aws_cloudwatch_event_target" "cron_send_monitored_link_alert" {
  rule = aws_cloudwatch_event_rule.cron_send_monitored_link_alert.name
  target_id = "cron_send_monitored_link_alert"
  arn = aws_lambda_function.cron_send_monitored_link_alert.arn
}