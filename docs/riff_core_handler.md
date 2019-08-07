---
id: riff-core-handler
title: "riff core handler"
---
## riff core handler

handlers deploy a workload

### Synopsis

Handlers can be created for a build or an image. Build based handlers
continuously watch for the latest image and will deploy new images. If the
underlying build is deleted, the handler will continue to run, but will no
longer self update. Image based handlers must be manually updated to trigger
roll out of an updated image.

Users wishing to perform checks on built images before deploying them can
provide their own external process to watch the build for new images and only
update the handler image once those checks pass.

The service to access the handler is available in the handler listing.

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

* [riff core](riff_core.md)	 - core runtime for riff workloads
* [riff core handler create](riff_core_handler_create.md)	 - create a handler to deploy a workload
* [riff core handler delete](riff_core_handler_delete.md)	 - delete handler(s)
* [riff core handler list](riff_core_handler_list.md)	 - table listing of handlers
* [riff core handler status](riff_core_handler_status.md)	 - show core handler status
* [riff core handler tail](riff_core_handler_tail.md)	 - watch handler logs

