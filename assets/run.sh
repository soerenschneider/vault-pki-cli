#!/bin/sh

set -e

SLEEP=0
TIMEOUT=120

if [ $(kind get clusters | wc -l) -eq 0 ]; then
  echo no cluster;
  cd kind && sh run-cluster.sh && cd -
  rm -f terraform/*.tfstate*
  SLEEP=15
fi

if kubectl config current-context | grep -e "^kind"; then
  echo "already using kind cluster"
else
  echo "please point kubectl to right cluster"
  exit 1
fi

kubectl apply -f k8s
sleep ${SLEEP}

echo ""
echo "waiting for vault to be ready... (timeout ${TIMEOUT}s)"
kubectl wait pod -l app=vault --for condition=Ready --timeout=${TIMEOUT}s

cd terraform && sh run-terraform.sh && cd -

cd .. && skaffold run && cd -

sleep 2
kubectl logs --tail=50 -f -l app=vault-pki-cli