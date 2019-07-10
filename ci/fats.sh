#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

version=`cat VERSION`
commit=$(git rev-parse HEAD)

# fetch FATS scripts
fats_dir=`dirname "${BASH_SOURCE[0]}"`/fats
fats_repo="projectriff/fats"
fats_refspec=1a462b4a35ebdb66d678d648b5803fa498871ec6 # projectriff/fats master as of 2019-07-09
source `dirname "${BASH_SOURCE[0]}"`/fats-fetch.sh $fats_dir $fats_refspec $fats_repo
source $fats_dir/.util.sh


# install riff-cli
echo "Installing riff"
if [ "$machine" == "MinGw" ]; then
  curl https://storage.googleapis.com/projectriff/riff-cli/releases/builds/v${version}-${commit}/riff-windows-amd64.zip > riff.zip
  unzip riff.zip -d /usr/bin/
  rm riff.zip
else
  curl https://storage.googleapis.com/projectriff/riff-cli/releases/builds/v${version}-${commit}/riff-linux-amd64.tgz | tar xz
  chmod +x riff
  sudo cp riff /usr/bin/riff
fi

# start FATS
source $fats_dir/start.sh

# install riff system
echo "Installing riff system"
$fats_dir/install.sh helm

kubectl create serviceaccount tiller -n kube-system
kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount kube-system:tiller
helm init --service-account tiller

helm repo add projectriff https://projectriff.storage.googleapis.com/charts/releases
helm repo update

helm install projectriff/projectriff-riff --name riff --set istio.enabled=true --set global.k8s.service.type=${K8S_SERVICE_TYPE} --devel

# health checks
echo "Checking for ready ingress"
wait_for_ingress_ready 'istio-ingressgateway' 'istio-system'

# setup namespace
kubectl create namespace $NAMESPACE
fats_create_push_credentials $NAMESPACE

# run test functions
source $fats_dir/functions/helpers.sh

for test in command; do
  path=${fats_dir}/functions/uppercase/${test}
  function_name=fats-cluster-uppercase-${test}
  image=$(fats_image_repo ${function_name})
  create_args="--git-repo https://github.com/${fats_repo}.git --git-revision ${fats_refspec} --sub-path functions/uppercase/${test}"
  input_data=riff
  expected_data=RIFF

  run_function $path $function_name $image "${create_args}" $input_data $expected_data
done

if [ "$machine" != "MinGw" ]; then
  for test in command; do
    path=${fats_dir}/functions/uppercase/${test}
    function_name=fats-local-uppercase-${test}
    image=$(fats_image_repo ${function_name})
    create_args="--local-path ."
    input_data=riff
    expected_data=RIFF

    run_function $path $function_name $image "${create_args}" $input_data $expected_data
  done
fi
