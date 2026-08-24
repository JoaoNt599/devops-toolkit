resource "kubernetes_service_account" "eks_alb_controller" {
  metadata {
    name      = "aws-load-balancer-controller"
    namespace =     "kube-system"
    labels = {
      "app.kubernetes.io/name"       = "aws-load-balancer-controller"
      "app.kubernetes.io/managed-by" = "terraform"
    }
    annotations = {
      "eks.amazonaws.com/role-arn" = module.iam_role.eks_alb_controller_role.arn
    }
  }
}