/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands_test

import (
	"context"
	"errors"
	"fmt"
	goruntime "runtime"
	"strings"
	"testing"

	"github.com/buildpacks/pack"
	"github.com/projectriff/cli/pkg/build/commands"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/k8s"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	kailtesting "github.com/projectriff/cli/pkg/testing/kail"
	packtesting "github.com/projectriff/cli/pkg/testing/pack"
	buildv1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	cachetesting "k8s.io/client-go/tools/cache/testing"
)

func TestApplicationCreateOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidResourceOptionsFieldError.Also(
				cli.ErrMissingField(cli.ImageFlagName),
				cli.ErrMissingOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName),
			),
		},
		{
			Name: "git source",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
			},
			ShouldValidate: true,
		},
		{
			Name: "local source",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				LocalPath:       ".",
			},
			ShouldValidate: true,
		},
		{
			Name: "no source",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
			},
			ExpectFieldErrors: cli.ErrMissingOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName),
		},
		{
			Name: "multiple sources",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				LocalPath:       ".",
			},
			ExpectFieldErrors: cli.ErrMultipleOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName),
		},
		{
			Name: "git source with cache",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				CacheSize:       "8Gi",
			},
			ShouldValidate: true,
		},
		{
			Name: "local source with cache",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				LocalPath:       ".",
				CacheSize:       "8Gi",
			},
			ExpectFieldErrors: cli.ErrDisallowedFields(cli.CacheSizeFlagName, ""),
		},
		{
			Name: "invalid cache",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				CacheSize:       "X",
			},
			ExpectFieldErrors: cli.ErrInvalidValue("X", cli.CacheSizeFlagName),
		},
		{
			Name: "with git subpath",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				SubPath:         "some/directory",
			},
			ShouldValidate: true,
		},
		{
			Name: "with local subpath",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				LocalPath:       ".",
				SubPath:         "some/directory",
			},
			ExpectFieldErrors: cli.ErrDisallowedFields(cli.SubPathFlagName, ""),
		},
		{
			Name: "missing git revision",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "",
			},
			ExpectFieldErrors: cli.ErrMissingField(cli.GitRevisionFlagName),
		},
		{
			Name: "with env",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				Env:             []string{"VAR1=foo", "VAR2=bar"},
			},
			ShouldValidate: true,
		},
		{
			Name: "with invalid env",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				Env:             []string{"=foo"},
			},
			ExpectFieldErrors: cli.ErrInvalidArrayValue("=foo", cli.EnvFlagName, 0),
		},
		{
			Name: "with limits",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				LimitCPU:        "500m",
				LimitMemory:     "512Mi",
			},
			ShouldValidate: true,
		},
		{
			Name: "with invalid limits",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				LimitCPU:        "50%",
				LimitMemory:     "NaN",
			},
			ExpectFieldErrors: cli.FieldErrors{}.Also(
				cli.ErrInvalidValue("50%", cli.LimitCPUFlagName),
				cli.ErrInvalidValue("NaN", cli.LimitMemoryFlagName),
			),
		},
		{
			Name: "git source, tail",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				Tail:            true,
				WaitTimeout:     "10m",
			},
			ShouldValidate: true,
		},
		{
			Name: "git source, tail missing timeout",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				Tail:            true,
			},
			ExpectFieldErrors: cli.ErrMissingField(cli.WaitTimeoutFlagName),
		},
		{
			Name: "git source, tail invalid timeout",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				Tail:            true,
				WaitTimeout:     "d",
			},
			ExpectFieldErrors: cli.ErrInvalidValue("d", cli.WaitTimeoutFlagName),
		},
		{
			Name: "dry run",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				DryRun:          true,
			},
			ShouldValidate: true,
		},
		{
			Name: "dry run, tail",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				Tail:            true,
				WaitTimeout:     "10m",
				DryRun:          true,
			},
			ExpectFieldErrors: cli.ErrMultipleOneOf(cli.DryRunFlagName, cli.TailFlagName),
		},
	}

	if goruntime.GOOS == "windows" {
		for i, tr := range table {
			opts, _ := tr.Options.(*commands.ApplicationCreateOptions)
			if opts.LocalPath != "" {
				tr.ShouldValidate = false
				tr.ExpectFieldErrors = tr.ExpectFieldErrors.Also(
					cli.ErrInvalidValue(fmt.Sprintf("%s is not available on Windows", cli.LocalPathFlagName), cli.LocalPathFlagName),
				)
			}
			table[i] = tr
		}
	}

	table.Run(t)
}

