#!/bin/bash

kubectl get deployments,services,pods --all-namespaces || true
echo ""
echo "RIFF:"
echo ""
kubectl get riff --all-namespaces || true
echo ""
echo "KPACK:"
echo ""
kubectl get clusterbuilders.build.pivotal.io,builders.build.pivotal.io,images.build.pivotal.io,sourceresolvers.build.pivotal.io,builds.build.pivotal.io --all-namespaces || true
echo ""
echo "KNATIVE:"
echo ""
kubectl get knative --all-namespaces || true
echo ""
echo "FAILING PODS:"
echo ""
kubectl get pods --all-namespaces --field-selector=status.phase!=Running \
| tail -n +2 | awk '{print "-n", $1, $2}' | xargs -L 1 kubectl describe pod || true
echo ""
echo "NODE:"
echo ""
kubectl describe node || true
echo ""
echo "RIFF:"
echo ""
kubectl describe riff --all-namespaces || true
echo ""
echo "KNATIVE:"
echo ""
kubectl describe knative --all-namespaces || true