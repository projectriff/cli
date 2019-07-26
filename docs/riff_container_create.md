---
id: riff-container-create
title: "riff container create"
---
## riff container create

create an container from source

### Synopsis

Create an container from source using the container Cloud Native Buildpack
builder.

Container source can be specified either as a Git repository or as a local
directory. Builds from Git are run in the cluster while builds from a local
directory are run inside a local Docker daemon and are orchestrated by this
command (in the future, builds from local source may also be run in the
cluster).

```
riff container create <name> [flags]
```

### Examples

```
riff container create my-app --image registry.example.com/image
```

### Options

```
      --dry-run                 print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
  -h, --help                    help for create
      --image repository        repository where the built images are pushed (default "_")
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --tail                    watch build logs
      --wait-timeout duration   duration to wait for the container to become ready when watching logs (default "10m")
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff container](riff_container.md)	 - containers built from source using container buildpacks

