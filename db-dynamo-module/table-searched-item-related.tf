#Local variables that are used in multiple files should be placed in ./locals.tf
#Put local variables that are only used in this file below
locals {
  
}

resource "aws_dynamodb_table" "user_searches_item" {
  lifecycle {
    prevent_destroy = true
  }

  #this table represents many-to-many relationship

  name         = "${var.deployment_environment_name}-${var.project_name_short}-StlUserSearchesItemDetail"
  billing_mode = "PAY_PER_REQUEST" #Default value is PROVISIONED. On demand is PAY_PER_REQUEST.

  #if no range_key, pk value must be unique, else pk-sk combination value must be unique
  hash_key  = "HubUserId" #attribute used as partition key (beware of hot partition problem)
  range_key = "HubSearchedItemId"           #attribute used as sort key
  global_secondary_index {
    name            = "ItemSearchedByUserGSI"
    hash_key        = "HubSearchedItemId"
    range_key       = "HubUserId"
    projection_type = "ALL"
  }
  ttl {
    attribute_name = "TimeExpired"
    enabled        = true
  }
  #other fields: check DAO on app-be code base

  #you can only specify indexed attributes (pk, sk, lsi, gsi)
  attribute {
    name = "HubUserId"
    type = "S" #(S)tring, (N)umber or (B)inary
  }
  attribute {
    name = "HubSearchedItemId"
    type = "S" #(S)tring, (N)umber or (B)inary
  }
}
