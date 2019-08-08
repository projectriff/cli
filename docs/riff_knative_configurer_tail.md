---
id: riff-knative-configurer-tail
title: "riff knative configurer tail"
---
## riff knative configurer tail

watch configurer logs

### Synopsis

Stream runtime logs for a configurer until canceled. To cancel, press Ctl-c in the
shell or kill the process.

As new configurer pods are started, the logs are displayed. To show historical logs
use --since.

```
riff knative configurer tail <name> [flags]
```

### Examples

```
riff knative configurer tail my-configurer
riff knative configurer tail my-configurer --since 1h
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

* [riff knative configurer](riff_knative_configurer.md)	 - configurers map HTTP requests to a workload

