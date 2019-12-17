---
id: riff-streaming-inmemory-provider-list
title: "riff streaming inmemory-provider list"
---
## riff streaming inmemory-provider list

table listing of in-memory providers

### Synopsis

List in-memory providers in a namespace or across all namespaces.

For detail regarding the status of a single in-memory provider, run:

    riff streaming inmemory-provider status <inmemory-provider-name>

```
riff streaming inmemory-provider list [flags]
```

### Examples

```
riff streaming inmemory-provider list
riff streaming inmemory-provider list --all-namespaces
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

* [riff streaming inmemory-provider](riff_streaming_inmemory-provider.md)	 - (experimental) in-memory stream provider

