---
id: riff-streaming-kafka-provider-status
title: "riff streaming kafka-provider status"
---
## riff streaming kafka-provider status

show kafka provider status

### Synopsis

Display status details for a kafka provider.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
kafka provider roll out is being processed.

```
riff streaming kafka-provider status <name> [flags]
```

### Examples

```
riff streamming kafka-provider status my-kafka-provider
```

### Options

```
  -h, --help             help for status
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff streaming kafka-provider](riff_streaming_kafka-provider.md)	 - (experimental) kafka stream provider

