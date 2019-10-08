---
id: riff-streaming-stream-create
title: "riff streaming stream create"
---
## riff streaming stream create

create a stream of messages

### Synopsis

<todo>

```
riff streaming stream create <name> [flags]
```

### Examples

```
riff streaming stream create --provider my-provider
```

### Options

```
      --content-type MIME type   MIME type for message payloads accepted by the stream
      --dry-run                  print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
  -h, --help                     help for create
  -n, --namespace name           kubernetes namespace (defaulted from kube config)
      --provider name            name of stream provider
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff streaming stream](riff_streaming_stream.md)	 - (experimental) streams of messages

