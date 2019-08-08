---
id: riff-knative-configurer
title: "riff knative configurer"
---
## riff knative configurer

configurers map HTTP requests to a workload

### Synopsis

Configurers can be created for a build reference or image. Build based configurers
continuously watch for the latest built image and will deploy new images. If the
underlying build resource is deleted, the configurer will continue to run, but will
no longer self update. Image based configurers must be manually updated to trigger
roll out of an updated image.

Users wishing to perform checks on built images before deploying them can
provide their own external process to watch the build resource for new images
and only update the configurer image once those checks pass.

The hostname to access the configurer is available in the configurer listing.

### Options

```
  -h, --help   help for configurer
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff knative](riff_knative.md)	 - Knative runtime for riff workloads
* [riff knative configurer create](riff_knative_configurer_create.md)	 - create a configurer to map HTTP requests to a workload
* [riff knative configurer delete](riff_knative_configurer_delete.md)	 - delete configurer(s)
* [riff knative configurer list](riff_knative_configurer_list.md)	 - table listing of configurers
* [riff knative configurer status](riff_knative_configurer_status.md)	 - show knative configurer status
* [riff knative configurer tail](riff_knative_configurer_tail.md)	 - watch configurer logs

