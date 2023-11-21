output "vpc_ids" {
  value = data.aws_vpcs.list.ids

  #optional arguments
  description = "The list of VPCs created from network-module."
  # sensitive   = false
  # depends_on = [
  #   #When a parent module accesses an output value exported by one of its child modules,
  #   # the dependencies of that output value allow Terraform to correctly determine
  #   # the dependencies between resources defined in different modules.
  #   #Terraform analyzes the value expression for an output value
  #   # and automatically determines a set of dependencies,
  #   # but in less-common cases there are dependencies that cannot be recognized implicitly.
  #   #In these rare cases, the depends_on argument can be used to create additional explicit dependencies.
  # ]

  # precondition {
  #   condition     = length(data.aws_vpcs.list.ids) > 0
  #   error_message = "The network-module VPC is undetected."
  # }
}
