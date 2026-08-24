output "subnet_public_1a" {
  value = aws_subnet.eks_subnet_public_1a.id
}

output "subnet_public_1b" {
  value = aws_subnet.eks_subnet_public_1b.id
}

output "subnet_private_1a" {
  value = aws_subnet.eks_subnet_private_1a.id
}

output "subnet_private_1b" {
  value = aws_subnet.eks_subnet_private_1b.id
}