---
id: riff-core-handler-status
title: "riff core handler status"
---
## riff core handler status

show core handler status

### Synopsis

Display status details for a handler.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
handler roll out is processed.

```
riff core handler status <name> [flags]
```

### Examples

```
riff core handler status my-handler
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

* [riff core handler](riff_core_handler.md)	 - handlers deploy a workload

