---
id: riff-streaming-kafka-provider-delete
title: "riff streaming kafka-provider delete"
---
## riff streaming kafka-provider delete

delete kafka provider(s)

### Synopsis

Delete one or more kafka providers by name or all kafka providers within a
namespace.

Deleting a kafka provider will disrupt all processors consuming streams managed
by the provider. Existing messages in the stream may be preserved by the
underlying kafka broker, depending on the implementation.

```
riff streaming kafka-provider delete <name(s)> [flags]
```

### Examples

```
riff streaming kafka-provider delete my-kafka-provider
riff streaming kafka-provider delete --all 
```

### Options

```
      --all              delete all kafka providers within the namespace
  -h, --help             help for delete
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

