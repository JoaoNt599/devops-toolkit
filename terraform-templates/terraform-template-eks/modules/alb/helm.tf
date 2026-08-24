resource "helm_release" "eks_helm_alb_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"
  version   = "1.4.4"

  set = [
    {
      name  = "clusterName"
      value = module.eks_cluster.eks_cluster_name
    },
    {
      name  = "serviceAccount.create"
      value = "false"
    },
    {
      name  = "serviceAccount.name"
      value = kubernetes_service_account.eks_alb_controller.metadata[0].name
    }
  ]
}