---
id: riff-handler-create
title: "riff handler create"
---
## riff handler create

create a handler to map HTTP requests to an application, function or image

### Synopsis

Create an HTTP request handler.

There are three way to create a handler:
- from an application reference
- from a function reference
- from an image

Application and function references are resolved within the same namespace as
the handler. As the build produces new images, the image will roll out
automatically.

Image based handlers must be updated manually to roll out new images.

The runtime environment can be configured by --env for static key-value pairs
and --env-from to map values from a ConfigMap or Secret.

```
riff handler create <name> [flags]
```

### Examples

```
riff handler create my-app-handler --application-ref my-app
riff handler create my-func-handler --function-ref my-func
riff handler create my-image-handler --image registry.example.com/my-image:latest
```

### Options

```
      --application-ref name    name of application to deploy
      --dry-run                 print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
      --env variable            environment variable defined as a key value pair separated by an equals sign, example "--env MY_VAR=my-value" (may be set multiple times)
      --env-from variable       environment variable from a config map or secret, example "--env-from MY_SECRET_VALUE=secretKeyRef:my-secret-name:key-in-secret", "--env-from MY_CONFIG_MAP_VALUE=configMapKeyRef:my-config-map-name:key-in-config-map" (may be set multiple times)
      --function-ref name       name of function to deploy
  -h, --help                    help for create
      --image image             container image to deploy
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --tail                    watch handler logs
      --wait-timeout duration   duration to wait for the handler to become ready when watching logs (default "10m")
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff handler](riff_handler.md)	 - handlers map HTTP requests to applications, functions or images

