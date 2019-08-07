---
id: riff-knative-handler-status
title: "riff knative handler status"
---
## riff knative handler status

show knative handler status

### Synopsis

Display status details for a handler.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
handler roll out is processed.

```
riff knative handler status <name> [flags]
```

### Examples

```
riff knative handler status my-handler
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

* [riff knative handler](riff_knative_handler.md)	 - handlers map HTTP requests to a workload

