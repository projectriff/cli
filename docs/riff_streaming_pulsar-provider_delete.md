---
id: riff-streaming-pulsar-provider-delete
title: "riff streaming pulsar-provider delete"
---
## riff streaming pulsar-provider delete

delete pulsar provider(s)

### Synopsis

Delete one or more pulsar providers by name or all pulsar providers within a
namespace.

Deleting a pulsar provider will disrupt all processors consuming streams managed
by the provider. Existing messages in the stream may be preserved by the
underlying pulsar broker, depending on the implementation.

```
riff streaming pulsar-provider delete <name(s)> [flags]
```

### Examples

```
riff streaming pulsar-provider delete my-pulsar-provider
riff streaming pulsar-provider delete --all 
```

### Options

```
      --all              delete all pulsar providers within the namespace
  -h, --help             help for delete
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

