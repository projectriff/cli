module github.com/projectriff/cli

go 1.13

require (
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/boz/go-logutil v0.1.0
	github.com/boz/kail v0.10.1
	github.com/buildpack/pack v0.4.1
	github.com/fatih/color v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-cmp v0.3.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/projectriff/system v0.0.0-20191003202856-6643ca21d1a3
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb
	k8s.io/apiextensions-apiserver v0.0.0-20190606210616-f848dc7be4a4
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

require (
	github.com/Azure/go-autorest/autorest v0.9.1 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/gophercloud/gophercloud v0.4.0 // indirect
)

replace (
	// force specific versions, go wants to revert normal dependencies
	// equivelent of kubernetes-1.14.3 tag for each repo
	k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190606210616-f848dc7be4a4
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)
