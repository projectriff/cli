---
id: riff-container
title: "riff container"
---
## riff container

containers built from source using container buildpacks

### Synopsis

Containers are a mechanism to convert web container source code into
container images that can be invoked over HTTP. Cloud Native Buildpacks are
provided to detect the language, provide a language runtime, install build and
runtime dependencies, compile the container, and packaging everything as a
container.

The container resource is only responsible for converting source code into a
container. The container container image may then be deployed as a request
handler. See `riff handler --help` for detail.

### Options

```
  -h, --help   help for container
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff](riff.md)	 - riff is for functions
* [riff container create](riff_container_create.md)	 - create an container from source
* [riff container delete](riff_container_delete.md)	 - delete container(s)
* [riff container list](riff_container_list.md)	 - table listing of containers
* [riff container status](riff_container_status.md)	 - show container status

