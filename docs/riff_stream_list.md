---
id: riff-stream-list
title: "riff stream list"
---
## riff stream list

table listing of streams

### Synopsis

List streams in a namespace or across all namespaces.

For detail regarding the status of a single stream, run:

	riff stream status <stream-name>

```
riff stream list [flags]
```

### Examples

```
riff stream list
riff stream list --all-namespaces
```

### Options

```
      --all-namespaces   use all kubernetes namespaces
  -h, --help             help for list
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff stream](riff_stream.md)	 - streams of messages

