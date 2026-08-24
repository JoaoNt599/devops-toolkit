locals {
  oidc = split("/", var.oidc_provider_arn)[4]
}