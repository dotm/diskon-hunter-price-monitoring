locals {
    route_of_monitored_link_edit_multiple = "POST /v1/monitoredLink.editMultiple"
}

resource "null_resource" "go_build_functions_output_monitored_link_edit_multiple" {
    triggers = {
        always_run = "${timestamp()}"
    }
    provisioner "local-exec" {
        command = "go run build_functions/main.go monitored-link-edit-multiple"
    }
}

data "archive_file" "lambda_monitored_link_edit_multiple" {
    type = "zip"
    depends_on = [null_resource.go_build_functions_output_monitored_link_edit_multiple]
    source_file = "${path.module}/dist/functions/monitored-link-edit-multiple"
    output_path = "${path.module}/dist/functions/monitored-link-edit-multiple.zip"
}

resource "aws_lambda_function" "monitored_link_edit_multiple" {
  function_name = "MonitoredLinkEditMultiple"
  filename      = data.archive_file.lambda_monitored_link_edit_multiple.output_path
  runtime       = "go1.x"
  handler       = "monitored-link-edit-multiple"
  publish       = true

  environment {
    variables = {
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

  source_code_hash = data.archive_file.lambda_monitored_link_edit_multiple.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "monitored_link_edit_multiple" {
  name = "/aws/lambda/${aws_lambda_function.monitored_link_edit_multiple.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "monitored_link_edit_multiple" {
  api_id = aws_apigatewayv2_api.lambda.id

  integration_uri    = aws_lambda_function.monitored_link_edit_multiple.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "monitored_link_edit_multiple" {
  api_id = aws_apigatewayv2_api.lambda.id

  route_key = local.route_of_monitored_link_edit_multiple
  target    = "integrations/${aws_apigatewayv2_integration.monitored_link_edit_multiple.id}"
}

resource "aws_lambda_permission" "monitored_link_edit_multiple" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.monitored_link_edit_multiple.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}



