variable "aws_deployment_access_key_id" {
  type      = string
  nullable  = false
  sensitive = true
}

variable "aws_deployment_secret_access_key" {
  type      = string
  nullable  = false
  sensitive = true
}

variable "aws_deployment_account_id" {
  type     = string
  nullable = false
}

variable "aws_deployment_region" {
  type     = string
  nullable = false
  default  = "ap-southeast-1"
}

variable "aws_deployment_region_short" {
  description = "Short/abbreviated version of aws_deployment_region. Useful for naming resources."

  type     = string
  nullable = false
  default  = "apse1"

  validation {
    //must refer and can only refer to self var 
    condition     = length(var.aws_deployment_region_short) < 6
    error_message = "aws_deployment_region_short must have length of less than 6"
  }
}

variable "deployment_environment_name" {
  type     = string
  nullable = false
}

variable "deployment_environment_purpose" {
  type     = string
  nullable = false
}

variable "project_name" {
  type     = string
  nullable = false
  default  = "diskon-hunter-price-monitoring"
}

variable "project_name_short" {
  description = "Short/abbreviated version of project_name. Useful for naming resources."

  type     = string
  nullable = false
  default  = "dhpm"
}

variable "module_name" {
  type     = string
  nullable = false
  default  = "db-dynamo"
}
