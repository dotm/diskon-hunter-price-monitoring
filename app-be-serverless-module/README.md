This is a serverless module for the backend layer.

AWS Lambda functions and API gateway are often used to create serverless applications.

Follow along with this [tutorial on HashiCorp Learn](https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway?in=terraform/aws).

## Directory structure explanation

- `./build_functions`
  - used to build the Go program inside `./dist/functions` that will be uploaded to AWS Lambda
- `./dist/functions` (git ignored)
  - used by build_functions to put latest compiled Go program.
- `./functions/*`
  - The entry point for the Lambda handlers inside `./src`
- `./playground`
  - used to test assumptions quickly. You can import from any file inside this module. And you can just fmt.Print to check the result of your tests.
- `./shared`
  - for any utilities or struct that is used by multiple Lambda functions.
  - TODO: explain all directory inside shared
- `./src/*`
  - the handlers of the Lambda functions.
  - TODO: explain DTO, DAO, command/query, etc.

## HOW TO create a new API endpoint

- Add new terraform file. For example: `./func-company-edit-multiple.tf`
  - You can copy, for example from `./func-company-add-multiple.tf` and then replace all `add` with `edit` and all `Add` with `Edit`.
- Add new directory in `./functions`. For example: `./functions/company-edit-multiple`
  - You can copy, for example from `./functions/company-add-multiple` and then replace all `add` with `edit` and all `Add` with `Edit`.
- Add the handler in the `./src` subdirectory. For example: `./src/company/editMultiple`

## Deploy to your local cloud

First deployment:

- Setup an AWS account
- Setup AWS CLI and configure it
- Create a local.tfvars based on local.tfvars.sample
- Make sure you use the directory this README is in as root directory
- `terraform init`
- `terraform workspace new local-[insert-your-name-here]`
- `terraform workspace select local-[insert-your-name-here]; terraform apply -var-file="local.tfvars"`
  - We use multiple command in one line to make sure we don't forget to select the correct workspace before applying the resources.

Not the first deployment:

- `terraform workspace select local-[insert-your-name-here]; terraform apply -var-file="local.tfvars"`
  - We use multiple command in one line to make sure we don't forget to select the correct workspace before applying the resources.
- To deploy specific lambda function:
    - `terraform workspace select local-[insert-your-name-here]; terraform apply -var-file="local.tfvars" -target="aws_lambda_function.hello_world"`
    - This will not deploy the respective API Gateway terraform resources. If you want to deploy the lambda function for the first time, don't use the target flag.

## Deploy to prod

First deployment:

- `terraform init`
- `terraform workspace new prod`
- `terraform workspace select prod; terraform apply -var-file="prod.tfvars"`
  - We use multiple command in one line to make sure we don't forget to select the correct workspace before applying the resources.
- Move AWS SES out of sandbox:
  - https://docs.aws.amazon.com/ses/latest/dg/request-production-access.html
- Configuration > Verified identities > Create identity:
  - Identity type: email address
  - Tip: use @yopmail.com for disposable email address

Not the first deployment:

- `terraform workspace select prod; terraform apply -var-file="prod.tfvars"`
  - We use multiple command in one line to make sure we don't forget to select the correct workspace before applying the resources.

## TODO

- add CRUD to dynamodb
- refactor main.tf
- add stages (prod, staging, qa) (don't forget to change api gateway tf config)
- add tests

## Initial Project Setup

- go mod init diskon-hunter/price-monitoring
- go get github.com/aws/aws-lambda-go

## Testing

- curl "$(terraform output -raw base_url)/hello?Name=Terraform"

## Run Backend in Local

- go run *.go
