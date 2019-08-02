---
id: riff-knative-handler
title: "riff knative handler"
---
## riff knative handler

handlers map HTTP requests to applications, functions or images

### Synopsis

Handlers can be created for one of an application, function or image.
Application and function based handlers continuously watch for the latest built
image and will deploy new images. If the underlying application or function is
deleted, the handler will continue to run, but will no longer self update. Image
based handlers must be manually updated to trigger roll out of an updated image.

Applications, functions and images are logically equivalent at runtime.
Functions with an invoker are more focused and opinionated applications, and
images are compiled applications.

Users wishing to perform checks on built images before deploying them can
provide their own external process to watch the application/function for new
images and only update the handler image once those checks pass.

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

* [riff knative](riff_knative.md)	 - knative runtime for riff functions
* [riff knative handler create](riff_knative_handler_create.md)	 - create a handler to map HTTP requests to an application, function, container or image
* [riff knative handler delete](riff_knative_handler_delete.md)	 - delete handler(s)
* [riff knative handler list](riff_knative_handler_list.md)	 - table listing of handlers
* [riff knative handler status](riff_knative_handler_status.md)	 - show knative handler status
* [riff knative handler tail](riff_knative_handler_tail.md)	 - watch handler logs

