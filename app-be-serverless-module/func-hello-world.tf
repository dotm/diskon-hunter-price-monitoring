#Put resources here.
#If this files become to long, you can move related resources to their own files.

#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
}

resource "null_resource" "go_build_functions_output_hello_world" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = "go run build_functions/main.go hello-world"
  }
}

data "archive_file" "lambda_hello_world" {
  type       = "zip"
  depends_on = [null_resource.go_build_functions_output_hello_world]

  source_file = "${path.module}/dist/functions/hello-world"
  output_path = "${path.module}/dist/functions/hello-world.zip"
}

resource "aws_lambda_function" "hello_world" {
  function_name = "HelloWorld"
  filename      = data.archive_file.lambda_hello_world.output_path
  runtime       = "go1.x"
  handler       = "hello-world"
  publish       = true

  #increasing memory size of a lambda will also increase it's cpu
  memory_size = 128 #in MB. from 128 (default) up to 10240.
  timeout     = 15  #in seconds. Amount of time your Lambda Function has to run. Defaults to 3

  source_code_hash = data.archive_file.lambda_hello_world.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "hello_world" {
  name = "/aws/lambda/${aws_lambda_function.hello_world.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "hello_world" {
  api_id = aws_apigatewayv2_api.lambda.id

  integration_uri    = aws_lambda_function.hello_world.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "hello_world" {
  api_id = aws_apigatewayv2_api.lambda.id

  route_key = "GET /hello"
  target    = "integrations/${aws_apigatewayv2_integration.hello_world.id}"
}

resource "aws_lambda_permission" "api_gw_hello_world" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.hello_world.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}
