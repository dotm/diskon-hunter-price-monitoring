#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
}

resource "aws_dynamodb_table" "user" {
  lifecycle {
    prevent_destroy = true
  }

  #user means anyone or anything that can sign in to access application features.

  name         = "${var.deployment_environment_name}-${var.project_name_short}-StlUserDetail"
  billing_mode = "PAY_PER_REQUEST" #Default value is PROVISIONED. On demand is PAY_PER_REQUEST.

  #if no range_key, pk value must be unique, else pk-sk combination value must be unique
  hash_key = "HubUserId" #attribute used as partition key (beware of hot partition problem)
  #other fields: check DAO on app-be code base

  #you can only specify indexed attributes (pk, sk, lsi, gsi)
  attribute {
    name = "HubUserId"
    type = "S" #(S)tring, (N)umber or (B)inary
  }
}

resource "aws_dynamodb_table" "user_email_authentication" {
  lifecycle {
    prevent_destroy = true
  }

  #user means anyone or anything that can sign in to access application features.

  name         = "${var.deployment_environment_name}-${var.project_name_short}-StlUserEmailAuthentication"
  billing_mode = "PAY_PER_REQUEST" #Default value is PROVISIONED. On demand is PAY_PER_REQUEST.

  #if no range_key, pk value must be unique, else pk-sk combination value must be unique
  hash_key = "Email" #attribute used as partition key (beware of hot partition problem)
  #other fields: check DAO on app-be code base

  #you can only specify indexed attributes (pk, sk, lsi, gsi)
  attribute {
    name = "Email"
    type = "S" #(S)tring, (N)umber or (B)inary
  }
}

resource "aws_dynamodb_table" "user_email_has_otp" {
  lifecycle {
    prevent_destroy = true
  }

  #user means anyone or anything that can sign in to access application features.

  name         = "${var.deployment_environment_name}-${var.project_name_short}-StlUserEmailHasOtpDetail"
  billing_mode = "PAY_PER_REQUEST" #Default value is PROVISIONED. On demand is PAY_PER_REQUEST.

  #if no range_key, pk value must be unique, else pk-sk combination value must be unique
  hash_key = "Email" #attribute used as partition key (beware of hot partition problem)
  ttl {
    attribute_name = "TimeExpired"
    enabled        = true
  }
  #other fields: check DAO on app-be code base

  #you can only specify indexed attributes (pk, sk, lsi, gsi)
  attribute {
    name = "Email"
    type = "S" #(S)tring, (N)umber or (B)inary
  }
}
