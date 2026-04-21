#!/bin/bash

set -e

# ==============================
# Validação de dependências
# ==============================

check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo "Erro: '$1' não está instalado."
        return 1
    fi
}

echo "Verificando dependências..."

check_command aws || {
    echo "Instale o AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html"
    exit 1
}

check_command docker || {
    echo "Instale o Docker: https://docs.docker.com/get-docker/"
    exit 1
}

echo "Dependências OK"
echo ""

# ==============================
# Variáveis configuráveis
# ==============================

REGION=${REGION:-"sua_regiao"}
ACCOUNT_ID=${ACCOUNT_ID:-"sua_conta_id"}
REPO=${REPO:-"seu_repositorio"}

ECR_URL="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

CLUSTER_NAME=${CLUSTER_NAME:-"seu_cluster_name"}
SERVICE_NAME=${SERVICE_NAME:-"seu_service_name"}

# ==============================
# Tag
# ==============================

echo "======================================"
echo "   Deploy Automatizado (QA - ECS)"
echo "======================================"
echo "Exemplo de tags: latest, v1, test"
echo "======================================"

read -p "Digite a tag da imagem: " TAG

if [ -z "$TAG" ]; then
    echo "Tag não pode ser vazia."
    exit 1
fi

echo ""
echo "Iniciando deploy..."
echo "Tag: $TAG"
echo "Cluster: $CLUSTER_NAME"
echo "Service: $SERVICE_NAME"
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
# Build
# ==============================

echo "Build da imagem..."
docker build -t $REPO:$TAG .

# ==============================
# Tag
# ==============================

echo "Tag da imagem..."
docker tag $REPO:$TAG $ECR_URL/$REPO:$TAG

# ==============================
# Push
# ==============================

echo "Enviando para o ECR..."
docker push $ECR_URL/$REPO:$TAG

# ==============================
# Deploy ECS
# ==============================

echo "Atualizando serviço ECS..."

aws ecs update-service \
    --cluster $CLUSTER_NAME \
    --service $SERVICE_NAME \
    --force-new-deployment \
    --region $REGION

echo ""
echo "Deploy da versão $TAG iniciado com sucesso no ambiente de QA!"
echo "Acompanhe a subida dos containers pelo painel da AWS."