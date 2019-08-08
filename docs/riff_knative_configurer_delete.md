---
id: riff-knative-configurer-delete
title: "riff knative configurer delete"
---
## riff knative configurer delete

delete configurer(s)

### Synopsis

Delete one or more configurers by name or all configurers within a namespace.

New HTTP requests addressed to the configurer will fail. A new configurer created with
the same name will start to receive new HTTP requests addressed to the same
configurer.

```
riff knative configurer delete <name(s)> [flags]
```

### Examples

```
riff knative configurer delete my-configurer
riff knative configurer delete --all
```

### Options

```
      --all              delete all configurers within the namespace
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

* [riff knative configurer](riff_knative_configurer.md)	 - configurers map HTTP requests to a workload

