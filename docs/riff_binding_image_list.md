---
id: riff-binding-image-list
title: "riff binding image list"
---
## riff binding image list

table listing of image bindings

### Synopsis

List image bindings in a namespace or across all namespaces.

For detail regarding the status of a single image, run:

    riff binding image status <image-binding-name>

```
riff binding image list [flags]
```

### Examples

```
riff binding image list
riff binding image list --all-namespaces
```

### Options

```
      --all-namespaces   use all kubernetes namespaces
  -h, --help             help for list
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

