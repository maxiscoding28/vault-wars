helm repo add hashicorp https://helm.releases.hashicorp.com
helm search repo hashicorp/vault
helm install luke hashicorp/vault --values yamls/luke.yaml




kubectl exec luke-vault-0 -- vault operator init \
    -key-shares=1 \
    -key-threshold=1 \
    -format=json > cluster-keys.json

VAULT_UNSEAL_KEY=$(jq -r ".unseal_keys_b64[]" cluster-keys.json)
kubectl exec luke-vault-0 -- vault operator unseal $VAULT_UNSEAL_KEY


kubectl exec luke-vault-1 -- vault operator unseal $VAULT_UNSEAL_KEY
kubectl exec luke-vault-2 -- vault operator unseal $VAULT_UNSEAL_KEY


kubectl exec -ti luke-vault-0 -- vault status
kubectl exec -ti luke-vault-1 -- vault status
kubectl exec -ti luke-vault-2 -- vault status

helm uninstall luke
k delete pvc --all




# Debug AWS CLI
kubectl run --rm -it aws-cli --image=amazon/aws-cli --command -- sh
