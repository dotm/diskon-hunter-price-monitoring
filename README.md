# README

## Onboarding for new developers

All developers:

- Read the README in the subdirectories to understand how to code in each subdirectory.
- Read the other sections of this README below.
- Recommended code editor: the one you're comfortable with, or Visual Studio Code
- Recommended development workflow: one code editor window per module
  - This is because Go lang has problem when one VSCode editor window contains multiple Go project.
  - It'll also avoid the case of you editing one file in module A when you intend to actually edit a file of the same name in module B.

Backend developer:

- Learn: Terraform, DynamoDB, Go Lang.
    - Important DynamoDB concepts: partition key (hash_key), sort key (range_key), put vs update operation.
- Ask for an AWS account (email and password).
  - This will be your personal account and will be used to deploy to `deployment_environment_name = "local"`.
- Best practice on newly created AWS account:
  - Create an admin IAM user (separate from root account) with password and access key.
    - In IAM create a new user. Assign the policy of AdministratorAccess.
    - Optionally, activate IAM access to billing information so that admin IAM user can access billing pages.
  - Create an AWS bugdet
    - Monthly cost budget (recommended amount is $1) with notification to your own email.
- Install AWS CLI (used in dynamodbhelper.CreateClientFromSession)
  - Run `aws configure`
  - Input AWS Access Key ID and AWS Secret Access Key based on the previously created access key
  - Default region is ap-southeast-1 (Singapore). Might be changed to Jakarta later.
  - Default output format is json
- Install terraform with the appropriate version.
  - Search `required_version` in this repo to find the appropriate terraform version.
- Deploy db and app-be module to your local environment in your personal AWS account.
- Start exploring the codebase and happy coding :)

Frontend developers:

- Learn: Flutter (+ Dart Lang), Android.
- Do the get started part from Flutter Docs (install flutter, setup code editor, test running the app on your phone)
- Start exploring the codebase and happy coding :)
- Deploy to your PC (Windows or Mac): `flutter run -d windows`
  - When developing, it's easier to deploy to PC first because:
    - no mobile phone is needed
    - you can easily resize the application window to check for all possible form factors (small phone, large phone, tablets, desktop, etc.)

## Offboarding developers

- Delete all Terraform workspace you've created:
  - Check README of this directory's subdirectories for more details.
  - For all the subdirectories
    that you've runned this command
      `terraform workspace new local-[insert-your-name-here]`
      and this command
      `terraform apply -var-file="local.tfvars"`,
    please run this command
      `terraform destroy -var-file="local.tfvars"`
      and this command
      `terraform workspace delete local-[insert-your-name-here]`.
  - Create a pull request to dev (or master if dev hasn't been created).
- If you have asked for an AWS account, please ask for that account to be closed.
  - Or you can close it yourself after destroying all the Terraform resources in your local workspace and deleting that local workspace.
- Ask for your read-write access to Git repository to be revoked.

## Long running branches and deployment environments

- The `master` branch is for `deployment_environment_name = "prod"`
  - If you need `prod.tfvars.sample`, you can fill:
    - `deployment_environment_name      = "prod"`
    - `deployment_environment_purpose   = "Enviroment accessed and used by end users"`
- The `staging` branch is for `deployment_environment_name = "staging"`
  - If you need `staging.tfvars.sample`, you can fill:
    - `deployment_environment_name      = "staging"`
    - `deployment_environment_purpose   = "Used by our partners to develop their products"`
- The `dev` branch is for `deployment_environment_name = "dev"`
  - If you need `dev.tfvars.sample`, you can fill:
    - `deployment_environment_name      = "dev"`
    - `deployment_environment_purpose   = "Latest features and fixes for all developers"`
- The `local/your-name-here` branch is for your personal `deployment_environment_name = "local"`
  - If you need `local.tfvars.sample`, you can fill:
    - `deployment_environment_name      = "local"`
    - `deployment_environment_purpose   = "Rapid local experimentation for an individual developer"`

## PR Ethics

- Create your own branch (example: `local/your-name-here`) for your personal `deployment_environment_name = "local"`
  - Base it from the correct branch (for example dev for development of new features, or staging/master for hotfix to the respective deployment environment).
  - Do your work in your own local branch, test it, and then create a Pull Request to the base branch once it's finished and ready to be merged.
    - Do backend and frontend changes for your features in one branch. Don't separate backend and frontend to their own branches.
  - Ask for reviewers to check your work and then merge the branch once approved.
- You can create multiple local branches. For example:
  - `local/your-name-here/feature/edit-company`
  - `local/your-name-here/hotfix/edit-company`
  - `local/your-name-here/whatever-name-you-want-here`

## Installing Terraform

We recommend using a version manager like `tfenv` to install Terraform
so that you can change between multiple Terraform versions easily.
But if it's too much work (especially in Windows), just install the appropriate version of Terraform.

## Sharing plugins across modules

This reduces the size of .terraform directory by providing a single source of aws plugin cache.

- `mkdir $HOME/.terraform.d/plugin-cache`
- `echo 'plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"' >> $HOME/.terraformrc`
  - You can also use: `export TF_PLUGIN_CACHE_DIR=$HOME/.terraform.d/plugin-cache`
- Restart shell

## Module Introduction

Each module should act as if they don't know other module.
Do NOT introduce any dependency that relies on the fact that they are in the same directory as this README.
Other dependencies through Terraform data sources, external key value storage, etc. is fine.

Notes on `0-` modules:

- `0-console-experiment`: use this for experimenting with `terraform console`
- `0-empty-module`: copy this to create a new module (don't forget to remove every `"0-empty-module"` especially in `variables.tf` file)

## Recommended Module Structure

https://learn.hashicorp.com/tutorials/terraform/pattern-module-creation

- Organization Module
- Network (network) Module
- DNS Zone???
- Database (db) Module: RDS, DynamoDB, ElasticSearch (ES), etc.
- App (app) Module:
  - Backend (BE)
  - Frontend (FE)
- Security Module
  - Authentication and/or Authorization (Auth)
    - This can be put inside App BE for custom auth solution.
  - Anti Distributed Denial of Service (DDoS)
