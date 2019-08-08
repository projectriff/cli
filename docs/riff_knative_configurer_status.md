---
id: riff-knative-configurer-status
title: "riff knative configurer status"
---
## riff knative configurer status

show knative configurer status

### Synopsis

Display status details for a configurer.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
configurer roll out is processed.

```
riff knative configurer status <name> [flags]
```

### Examples

```
riff knative configurer status my-configurer
```

### Options

```
  -h, --help             help for status
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

