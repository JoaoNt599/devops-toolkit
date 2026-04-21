#!/bin/bash

# ==============================================================================
# Script para Diagnóstico de Conexões TLS/HTTPS
# ==============================================================================

# Definição de variáveis (podem ser passadas como parâmetro ou editadas aqui)
HOST="${1:-www.google.com}"
PORT="${2:-443}"
CERT="/tmp/serpro.crt"

# Criação do arquivo de log baseado na data e hora (ex: tls_check_20251225_112305.txt)
LOG_FILE="tls_check_$(date +%Y%m%d_%H%M%S).txt"

# Redireciona toda a saída (stdout e stderr) para o terminal E para o arquivo de log
exec > >(tee -i "$LOG_FILE") 2>&1

echo "########################################################"
echo " TLS / HTTPS CHECK"
echo " Host: $HOST"
echo " Date: $(date)"
echo "########################################################"
echo

# 1. Inspecionar a cadeia completa de certificados e validar handshake
echo ">>> OpenSSL - s_client (showcerts)"
# O redirecionamento < /dev/null evita que o openssl fique aguardando input do usuário
openssl s_client -showcerts -connect "$HOST:$PORT" < /dev/null
echo

# 2. Teste de conectividade simples via cURL
echo ">>> CURL (simple)"
# Adicionado -s -o /dev/null -w para testar a resposta HTTP de forma mais limpa, 
# mas mantendo a essência do seu script original que chamava direto.
curl -I "https://$HOST"
echo

# 3. Teste de conectividade detalhada via cURL (Verbose para ver detalhes do handshake)
echo ">>> CURL (verbose)"
curl -v "https://$HOST" > /dev/null
echo

# 4. Verificação de confiança do Sistema Operacional (CA Store)
echo ">>> VERIFICAÇÃO DE TRUST STORE LOCAIS"
if [ -f "$CERT" ]; then
    echo ">>> Certificado encontrado: $CERT"
    
    # Copia o certificado para a pasta de CAs do sistema (Padrão Debian/Ubuntu)
    sudo cp "$CERT" /usr/local/share/ca-certificates/
    
    # Atualiza a store do sistema
    sudo update-ca-certificates
else
    echo ">>> Certificado $CERT NÃO encontrado"
fi

echo
echo "########################################################"
echo ">>> Diagnóstico concluído. Evidências salvas em: $LOG_FILE"
echo "########################################################"