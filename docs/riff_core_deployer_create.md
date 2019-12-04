---
id: riff-core-deployer-create
title: "riff core deployer create"
---
## riff core deployer create

create a deployer to deploy a workload

### Synopsis

Create a core deployer.

Build references are resolved within the same namespace as the deployer. As the
build produces new images, the image will roll out automatically. Image based
deployers must be updated manually to roll out new images. Rolling out an image
means creating a deployment with a pod referencing the image and a service
referencing the deployment.

The runtime environment can be configured by --env for static key-value pairs
and --env-from to map values from a ConfigMap or Secret.

```
riff core deployer create <name> [flags]
```

### Examples

```
riff core deployer create my-app-deployer --application-ref my-app
riff core deployer create my-func-deployer --function-ref my-func
riff core deployer create my-func-deployer --container-ref my-container
riff core deployer create my-image-deployer --image registry.example.com/my-image:latest
```

### Options

```
      --application-ref name    name of application to deploy
      --container-ref name      name of container to deploy
      --dry-run                 print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
      --env variable            environment variable defined as a key value pair separated by an equals sign, example "--env MY_VAR=my-value" (may be set multiple times)
      --env-from variable       environment variable from a config map or secret, example "--env-from MY_SECRET_VALUE=secretKeyRef:my-secret-name:key-in-secret", "--env-from MY_CONFIG_MAP_VALUE=configMapKeyRef:my-config-map-name:key-in-config-map" (may be set multiple times)
      --function-ref name       name of function to deploy
  -h, --help                    help for create
      --image image             container image to deploy
      --ingress-policy policy   ingress policy for network access to the workload, one of "ClusterLocal" or "External" (default "ClusterLocal")
      --limit-cpu cores         the maximum amount of cpu allowed, in CPU cores (500m = .5 cores)
      --limit-memory bytes      the maximum amount of memory allowed, in bytes (500Mi = 500MiB = 500 * 1024 * 1024)
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --tail                    watch deployer logs
      --target-port port        port that the workload listens on for traffic. The value is exposed to the workload as the PORT environment variable
      --wait-timeout duration   duration to wait for the deployer to become ready when watching logs (default "10m")
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff core deployer](riff_core_deployer.md)	 - deployers deploy a workload

