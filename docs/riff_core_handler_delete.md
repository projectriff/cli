---
id: riff-core-handler-delete
title: "riff core handler delete"
---
## riff core handler delete

delete handler(s)

### Synopsis

Delete one or more handlers by name or all handlers within a namespace.

```
riff core handler delete <name(s)> [flags]
```

### Examples

```
riff core handler delete my-handler
riff core handler delete --all
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

* [riff core handler](riff_core_handler.md)	 - handlers deploy a workload

