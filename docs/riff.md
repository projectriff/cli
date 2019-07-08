## riff

riff is for functions

### Synopsis

The riff CLI combines with the projectriff system CRDs to build, run and wire
workloads (applications and functions). This CLI is one interface to the riff
API.

Before running riff, please install the projectriff system and its dependencies.
See https://projectriff.io/docs/getting-started/

This CLI contains commands that fit into four themes:
- builds
- requests
- streams
- misc

Build commands include the application and function command to define build
plans and the credential commands to authenticate builds to container
registries.

Requests commands include the handler command to map HTTP requests to a built
application, function or container image.

Streams commands include the stream and processor commands to define streams of
messages and map those streams to function inputs and outputs with processors.

Miscellaneous commands include shell completion and the doctor to diagnose
common issues.

### Options

```
      --config file        config file (default is $HOME/.riff.yaml)
  -h, --help               help for riff
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff application](riff_application.md)	 - applications built from source using application buildpacks
* [riff completion](riff_completion.md)	 - generate shell completion script
* [riff credential](riff_credential.md)	 - credentials for container registries
* [riff doctor](riff_doctor.md)	 - check riff's requirements are installed
* [riff function](riff_function.md)	 - functions built from source using function buildpacks
* [riff handler](riff_handler.md)	 - handlers map HTTP requests to applications, functions or images
* [riff processor](riff_processor.md)	 - processors apply functions to messages on streams
* [riff stream](riff_stream.md)	 - streams of messages

