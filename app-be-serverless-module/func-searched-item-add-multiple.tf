#Put resources here.
#If this files become to long, you can move related resources to their own files.

#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
  route_of_searched_item_add_multiple = "POST /v1/searchedItem.addMultiple"
}

resource "null_resource" "go_build_functions_output_searched_item_add_multiple" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = "go run build_functions/main.go searched-item-add-multiple"
  }
}

data "archive_file" "lambda_searched_item_add_multiple" {
  type       = "zip"
  depends_on = [null_resource.go_build_functions_output_searched_item_add_multiple]

  source_file = "${path.module}/dist/functions/searched-item-add-multiple"
  output_path = "${path.module}/dist/functions/searched-item-add-multiple.zip"
}

resource "aws_lambda_function" "searched_item_add_multiple" {
  function_name = "SearchedItemAddMultiple"
  filename      = data.archive_file.lambda_searched_item_add_multiple.output_path
  runtime       = "go1.x"
  handler       = "searched-item-add-multiple"
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

  source_code_hash = data.archive_file.lambda_searched_item_add_multiple.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "searched_item_add_multiple" {
  name = "/aws/lambda/${aws_lambda_function.searched_item_add_multiple.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "searched_item_add_multiple" {
  api_id = aws_apigatewayv2_api.lambda.id

  integration_uri    = aws_lambda_function.searched_item_add_multiple.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "searched_item_add_multiple" {
  api_id = aws_apigatewayv2_api.lambda.id

  route_key = local.route_of_searched_item_add_multiple
  target    = "integrations/${aws_apigatewayv2_integration.searched_item_add_multiple.id}"
}

resource "aws_lambda_permission" "searched_item_add_multiple" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.searched_item_add_multiple.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}
