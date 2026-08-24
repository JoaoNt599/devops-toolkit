data "tls_certificate" "eks_oidc_tls_certificate" {
  url = aws_eks_cluster.eks_cluster.identity[0].oidc[0].issuer
}

resource "aws_aim_open_id_connect_provider" "eks_oidc_provider" {
  client_id_list  = [
    "sts.amazonaws.com"
    ]
  thumbprint_list = [data.tls_certificate.eks_oidc_tls_certificate.certificates[*].sha1_fingerprint]
  url             = data.tls_certificate.eks_oidc_tls_certificate.certificates.url

  tags = merge(
    local.tags,
    {
      Name = "${var.project_name}-oidc-provider"
    }
  )
}