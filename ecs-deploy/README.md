# ECS Automated Deployment Script

Este repositório contém o script de automação de deploy (deploy.sh) desenvolvido para simplificar e padronizar a entrega de aplicações em containers no Amazon Elastic Container Service (ECS) via Amazon Elastic Container Registry (ECR).

## Visão Geral

O script automatiza o ciclo completo de build e deploy de uma imagem Docker, executando as seguintes etapas de forma sequencial e segura:

1. Autenticação segura no AWS ECR.
2. Build da imagem Docker local.
3. Aplicação da estratégia de Dupla Tag (Double Tagging).
4. Push das imagens para o repositório ECR.
5. Gatilho de atualização forçada (Force New Deployment) no cluster ECS.

## A Estratégia de Dupla Tag (Double Tagging)

Para garantir um histórico de versionamento seguro sem a necessidade de alterar a Task Definition do ECS a cada deploy manual, o script aplica duas tags à mesma imagem gerada:

- **Tag de Histórico (ex: qa-v1.0.0):** Utilizada para manter o histórico imutável no ECR, permitindo auditoria e facilitando o processo de rollback em caso de falhas.
- **Tag de Target do ECS (ex: latest):** É a tag que a Task Definition atual do ECS está configurada para buscar. Ao sobrescrever esta tag, o ECS automaticamente puxa o código atualizado durante o Force New Deployment.

## Pré-requisitos

Antes de executar o script, certifique-se de que o seu ambiente atende aos seguintes requisitos:

1. **Docker:** Instalado e com o daemon em execução. O usuário logado deve ter permissão de execução do Docker (pertencer ao grupo docker).
2. **AWS CLI:** Instalado (versão 2 recomendada).
3. **Credenciais AWS:** Configuradas localmente via `aws configure` ou via variáveis de ambiente, com um perfil (IAM User/Role) que possua as seguintes permissões:
   - `ecr:GetAuthorizationToken`
   - `ecr:BatchCheckLayerAvailability`, `ecr:PutImage`, `ecr:InitiateLayerUpload`, etc. (Permissões de Push no ECR).
   - `ecs:UpdateService` (Permissão no cluster e serviço alvo).

## Configuração

Antes do primeiro uso, edite o arquivo `deploy.sh` e valide/configure as variáveis globais de ambiente de acordo com a sua infraestrutura:

```bash
REGION="us-east-1"
REPO="seu-repositorio/nome-da-imagem"
CLUSTER_NAME="seu-cluster"
SERVICE_NAME="seu-servico"
ECS_TARGET_TAG="latest"

(Nota: O ACCOUNT_ID é capturado dinamicamente com base nas credenciais logadas).

## Como utilizar

1. Clone o repositório e navegue até a raiz do projeto (onde encontra-se o Dockerfile).

2. Conceda permissão de execução ao script (necessário apenas na primeira vez):

    chmod +x deploy.sh

3. Execute o script:

    ./deploy.sh

4. Siga as instruções no terminal. O script solicitará a tag de versionamento (ex: v1.0.0) e pedirá uma confirmação final antes de iniciar o processo de build e push.

## Troubleshooting (Solução de Problemas)

- Erro: invalid reference format: Certifique-se de que o nome do seu repositório ECR não começa com uma barra (/) e não contém letras maiúsculas.

- Erro: permission denied while trying to connect to the docker API: O seu usuário não tem permissões para executar o Docker. Execute sudo usermod -aG docker $USER e reinicie a sessão do terminal, ou execute o script com sudo (garantindo que as credenciais da AWS estejam acessíveis para o usuário root).

- Erro de Autenticação no ECR: Verifique se o seu token de sessão (MFA) expirou ou se as chaves configuradas no aws configure possuem as permissões corretas.