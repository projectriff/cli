---
id: riff-streaming-inmemory-provider-create
title: "riff streaming inmemory-provider create"
---
## riff streaming inmemory-provider create

create a in-memory provider of messages

### Synopsis

<todo>

```
riff streaming inmemory-provider create <name> [flags]
```

### Examples

```
riff streaming inmemory-provider create my-inmemory-provider
```

### Options

```
      --dry-run                 print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
  -h, --help                    help for create
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --tail                    watch creation progress
      --wait-timeout duration   duration to wait for the provider to become ready when watching progress (default 1m0s)
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff streaming inmemory-provider](riff_streaming_inmemory-provider.md)	 - (experimental) in-memory stream provider

