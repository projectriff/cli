---
id: riff-core
title: "riff core"
---
## riff core

core runtime for riff functions and applications

### Synopsis

The core runtime uses stock kubernetes resources to deploy a function or
application. A Deployment is created along with a Service to forward traffic to
the deployment.

Ingress and autoscalers are not provided.

### Options

```
  -h, --help   help for core
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff](riff.md)	 - riff is for functions
* [riff core handler](riff_core_handler.md)	 - handlers map HTTP requests to applications, functions or images

