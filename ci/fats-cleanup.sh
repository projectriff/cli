#!/bin/bash

set -o nounset

fats_dir=`dirname "${BASH_SOURCE[0]}"`/fats

# attempt to cleanup fats
if [ -d "$fats_dir" ]; then
  source $fats_dir/macros/cleanup-user-resources.sh
  kubectl delete namespace $NAMESPACE

  echo "Uninstall riff system"
  
  helm delete --purge riff
  kubectl delete customresourcedefinitions.apiextensions.k8s.io -l app.kubernetes.io/managed-by=Tiller,app.kubernetes.io/instance=riff

  helm delete --purge istio
  kubectl delete customresourcedefinitions.apiextensions.k8s.io -l app.kubernetes.io/managed-by=Tiller,app.kubernetes.io/instance=istio
  kubectl delete namespace istio-system

  source $fats_dir/macros/helm-reset.sh
  source $fats_dir/cleanup.sh
fi
