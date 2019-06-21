#!/bin/bash

set -o nounset

fats_dir=`dirname "${BASH_SOURCE[0]}"`/fats

# attempt to cleanup fats
if [ -d "$fats_dir" ]; then
  echo "Uninstall riff system"
  duffle_k8s_service_account=${duffle_k8s_service_account:-duffle-runtime}
  duffle_k8s_namespace=${duffle_k8s_namespace:-kube-system}
  SERVICE_ACCOUNT=${duffle_k8s_service_account} KUBE_NAMESPACE=${duffle_k8s_namespace} duffle uninstall riff -d k8s || true
  kubectl delete namespace $NAMESPACE || true

  source $fats_dir/cleanup.sh
fi
