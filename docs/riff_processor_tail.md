---
id: riff-processor-tail
title: "riff processor tail"
---
## riff processor tail

watch processor logs

### Synopsis

<todo>

```
riff processor tail [flags]
```

### Examples

```
riff processor tail my-processor
riff processor tail my-processor --since 1h
```

### Options

```
  -h, --help             help for tail
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff processor](riff_processor.md)	 - processors apply functions to messages on streams

