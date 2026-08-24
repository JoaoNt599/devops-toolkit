module "eks_network" {
  source       = "./modules/network"
  cidr_block   = var.cidr_block
  project_name = var.project_name
  tags         = var.tags
}

module "eks_cluster" {
  source           = "./modules/cluster"
  project_name     = var.project_name
  tags             = var.tags
  public_subnet_1a = module.eks_network.public_subnet_1a
  public_subnet_1b = module.eks_network.public_subnet_1b
}

module "eks_managed_node_group" {
  source            = "./modules/managed-node-group"
  project_name      = var.project_name
  tags              = var.tags
  cluster_name      = module.eks_cluster.cluster_name
  private_subnet_1a = module.network.private_subnet_1a
  private_subnet_1b = module.network.private_subnet_1b
}

module "eks_alb" {
  source            = "./modules/alb"
  project_name      = var.project_name
  tags              = var.tags
  oidc_provider_arn = module.eks_cluster.oidc
  eks_cluster_name  = module.eks_cluster.cluster_name
}