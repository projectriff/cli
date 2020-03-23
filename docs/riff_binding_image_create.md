---
id: riff-binding-image-create
title: "riff binding image create"
---
## riff binding image create

create a image to deploy a workload

### Synopsis

Create an image binding.

<todo>

```
riff binding image create <name> [flags]
```

### Examples

```
riff binding image create my-image-binding
```

### Options

```
      --dry-run                     print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
  -h, --help                        help for create
  -n, --namespace name              kubernetes namespace (defaulted from kube config)
      --provider object reference   provider object reference to get images from
      --subject object reference    subject object reference to inject images into
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff binding image](riff_binding_image.md)	 - <todo>

