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
	github.com/projectriff/system v0.0.0-20191106092351-69e4bee9b2e8
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	// equivalent of kubernetes-1.15.4 tag for each k8s.io repo
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apiextensions-apiserver v0.0.0-20190918201827-3de75813f604
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
)

replace github.com/projectriff/system => github.com/scothis/system v0.0.0-20191106210539-b6f64de620f8
