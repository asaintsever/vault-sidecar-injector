#!/bin/bash

set -e

SCRIPT_PATH="$(cd "$(dirname $0)" && pwd)"
VAULT_POD="kubectl exec -i vault-0 -- sh -c"

# Create policies to allow read access to our secrets
cat "$SCRIPT_PATH"/vault-server-policy-test.hcl | ${VAULT_POD} "VAULT_TOKEN=root vault policy write test_pol -"
cat "$SCRIPT_PATH"/vault-server-policy-test2.hcl | ${VAULT_POD} "VAULT_TOKEN=root vault policy write test_pol2 -"
cat "$SCRIPT_PATH"/vault-server-policy-test3.hcl | ${VAULT_POD} "VAULT_TOKEN=root vault policy write test_pol3 -"

# Enable KV v1 (KV v2 is enabled with Vault server in dev mode whereas KV v1 is enabled in prod mode: https://www.vaultproject.io/docs/secrets/kv/kv-v2.html#setup)
${VAULT_POD} "VAULT_TOKEN=root vault secrets disable secret/"
${VAULT_POD} "VAULT_TOKEN=root vault secrets enable -version=1 -path=secret kv"

# Enable Transit Secrets Engine and create a test key
${VAULT_POD} "VAULT_TOKEN=root vault secrets enable transit" || true
${VAULT_POD} "VAULT_TOKEN=root vault write -f transit/keys/test-key"

# Enable Vault K8S Auth Method
echo "-> Enable & set up Vault Kubernetes Auth Method"
${VAULT_POD} "VAULT_TOKEN=root vault auth enable kubernetes" || true

# Config Vault K8S Auth Method
export VAULT_SA_NAME=$(kubectl get sa vault -o jsonpath="{.secrets[*]['name']}")
export SA_JWT_TOKEN=$(kubectl get secret $VAULT_SA_NAME -o jsonpath="{.data.token}" | base64 --decode; echo)
export SA_CA_CRT=$(kubectl get secret $VAULT_SA_NAME -o jsonpath="{.data['ca\.crt']}" | base64 --decode; echo)

K8S_VER_MAJOR=$(kubectl version --short -o json | jq -r '.serverVersion.major')
K8S_VER_MINOR=$(kubectl version --short -o json | jq -r '.serverVersion.minor')

if [ $K8S_VER_MAJOR -ge 1 ] && [ $K8S_VER_MINOR -gt 20 ];then
    echo "Kubernetes 1.21+: get service account issuer"
    # See ref: https://www.vaultproject.io/docs/auth/kubernetes#discovering-the-service-account-issuer
    kubectl proxy &
    echo "Wait ..."
    sleep 10
    export SA_ISSUER=$(curl -s http://127.0.0.1:8001/.well-known/openid-configuration | jq -r .issuer)
    echo "Get issuer for cluster: $SA_ISSUER"

    ${VAULT_POD} "VAULT_TOKEN=root vault write auth/kubernetes/config kubernetes_host=\"https://kubernetes:443\" kubernetes_ca_cert=\"$SA_CA_CRT\" token_reviewer_jwt=\"$SA_JWT_TOKEN\" issuer=\"$SA_ISSUER\""
else
    ${VAULT_POD} "VAULT_TOKEN=root vault write auth/kubernetes/config kubernetes_host=\"https://kubernetes:443\" kubernetes_ca_cert=\"$SA_CA_CRT\" token_reviewer_jwt=\"$SA_JWT_TOKEN\""
fi

# Create roles for Vault K8S Auth Method
${VAULT_POD} "VAULT_TOKEN=root vault write auth/kubernetes/role/test bound_service_account_names=default,job-sa bound_service_account_namespaces=default policies=test_pol ttl=5m"
${VAULT_POD} "VAULT_TOKEN=root vault write auth/kubernetes/role/test2 bound_service_account_names=default,job-sa bound_service_account_namespaces=default policies=test_pol2 ttl=5m"
${VAULT_POD} "VAULT_TOKEN=root vault write auth/kubernetes/role/test3 bound_service_account_names=default,job-sa bound_service_account_namespaces=default policies=test_pol3 ttl=5m"

# Enable Vault AppRole Auth Method
echo "-> Enable & set up Vault AppRole Auth Method"
${VAULT_POD} "VAULT_TOKEN=root vault auth enable approle" || true

# Create roles for Vault AppRole Auth Method
${VAULT_POD} "VAULT_TOKEN=root vault write auth/approle/role/test secret_id_ttl=60m token_num_uses=0 token_ttl=20m token_max_ttl=30m secret_id_num_uses=0 policies=test_pol"
${VAULT_POD} "VAULT_TOKEN=root vault write auth/approle/role/test2 secret_id_ttl=60m token_num_uses=0 token_ttl=20m token_max_ttl=30m secret_id_num_uses=0 policies=test_pol2"
${VAULT_POD} "VAULT_TOKEN=root vault write auth/approle/role/test3 secret_id_ttl=60m token_num_uses=0 token_ttl=20m token_max_ttl=30m secret_id_num_uses=0 policies=test_pol3"

# Add some secrets
${VAULT_POD} "VAULT_TOKEN=root vault kv put secret/test/test-app-svc ttl=10s SECRET1=Batman SECRET2=BruceWayne"
${VAULT_POD} "VAULT_TOKEN=root vault kv put secret/test2/test-app2-svc ttl=5s SECRET1=my SECRET2=name SECRET3=is SECRET4=James"
${VAULT_POD} "VAULT_TOKEN=root vault kv put secret/test3/test-app3-svc1 ttl=15s SECRET1=svc1_sec1_value SECRET2=svc1_sec2_value"
${VAULT_POD} "VAULT_TOKEN=root vault kv put secret/test3/test-app3-svc2 ttl=15s SECRET1=svc2_sec1_value SECRET2=svc2_sec2_value"

# List Auth Methods and Secrets Engines
echo
echo "Auth Methods"
echo "============"
${VAULT_POD} "VAULT_TOKEN=root vault auth list -detailed"
echo
echo "Secrets Engines"
echo "==============="
${VAULT_POD} "VAULT_TOKEN=root vault secrets list -detailed"