func TestApplicationCreateCommand(t *testing.T) {
	defaultNamespace := "default"
	applicationName := "my-application"
	imageDefault := "_"
	imageTag := "registry.example.com/repo:tag"
	registryHost := "registry.example.com"
	gitRepo := "https://example.com/repo.git"
	gitMaster := "master"
	gitSha := "deadbeefdeadbeefdeadbeefdeadbeef"
	subPath := "some/path"
	cacheSize := "8Gi"
	cacheSizeQuantity := resource.MustParse(cacheSize)
	localPath := "."

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "git repo",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo, dry run",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.DryRunFlagName},
			ExpectOutput: `
---
apiVersion: build.projectriff.io/v1alpha1
kind: Application
metadata:
  creationTimestamp: null
  name: my-application
  namespace: default
spec:
  build:
    resources: {}
  image: registry.example.com/repo:tag
  source:
    git:
      revision: master
      url: https://example.com/repo.git
status: {}

Created application "my-application"
`,
		},
		{
			Name: "git repo with revision",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.GitRevisionFlagName, gitSha},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitSha,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo with subpath",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.SubPathFlagName, subPath},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
							SubPath: subPath,
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo with cache",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.CacheSizeFlagName, cacheSize},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image:     imageTag,
						CacheSize: &cacheSizeQuantity,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo with env",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.EnvFlagName, "MY_VAR1=value1", cli.EnvFlagName, "MY_VAR2=value2"},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Build: buildv1alpha1.ImageBuild{
							Env: []corev1.EnvVar{
								{Name: "MY_VAR1", Value: "value1"},
								{Name: "MY_VAR2", Value: "value2"},
							},
						},
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo with limits",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.LimitCPUFlagName, "100m", cli.LimitMemoryFlagName, "128Mi"},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Build: buildv1alpha1.ImageBuild{
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("128Mi"),
								},
							},
						},
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "local path",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.LocalPathFlagName, localPath},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				packClient := &packtesting.Client{}
				c.Pack = packClient
				packClient.On("Build", mock.Anything, pack.BuildOptions{
					Image:   imageTag,
					AppPath: localPath,
					Builder: "cloudfoundry/cnb:bionic",
					Env:     map[string]string{},
					Publish: true,
				}).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...build output...\n")
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				packClient := c.Pack.(*packtesting.Client)
				packClient.AssertExpectations(t)
				return nil
			},
			GivenObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "riff-system",
						Name:      "builders",
					},
					Data: map[string]string{
						"riff-application": "cloudfoundry/cnb:bionic",
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
					},
				},
			},
			ExpectOutput: `
...build output...
Created application "my-application"
`,
		},
		{
			Name: "local path with env",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.LocalPathFlagName, localPath, cli.EnvFlagName, "MY_VAR1=value1", cli.EnvFlagName, "MY_VAR2=value2"},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				packClient := &packtesting.Client{}
				c.Pack = packClient
				packClient.On("Build", mock.Anything, pack.BuildOptions{
					Image:   imageTag,
					AppPath: localPath,
					Builder: "cloudfoundry/cnb:bionic",
					Env: map[string]string{
						"MY_VAR1": "value1",
						"MY_VAR2": "value2",
					},
					Publish: true,
				}).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...build output...\n")
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				packClient := c.Pack.(*packtesting.Client)
				packClient.AssertExpectations(t)
				return nil
			},
			GivenObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "riff-system",
						Name:      "builders",
					},
					Data: map[string]string{
						"riff-application": "cloudfoundry/cnb:bionic",
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Build: buildv1alpha1.ImageBuild{
							Env: []corev1.EnvVar{
								{Name: "MY_VAR1", Value: "value1"},
								{Name: "MY_VAR2", Value: "value2"},
							},
						},
						Image: imageTag,
					},
				},
			},
			ExpectOutput: `
...build output...
Created application "my-application"
`,
		},
		{
			Name: "local path, dry run",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.LocalPathFlagName, localPath, cli.DryRunFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				packClient := &packtesting.Client{}
				c.Pack = packClient
				packClient.On("Build", mock.Anything, pack.BuildOptions{
					Image:   imageTag,
					AppPath: localPath,
					Builder: "cloudfoundry/cnb:bionic",
					Env:     map[string]string{},
					Publish: true,
				}).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...build output...\n")
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				packClient := c.Pack.(*packtesting.Client)
				packClient.AssertExpectations(t)
				return nil
			},
			GivenObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "riff-system",
						Name:      "builders",
					},
					Data: map[string]string{
						"riff-application": "cloudfoundry/cnb:bionic",
					},
				},
			},
			ExpectOutput: `
...build output...
---
apiVersion: build.projectriff.io/v1alpha1
kind: Application
metadata:
  creationTimestamp: null
  name: my-application
  namespace: default
spec:
  build:
    resources: {}
  image: registry.example.com/repo:tag
status: {}

Created application "my-application"
`,
		},
		{
			Name: "local path, no builders",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.LocalPathFlagName, localPath},
			ExpectOutput: `
`,
			ShouldError: true,
			Verify: func(t *testing.T, output string, err error) {
				if expected, actual := `configmaps "builders" not found`, err.Error(); expected != actual {
					t.Errorf("expected error %q, actual %q", expected, actual)
				}
			},
		},
		{
			Name: "local path, no application builder",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.LocalPathFlagName, localPath},
			GivenObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "riff-system",
						Name:      "builders",
					},
					Data: map[string]string{},
				},
			},
			ExpectOutput: `
`,
			ShouldError: true,
			Verify: func(t *testing.T, output string, err error) {
				if expected, actual := `unknown builder for "riff-application"`, err.Error(); expected != actual {
					t.Errorf("expected error %q, actual %q", expected, actual)
				}
			},
		},
		{
			Name: "local path, pack error",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.LocalPathFlagName, localPath},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				packClient := &packtesting.Client{}
				c.Pack = packClient
				packClient.On("Build", mock.Anything, pack.BuildOptions{
					Image:   imageTag,
					AppPath: localPath,
					Builder: "cloudfoundry/cnb:bionic",
					Env:     map[string]string{},
					Publish: true,
				}).Return(fmt.Errorf("pack error")).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...build output...\n")
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				packClient := c.Pack.(*packtesting.Client)
				packClient.AssertExpectations(t)
				return nil
			},
			GivenObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "riff-system",
						Name:      "builders",
					},
					Data: map[string]string{
						"riff-application": "cloudfoundry/cnb:bionic",
					},
				},
			},
			ExpectOutput: `
...build output...
`,
			ShouldError: true,
			Verify: func(t *testing.T, output string, err error) {
				if expected, actual := "pack error", err.Error(); expected != actual {
					t.Errorf("expected error %q, actual %q", expected, actual)
				}
			},
		},
		{
			Name: "local path, default image",
			Args: []string{applicationName, cli.LocalPathFlagName, localPath},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				packClient := &packtesting.Client{}
				c.Pack = packClient
				packClient.On("Build", mock.Anything, pack.BuildOptions{
					Image:   fmt.Sprintf("%s/%s", registryHost, applicationName),
					AppPath: localPath,
					Builder: "cloudfoundry/cnb:bionic",
					Env:     map[string]string{},
					Publish: true,
				}).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...build output...\n")
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				packClient := c.Pack.(*packtesting.Client)
				packClient.AssertExpectations(t)
				return nil
			},
			GivenObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "riff-build",
					},
					Data: map[string]string{
						"default-image-prefix": registryHost,
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "riff-system",
						Name:      "builders",
					},
					Data: map[string]string{
						"riff-application": "cloudfoundry/cnb:bionic",
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageDefault,
					},
				},
			},
			ExpectOutput: `
...build output...
Created application "my-application"
`,
		},
		{
			Name:        "local path, default image, missing default",
			Args:        []string{applicationName, cli.LocalPathFlagName, localPath},
			ShouldError: true,
			Verify: func(t *testing.T, output string, err error) {
				if expected, actual := "default image prefix requires initialized credentials, run `riff help credentials`", err.Error(); expected != actual {
					t.Errorf("expected error %q, actual %q", expected, actual)
				}
			},
		},
		{
			Name: "local path, default image, configmap load error",
			Args: []string{applicationName, cli.LocalPathFlagName, localPath},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("get", "configmaps"),
			},
			ShouldError: true,
		},
		{
			Name: "local path, default image",
			Args: []string{applicationName, cli.LocalPathFlagName, localPath},
			GivenObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "riff-build",
					},
					Data: map[string]string{
						"default-image-prefix": "",
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error existing application",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo},
			GivenObjects: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error during create",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("create", "applications"),
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "tail logs",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("ApplicationLogs", mock.Anything, &buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				}, cli.TailSinceCreateDefault, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...log output...\n")
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				if lw, ok := k8s.GetListerWatcher(ctx, nil, "", nil).(*cachetesting.FakeControllerSource); ok {
					lw.Shutdown()
				}

				kail := c.Kail.(*kailtesting.Logger)
				kail.AssertExpectations(t)
				return nil
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
Waiting for application "my-application" to become ready...
...log output...
Application "my-application" is ready
`,
		},
		{
			Name: "tail timeout",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.TailFlagName, cli.WaitTimeoutFlagName, "5ms"},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("ApplicationLogs", mock.Anything, &buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				}, cli.TailSinceCreateDefault, mock.Anything).Return(k8s.ErrWaitTimeout).Run(func(args mock.Arguments) {
					ctx := args[0].(context.Context)
					fmt.Fprintf(c.Stdout, "...log output...\n")
					// wait for context to be cancelled
					<-ctx.Done()
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				if lw, ok := k8s.GetListerWatcher(ctx, nil, "", nil).(*cachetesting.FakeControllerSource); ok {
					lw.Shutdown()
				}

				kail := c.Kail.(*kailtesting.Logger)
				kail.AssertExpectations(t)
				return nil
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
Waiting for application "my-application" to become ready...
...log output...
Timeout after "5ms" waiting for "my-application" to become ready
To view status run: riff application list --namespace default
To continue watching logs run: riff application tail my-application --namespace default
`,
			ShouldError: true,
			Verify: func(t *testing.T, output string, err error) {
				if actual := err; !errors.Is(err, cli.SilentError) {
					t.Errorf("expected error to be silent, actual %#v", actual)
				}
			},
		},
		{
			Name: "tail error",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("ApplicationLogs", mock.Anything, &buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				}, cli.TailSinceCreateDefault, mock.Anything).Return(fmt.Errorf("kail error"))
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				if lw, ok := k8s.GetListerWatcher(ctx, nil, "", nil).(*cachetesting.FakeControllerSource); ok {
					lw.Shutdown()
				}

				kail := c.Kail.(*kailtesting.Logger)
				kail.AssertExpectations(t)
				return nil
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: &buildv1alpha1.Source{
							Git: &buildv1alpha1.Git{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
Waiting for application "my-application" to become ready...
`,
			ShouldError: true,
		},
	}

	if goruntime.GOOS == "windows" {
		for i, tr := range table {
			if strings.Contains(strings.Join(tr.Args, "|"), cli.LocalPathFlagName) {
				tr.Skip = true
				table[i] = tr
			}
		}
	}

	table.Run(t, commands.NewApplicationCreateCommand)
}
