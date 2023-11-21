data "aws_vpcs" "list" {
  tags = {
    deployment_environment_name = var.deployment_environment_name
    project_name                = var.project_name
  }
}
