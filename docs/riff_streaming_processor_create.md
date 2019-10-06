---
id: riff-streaming-processor-create
title: "riff streaming processor create"
---
## riff streaming processor create

create a processor to apply a function to messages on streams

### Synopsis

<todo>

```
riff streaming processor create <name> [flags]
```

### Examples

```
riff streaming processor create my-processor --function-ref my-func --input my-input-stream
riff streaming processor create my-processor --function-ref my-func --input my-input-stream --input my-join-stream --output my-output-stream
```

### Options

```
      --dry-run                 print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
      --function-ref name       name of function build to deploy
  -h, --help                    help for create
      --input name              name of stream to read messages from (may be set multiple times)
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --output name             name of stream to write messages to (may be set multiple times)
      --tail                    watch processor logs
      --wait-timeout duration   duration to wait for the processor to become ready when watching logs (default "10m")
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff streaming processor](riff_streaming_processor.md)	 - (experimental) processors apply functions to messages on streams

