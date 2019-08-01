---
id: riff-knative-handler-delete
title: "riff knative handler delete"
---
## riff knative handler delete

delete handler(s)

### Synopsis

Delete one or more handlers by name or all handlers within a namespace.

New HTTP requests addressed to the handler will fail. A new handler created with
the same name will start to receive new HTTP requests addressed to the same
handler.

```
riff knative handler delete <name(s)> [flags]
```

### Examples

```
riff handler delete my-handler
riff handler delete --all 
```

### Options

```
      --all              delete all handlers within the namespace
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

* [riff knative handler](riff_knative_handler.md)	 - handlers map HTTP requests to applications, functions or images

