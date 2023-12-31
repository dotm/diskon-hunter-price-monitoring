#Put resources here.
#If this files become to long, you can move related resources to their own files.

#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
}

resource "null_resource" "go_build_functions_output_app_version" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = "go run build_functions/main.go app-version"
  }
}

data "archive_file" "lambda_app_version" {
  type       = "zip"
  depends_on = [null_resource.go_build_functions_output_app_version]

  source_file = "${path.module}/dist/functions/app-version"
  output_path = "${path.module}/dist/functions/app-version.zip"
}

resource "aws_lambda_function" "app_version" {
  function_name = "AppVersion"
  filename      = data.archive_file.lambda_app_version.output_path
  runtime       = "go1.x"
  handler       = "app-version"
  publish       = true

  #increasing memory size of a lambda will also increase it's cpu
  memory_size = 128 #in MB. from 128 (default) up to 10240.
  timeout     = 15  #in seconds. Amount of time your Lambda Function has to run. Defaults to 3

  source_code_hash = data.archive_file.lambda_app_version.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "app_version" {
  name = "/aws/lambda/${aws_lambda_function.app_version.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "app_version" {
  api_id = aws_apigatewayv2_api.lambda.id

  integration_uri    = aws_lambda_function.app_version.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "app_version" {
  api_id = aws_apigatewayv2_api.lambda.id

  route_key = "GET /appVersion"
  target    = "integrations/${aws_apigatewayv2_integration.app_version.id}"
}

resource "aws_lambda_permission" "api_gw_app_version" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.app_version.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}
