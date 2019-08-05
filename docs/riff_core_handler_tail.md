---
id: riff-core-handler-tail
title: "riff core handler tail"
---
## riff core handler tail

watch handler logs

### Synopsis

Stream runtime logs for a handler until canceled. To cancel, press Ctl-c in the
shell or kill the process.

As new handler pods are started, the logs are displayed. To show historical logs
use --since.

```
riff core handler tail <name> [flags]
```

### Examples

```
riff core handler tail my-handler
riff core handler tail my-handler --since 1h
```

### Options

```
  -h, --help             help for tail
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
      --since duration   time duration to start reading logs from
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff core handler](riff_core_handler.md)	 - handlers deploy a workload

