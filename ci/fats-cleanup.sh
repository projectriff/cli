#!/bin/bash

set -o nounset

fats_dir=`dirname "${BASH_SOURCE[0]}"`/fats

# attempt to cleanup fats
if [ -d "$fats_dir" ]; then
  echo "Uninstall riff system"
  kubectl delete riff --all-namespaces --all
  kubectl delete knative --all-namespaces --all

  kubectl delete namespace $NAMESPACE
  
  helm delete --purge riff
  kubectl delete customresourcedefinitions.apiextensions.k8s.io -l app.kubernetes.io/managed-by=Tiller,app.kubernetes.io/instance=riff

  helm delete --purge istio
  kubectl delete customresourcedefinitions.apiextensions.k8s.io -l app.kubernetes.io/managed-by=Tiller,app.kubernetes.io/instance=istio
  kubectl delete namespace istio-system

  helm reset
  kubectl delete serviceaccount tiller -n kube-system
  kubectl delete clusterrolebinding tiller

  source $fats_dir/cleanup.sh
fi
