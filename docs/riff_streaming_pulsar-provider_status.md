---
id: riff-streaming-pulsar-provider-status
title: "riff streaming pulsar-provider status"
---
## riff streaming pulsar-provider status

show pulsar provider status

### Synopsis

Display status details for a pulsar provider.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
pulsar provider roll out is being processed.

```
riff streaming pulsar-provider status <name> [flags]
```

### Examples

```
riff streamming pulsar-provider status my-pulsar-provider
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

* [riff streaming pulsar-provider](riff_streaming_pulsar-provider.md)	 - (experimental) pulsar stream provider

