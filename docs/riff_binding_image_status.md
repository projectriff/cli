---
id: riff-binding-image-status
title: "riff binding image status"
---
## riff binding image status

show image binding status

### Synopsis

Display status details for an image binding.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
image roll out is processed.

```
riff binding image status <name> [flags]
```

### Examples

```
riff binding image status my-imagebinding
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

* [riff binding image](riff_binding_image.md)	 - <todo>

