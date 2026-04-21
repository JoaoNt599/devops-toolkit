# TLS/HTTPS Diagnostic Tool

Este repositório contém um script Shell focado em diagnosticar e validar conexões TLS/HTTPS. Ele foi projetado para identificar falhas silenciosas de segurança e conectividade que ocorrem em tempo de execução (*runtime*), ajudando a isolar problemas de infraestrutura antes que afetem a produção.

## Descrição do Projeto

Erros de TLS/SSL costumam ser vagos em logs de aplicação. Este script atua como uma ferramenta de diagnóstico de baixo nível para:
- Validar o handshake TLS completo.
- Inspecionar a cadeia de certificados enviada pelo servidor.
- Verificar a integridade e confiança da *Trust Store* do sistema operacional.
- Automatizar a atualização de CAs (Autoridades Certificadoras) locais.

## Funcionalidades

- **Handshake Check:** Valida se o servidor suporta os protocolos e cifras esperados.
- **Certificate Inspection:** Exibe detalhadamente o `Subject`, `Issuer` e a cadeia completa (Root, Intermediate e Leaf).
- **Update CA Store:** Possui lógica integrada para detectar certificados locais (ex: `serpro.crt`) e atualizar o repositório de confiança do Linux automaticamente.
- **Logging de Evidências:** Cada execução gera um arquivo `.txt` único com carimbo de data e hora para fins de auditoria e depuração.

## Tecnologias Utilizadas

- **Shell Script (Bash):** Orquestração e lógica de automação.
- **OpenSSL (s_client):** Ferramenta principal para diagnóstico de SSL/TLS.
- **cURL:** Validação de requisições HTTPS e negociação de cifras.
- **ca-certificates:** Utilitário nativo do Linux para gerenciamento de certificados confiáveis.

## Como Usar

### 1. Preparação
Dê permissão de execução ao arquivo:
```bash
chmod +x tls_check.sh

### 2. Execução Simples
O script aceita o Host e a Porta como argumentos (padrão: www.google.com na porta 443):

    ./tls_check.sh api.meudominio.com 443

### 3. Atualização de Certificados

Se você possuir um arquivo de certificado em /tmp/serpro.crt, o script tentará automaticamente copiá-lo para a pasta de CAs do sistema e rodar o comando update-ca-certificates (requer privilégios de sudo).

## Estrutura de Saída

A cada execução, o script gera um log com o seguinte padrão de nomenclatura:
tls_check_YYYYMMDD_HHMMSS.txt

Este log contém o dump do OpenSSL, o cabeçalho HTTP retornado pelo cURL e o status da atualização da Trust Store.


## Problemas Comuns Resolvidos

- Certificados Expirados: Identificação imediata na saída do OpenSSL.
- Cadeia Incompleta: Quando o servidor esquece de enviar o certificado intermediário.
- Self-Signed sem Trust: Ajuda a validar se o certificado de uma infraestrutura interna foi corretamente instalado no SO.
- Incompatibilidade de Cipher: Verifica se o cliente consegue negociar cifras modernas (como AES-GCM).


####  Script desenvolvido para automação de rotinas de DevOps, SecOps e SRE.