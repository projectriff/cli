---
id: riff-streaming-pulsar-provider-create
title: "riff streaming pulsar-provider create"
---
## riff streaming pulsar-provider create

create a pulsar provider of messages

### Synopsis

<todo>

```
riff streaming pulsar-provider create <name> [flags]
```

### Examples

```
riff streaming pulsar-provider create my-pulsar-provider --service-url pulsar://localhost:6650
```

### Options

```
      --dry-run                 print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
  -h, --help                    help for create
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --service-url url         url of the pulsar service
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

* [riff streaming pulsar-provider](riff_streaming_pulsar-provider.md)	 - (experimental) pulsar stream provider

