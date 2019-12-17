---
id: riff-streaming-pulsar-provider-list
title: "riff streaming pulsar-provider list"
---
## riff streaming pulsar-provider list

table listing of pulsar providers

### Synopsis

List pulsar providers in a namespace or across all namespaces.

For detail regarding the status of a single pulsar provider, run:

    riff streaming pulsar-provider status <pulsar-provider-name>

```
riff streaming pulsar-provider list [flags]
```

### Examples

```
riff streaming pulsar-provider list
riff streaming pulsar-provider list --all-namespaces
```

### Options

```
      --all-namespaces   use all kubernetes namespaces
  -h, --help             help for list
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff streaming pulsar-provider](riff_streaming_pulsar-provider.md)	 - (experimental) pulsar stream provider

