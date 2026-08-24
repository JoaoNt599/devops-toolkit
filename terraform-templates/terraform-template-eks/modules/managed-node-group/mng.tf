resource "aws_eks_node_group" "eks_managed_node_group" {
  cluster_name = var.cluster_name
  node_group_name = "${var.project_name}-managed-node-group"
  node_role_arn = aws_iam_role.eks_mng_role.arn
  subnet_ids = [
    var.private_subnet_1a,
    var.private_subnet_1b
  ]

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  tags = merge(
    var.tags,
    {
      "Name" = "${var.project_name}-managed-node-group"
    }
  )

    depends_on = [
    aws_iam_role_policy_attachment.eks_mng_role_attachment_worker,
    aws_iam_role_policy_attachment.eks_mng_role_attachment_ecr,
    aws_iam_role_policy_attachment.eks_mng_role_attachment_cni
    ]
}