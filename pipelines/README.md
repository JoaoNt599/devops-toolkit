# Pipelines Templates


## Objetivo: Criar um kit reutilizável, organizado e explicável que permita:

- Padronizar processos de deploy
- Melhorar rastreabilidade de versões
- Reduzir riscos em deploys
- Facilitar rollback em caso de falha
- Servir como base para novos projetos

## Princípios adotados

- Rastreabilidade: cada deploy deve ser identificável
- Imutabilidade: evitar sobrescrita de artefatos
- Automação: reduzir intervenção manual
- Resiliência: rollback automático em falhas
- Clareza: tudo deve ser explicável e reutilizáve


## O que já está implementado

CI/CD (GitHub Actions)

Pipeline automatizada com:

- Build de imagem Docker
- Versionamento com SHA (imutabilidade)
- Push para ECR
- Deploy no ECS
- Espera por estabilidade do serviço
- Health check após deploy
- Rollback automático em caso de falha


## Problema resolvido:

- Deploy sem rastreabilidade e com alto risco de falha sem recuperação automática.
