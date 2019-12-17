---
id: riff-streaming-inmemory-provider-status
title: "riff streaming inmemory-provider status"
---
## riff streaming inmemory-provider status

show inmemory provider status

### Synopsis

Display status details for a in-memory provider.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
in-memory provider roll out is being processed.

```
riff streaming inmemory-provider status <name> [flags]
```

### Examples

```
riff streamming inmemory-provider status my-inmemory-provider
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

* [riff streaming inmemory-provider](riff_streaming_inmemory-provider.md)	 - (experimental) in-memory stream provider

