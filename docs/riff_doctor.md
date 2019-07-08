## riff doctor

check riff's requirements are installed

### Synopsis

Check that riff is install and configured for usage.

The doctor checks:
- necessary system components are installed

The checkup will include more checks in the future as we discover common issues.
The doctor is not a tool for monitoring the health of a cluster or the install.
Usage is contextualized to a specific user and namespace. An issue with one user
or namespace may not indicate systemic issues.

```
riff doctor [flags]
```

### Examples

```
riff doctor
```

### Options

```
  -h, --help   help for doctor
```

### Options inherited from parent commands

```
      --config file        config file (default is $HOME/.riff.yaml)
      --kube-config file   kubectl config file (default is $HOME/.kube/config)
      --no-color           disable color output in terminals
```

### SEE ALSO

* [riff](riff.md)	 - riff is for functions

