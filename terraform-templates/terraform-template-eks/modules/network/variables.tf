variable "cidr_block" {
  type        = string
  description = "The CIDR block for the subnet (VPC)"
}

variable "project_name" {
  type        = string
  description = "The name to be used for resources and tags"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to assign to resources"
  default     = {}
}