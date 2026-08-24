variable "project_name" {
  type        = string
  description = "The name to be used for resources and tags"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to assign to resources"
  default     = {}
}

variable "public_subnet_1a" {
  type        = string
  description = "Subnet to create EKS cluster AZ 1a"
}

variable "public_subnet_1b" {
  type        = string
  description = "Subnet to create EKS cluster AZ 1b"
}

