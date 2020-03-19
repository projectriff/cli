---
id: riff-binding-image-delete
title: "riff binding image delete"
---
## riff binding image delete

delete image(s)

### Synopsis

Delete one or more images by name or all images within a namespace.

```
riff binding image delete <name(s)> [flags]
```

### Examples

```
riff binding image delete my-image-binding
riff binding image delete --all
```

### Options

```
      --all              delete all images within the namespace
  -h, --help             help for delete
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

