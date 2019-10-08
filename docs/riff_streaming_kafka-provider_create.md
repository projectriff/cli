---
id: riff-streaming-kafka-provider-create
title: "riff streaming kafka-provider create"
---
## riff streaming kafka-provider create

create a kafka provider of messages

### Synopsis

<todo>

```
riff streaming kafka-provider create <name> [flags]
```

### Examples

```
riff streaming kafka-provider create my-kafka-provider --bootstrap-servers kafka.local:9092
```

### Options

```
      --bootstrap-servers address   address of the kafka broker
      --dry-run                     print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr
  -h, --help                        help for create
  -n, --namespace name              kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff streaming kafka-provider](riff_streaming_kafka-provider.md)	 - (experimental) kafka stream provider

