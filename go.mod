module github.com/projectriff/cli

go 1.13

require (
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/boz/go-logutil v0.1.0
	github.com/boz/kail v0.12.0
	github.com/buildpack/pack v0.5.0
	github.com/fatih/color v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-cmp v0.3.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/projectriff/system v0.0.0-20191209165036-1a0ef53e2e5c
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	// equivalent of kubernetes-1.16.3 tag for each k8s.io repo
	k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apiextensions-apiserver v0.0.0-20191114105449-027877536833
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
)
