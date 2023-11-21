# README

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
