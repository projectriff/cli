---
id: riff-streaming-inmemory-provider-delete
title: "riff streaming inmemory-provider delete"
---
## riff streaming inmemory-provider delete

delete in-memory provider(s)

### Synopsis

Delete one or more in-memory providers by name or all in-memory providers within
a namespace.

Deleting a in-memory provider will disrupt all processors consuming streams
managed by the provider. Existing messages in the stream may be preserved by the
underlying in-memory broker, depending on the implementation.

```
riff streaming inmemory-provider delete <name(s)> [flags]
```

### Examples

```
riff streaming inmemory-provider delete my-inmemory-provider
riff streaming inmemory-provider delete --all 
```

### Options

```
      --all              delete all inmemory providers within the namespace
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

* [riff streaming inmemory-provider](riff_streaming_inmemory-provider.md)	 - (experimental) in-memory stream provider

