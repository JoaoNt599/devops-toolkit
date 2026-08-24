variable "project_name" {
  type = string
  description = "Project name to be used to name the resources"
}

variable "tags" {
  type = map(any)
  description = "Tags to be added to AWS resources"
  default = {}
}

variable "cluster_name" {
  type = string
  description = "EKS cluster name to associate the managed node group with"
}

variable "private_subnet_1a" {
  type = string
  description = "Private subnet ID for the managed node group"
}

variable "private_subnet_1b" {
  type = string
  description = "Private subnet ID for the managed node group"
}