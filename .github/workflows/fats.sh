#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly version=$(cat VERSION)
readonly git_sha=$(git rev-parse HEAD)
readonly git_timestamp=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d%H%M%S' --format="%cd")
readonly slug=${version}-${git_timestamp}-${git_sha:0:16}

# fetch FATS scripts
fats_dir=`dirname "${BASH_SOURCE[0]}"`/fats
fats_repo="projectriff/fats"
fats_refspec=696bfb86ab8111c1945b81a661629d7dd70388e7 # master as of 2019-12-08
source `dirname "${BASH_SOURCE[0]}"`/fats-fetch.sh $fats_dir $fats_refspec $fats_repo
source $fats_dir/.util.sh

# install riff-cli
echo "Installing riff"
if [ "$machine" == "MinGw" ]; then
  curl https://storage.googleapis.com/projectriff/riff-cli/releases/builds/v${slug}/riff-windows-amd64.zip > riff.zip
  unzip riff.zip -d /usr/bin/
  rm riff.zip
else
  curl https://storage.googleapis.com/projectriff/riff-cli/releases/builds/v${slug}/riff-linux-amd64.tgz | tar xz
  chmod +x riff
  sudo cp riff /usr/bin/riff
fi

# start FATS
source $fats_dir/start.sh

$fats_dir/install.sh helm
source $fats_dir/macros/helm-init.sh
helm repo add projectriff https://projectriff.storage.googleapis.com/charts/releases
helm repo update

echo "Installing Cert Manager"
helm install projectriff/cert-manager --name cert-manager --devel --wait

source $fats_dir/macros/no-resource-requests.sh

echo "Installing riff system"
helm install projectriff/istio --name istio --namespace istio-system --devel --wait \
  --set gateways.istio-ingressgateway.type=${K8S_SERVICE_TYPE}
helm install projectriff/riff --name riff --devel --wait \
  --set tags.core-runtime=true \
  --set tags.knative-runtime=true \
  --set cert-manager.enabled=false

# health checks
echo "Checking for ready ingress"
wait_for_ingress_ready 'istio-ingressgateway' 'istio-system'

# setup namespace
kubectl create namespace $NAMESPACE
fats_create_push_credentials $NAMESPACE
source $fats_dir/macros/create-riff-dev-pod.sh

# run test functions
for test in command; do
  name=fats-cluster-uppercase-${test}
  image=$(fats_image_repo ${name})
  curl_opts="-H Content-Type:text/plain -H Accept:text/plain -d cli"
  expected_data="CLI"

  echo "##[group]Run function $name"

  riff function create $name --image $image --namespace $NAMESPACE --tail \
    --git-repo https://github.com/$fats_repo --git-revision $fats_refspec --sub-path functions/uppercase/${test} &

  riff core deployer create $name \
    --function-ref $name \
    --ingress-policy External \
    --namespace $NAMESPACE \
    --tail
  source $fats_dir/macros/invoke_incluster.sh \
    "$(kubectl get deployers.core.projectriff.io ${name} --namespace ${NAMESPACE} -ojsonpath='{.status.address.url}')" \
    "${curl_opts}" \
    "${expected_data}"
  # TODO invoke via ingress as well
  riff core deployer delete $name --namespace $NAMESPACE

  riff knative deployer create $name \
    --function-ref $name \
    --ingress-policy External \
    --namespace $NAMESPACE \
    --tail
  source $fats_dir/macros/invoke_incluster.sh \
    "$(kubectl get deployers.knative.projectriff.io ${name} --namespace ${NAMESPACE} -ojsonpath='{.status.address.url}')" \
    "${curl_opts}" \
    "${expected_data}"
  source $fats_dir/macros/invoke_knative_deployer.sh \
    "${name}" \
    "${curl_opts}" \
    "${expected_data}"
  riff knative deployer delete $name --namespace $NAMESPACE

  riff function delete $name --namespace $NAMESPACE
  fats_delete_image $image

  echo "##[endgroup]"
done

if [ "$machine" != "MinGw" ]; then
  for test in command; do
    name=fats-local-uppercase-${test}
    image=$(fats_image_repo ${name})
    curl_opts="-H Content-Type:text/plain -H Accept:text/plain -d cli"
    expected_data="CLI"

    echo "##[group]Run function $name"

    riff function create $name --image $image --namespace $NAMESPACE --tail \
      --local-path $fats_dir/functions/uppercase/${test} &

    riff core deployer create $name \
      --function-ref $name \
      --ingress-policy External \
      --namespace $NAMESPACE \
      --tail
    source $fats_dir/macros/invoke_incluster.sh \
      "$(kubectl get deployers.core.projectriff.io ${name} --namespace ${NAMESPACE} -ojsonpath='{.status.address.url}')" \
      "${curl_opts}" \
      "${expected_data}"
    # TODO invoke via ingress as well
    source $fats_dir/macros/invoke_core_deployer.sh $name "-H Content-Type:text/plain -H Accept:text/plain -d cli" CLI
    riff core deployer delete $name --namespace $NAMESPACE

    riff knative deployer create $name \
      --function-ref $name \
      --ingress-policy External \
      --namespace $NAMESPACE \
      --tail
    source $fats_dir/macros/invoke_incluster.sh \
      "$(kubectl get deployers.knative.projectriff.io ${name} --namespace ${NAMESPACE} -ojsonpath='{.status.address.url}')" \
      "${curl_opts}" \
      "${expected_data}"
    source $fats_dir/macros/invoke_knative_deployer.sh \
      "${name}" \
      "${curl_opts}" \
      "${expected_data}"
    riff knative deployer delete $name --namespace $NAMESPACE

    riff function delete $name --namespace $NAMESPACE
    fats_delete_image $image

    echo "##[endgroup]"
  done
fi
