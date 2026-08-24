output "oidc" {
  value = modules.eks_cluster.oidc
}

output "ca" {
  value = modules.eks_cluster.certificate_authority
}

output "endpoint" {
  value = modules.eks_cluster.eks_cluster_endpoint
}