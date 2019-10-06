---
id: riff-streaming-kafka-provider-list
title: "riff streaming kafka-provider list"
---
## riff streaming kafka-provider list

table listing of kafka providers

### Synopsis

List kafka providers in a namespace or across all namespaces.

For detail regarding the status of a single kafka provider, run:

    riff streaming kafka-provider status <kafka-provider-name>

```
riff streaming kafka-provider list [flags]
```

### Examples

```
riff streaming kafka-provider list
riff streaming kafka-provider list --all-namespaces
```

### Options

```
      --all-namespaces   use all kubernetes namespaces
  -h, --help             help for list
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff streaming kafka-provider](riff_streaming_kafka-provider.md)	 - (experimental) kafka stream provider

