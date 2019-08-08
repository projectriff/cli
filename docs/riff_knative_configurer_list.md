---
id: riff-knative-configurer-list
title: "riff knative configurer list"
---
## riff knative configurer list

table listing of configurers

### Synopsis

List configurers in a namespace or across all namespaces.

For detail regarding the status of a single configurer, run:

    riff knative configurer status <configurer-name>

```
riff knative configurer list [flags]
```

### Examples

```
riff knative configurer list
riff knative configurer list --all-namespaces
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

* [riff knative configurer](riff_knative_configurer.md)	 - configurers map HTTP requests to a workload

