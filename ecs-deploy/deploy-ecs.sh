#!/bin/bash

set -e

# ==============================
# Validação de dependências
# ==============================

check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo " Erro: '$1' não está instalado."
        return 1
    fi
}

echo "Verificando dependências..."
check_command aws || { echo " Instale o AWS CLI"; exit 1; }
check_command docker || { echo "Instale o Docker"; exit 1; }
echo "Dependências OK"
echo ""

# ==============================
# Variáveis configuráveis
# ==============================

REGION=${REGION:-"sua-região"}
# Pega o Account ID automaticamente via AWS CLI
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
REPO=${REPO:-"seu-repo-name"}

ECR_URL="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

CLUSTER_NAME=${CLUSTER_NAME:-"seu-cluster-name"}
SERVICE_NAME=${SERVICE_NAME:-"seu-service-name"}

# Tag fixa que o ECS está esperando atualmente
ECS_TARGET_TAG="latest"

# ==============================
# Input da Tag
# ==============================

echo "======================================"
echo "    Deploy Automatizado (QA - ECS)"
echo "======================================"
echo "Use tags de versionamento: qa-v1.0.0, qa-v1.0.1..."
echo "======================================"

read -p "Digite a tag da imagem: " TAG

if [ -z "$TAG" ]; then
    echo " Erro: Tag não pode ser vazia."
    exit 1
fi

echo ""
echo "Iniciando deploy..."
echo "Tag Histórico: $TAG"
echo "Tag ECS: $ECS_TARGET_TAG"
echo "Cluster: $CLUSTER_NAME"
echo ""

read -p "Deseja continuar? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Deploy cancelado."
    exit 1
fi

# ==============================
# Login no ECR
# ==============================

echo "Autenticando no ECR..."
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ECR_URL

# ==============================
# Build & Tags
# ==============================

echo "Fazendo o build da imagem..."
docker build -t $REPO:$TAG .

echo "Aplicando a Dupla Tag..."
# Tag de histórico (ex: qa-v1.0.0)
docker tag $REPO:$TAG $ECR_URL/$REPO:$TAG
# Tag que o ECS vai puxar (latest)
docker tag $REPO:$TAG $ECR_URL/$REPO:$ECS_TARGET_TAG

# ==============================
# Push
# ==============================

echo "Enviando a versão $TAG para o ECR..."
docker push $ECR_URL/$REPO:$TAG

echo "Atualizando a tag $ECS_TARGET_TAG no ECR..."
docker push $ECR_URL/$REPO:$ECS_TARGET_TAG

# ==============================
# Deploy ECS
# ==============================

echo "Forçando o ECS a puxar a nova imagem..."
aws ecs update-service \
    --cluster $CLUSTER_NAME \
    --service $SERVICE_NAME \
    --force-new-deployment \
    --region $REGION > /dev/null

echo ""
echo "Deploy da versão $TAG iniciado com sucesso no ambiente de QA!"
echo "Acompanhe a subida dos containers pelo painel da AWS."