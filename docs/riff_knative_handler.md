---
id: riff-knative-handler
title: "riff knative handler"
---
## riff knative handler

handlers map HTTP requests to a workload

### Synopsis

Handlers can be created for a build reference or image. Build based handlers
continuously watch for the latest built image and will deploy new images. If the
underlying build resource is deleted, the handler will continue to run, but will
no longer self update. Image based handlers must be manually updated to trigger
roll out of an updated image.

Users wishing to perform checks on built images before deploying them can
provide their own external process to watch the build resource for new images
and only update the handler image once those checks pass.

The hostname to access the handler is available in the handler listing.

### Options

```
  -h, --help   help for handler
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff knative](riff_knative.md)	 - Knative runtime for riff workloads
* [riff knative handler create](riff_knative_handler_create.md)	 - create a handler to map HTTP requests to a workload
* [riff knative handler delete](riff_knative_handler_delete.md)	 - delete handler(s)
* [riff knative handler list](riff_knative_handler_list.md)	 - table listing of handlers
* [riff knative handler status](riff_knative_handler_status.md)	 - show knative handler status
* [riff knative handler tail](riff_knative_handler_tail.md)	 - watch handler logs

