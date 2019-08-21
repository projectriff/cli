#!/bin/bash

set -o nounset

fats_dir=`dirname "${BASH_SOURCE[0]}"`/fats

uninstall_chart() {
  local name=$1

  helm delete --purge $name
  kubectl delete customresourcedefinitions.apiextensions.k8s.io -l app.kubernetes.io/managed-by=Tiller,app.kubernetes.io/instance=$name 
}

# attempt to cleanup fats
if [ -d "$fats_dir" ]; then
  source $fats_dir/macros/cleanup-user-resources.sh
  kubectl delete namespace $NAMESPACE

  echo "Uninstall riff system"
  
  uninstall_chart riff

  uninstall_chart istio
  kubectl get customresourcedefinitions.apiextensions.k8s.io -oname | grep istio.io | xargs -L1 kubectl delete
  kubectl delete namespace istio-system

  uninstall_chart cert-manager

  source $fats_dir/macros/helm-reset.sh
  source $fats_dir/cleanup.sh
fi
