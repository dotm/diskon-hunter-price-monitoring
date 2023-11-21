This is an example module that could be copied.

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

## Deploy to prod

First deployment:

- `terraform init`
- `terraform workspace new prod`
- `terraform workspace select prod; terraform apply -var-file="prod.tfvars"`
  - We use multiple command in one line to make sure we don't forget to select the correct workspace before applying the resources.

Not the first deployment:

- `terraform workspace select prod; terraform apply -var-file="prod.tfvars"`
  - We use multiple command in one line to make sure we don't forget to select the correct workspace before applying the resources.