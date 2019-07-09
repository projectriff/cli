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

package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/buildpack/pack"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/k8s"
	"github.com/projectriff/cli/pkg/race"
	buildv1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FunctionCreateOptions struct {
	cli.ResourceOptions

	Image     string
	CacheSize string

	Artifact string
	Handler  string
	Invoker  string

	LocalPath   string
	GitRepo     string
	GitRevision string
	SubPath     string

	Tail        bool
	WaitTimeout string

	DryRun bool
}

var (
	_ cli.Validatable = (*FunctionCreateOptions)(nil)
	_ cli.Executable  = (*FunctionCreateOptions)(nil)
	_ cli.DryRunable  = (*FunctionCreateOptions)(nil)
)

func (opts *FunctionCreateOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	if opts.Image == "" {
		errs = errs.Also(cli.ErrMissingField(cli.ImageFlagName))
	}

	if opts.CacheSize != "" {
		// must parse as a resource quantity
		if _, err := resource.ParseQuantity(opts.CacheSize); err != nil {
			errs = errs.Also(cli.ErrInvalidValue(opts.CacheSize, cli.CacheSizeFlagName))
		}
	}

	// git-repo and local-path are mutually exclusive
	if opts.GitRepo == "" && opts.LocalPath == "" {
		errs = errs.Also(cli.ErrMissingOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName))
	} else if opts.GitRepo != "" && opts.LocalPath != "" {
		errs = errs.Also(cli.ErrMultipleOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName))
	}

	// git-revision is required for git-repo
	if opts.GitRepo != "" && opts.GitRevision == "" {
		errs = errs.Also(cli.ErrMissingField(cli.GitRevisionFlagName))
	}

	if opts.LocalPath != "" {
		if opts.SubPath != "" {
			// sub-path cannot be used with local-path
			errs = errs.Also(cli.ErrDisallowedFields(cli.SubPathFlagName))
		}
		if opts.CacheSize != "" {
			// cache-size cannot be used with local-path
			errs = errs.Also(cli.ErrDisallowedFields(cli.CacheSizeFlagName))
		}
	}

	// nothing to do for artifact, handler, and invoker

	if opts.Tail {
		if opts.WaitTimeout == "" {
			errs = errs.Also(cli.ErrMissingField(cli.WaitTimeoutFlagName))
		} else if _, err := time.ParseDuration(opts.WaitTimeout); err != nil {
			errs = errs.Also(cli.ErrInvalidValue(opts.WaitTimeout, cli.WaitTimeoutFlagName))
		}
	}

	if opts.DryRun && opts.Tail {
		errs = errs.Also(cli.ErrMultipleOneOf(cli.DryRunFlagName, cli.TailFlagName))
	}

	return errs
}

func (opts *FunctionCreateOptions) Exec(ctx context.Context, c *cli.Config) error {
	function := &buildv1alpha1.Function{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      opts.Name,
		},
		Spec: buildv1alpha1.FunctionSpec{
			Image:    opts.Image,
			Artifact: opts.Artifact,
			Handler:  opts.Handler,
			Invoker:  opts.Invoker,
		},
	}
	if opts.CacheSize != "" {
		quantity := resource.MustParse(opts.CacheSize)
		function.Spec.CacheSize = &quantity
	}
	if opts.GitRepo != "" {
		function.Spec.Source = &buildv1alpha1.Source{
			Git: &buildv1alpha1.GitSource{
				URL:      opts.GitRepo,
				Revision: opts.GitRevision,
			},
			SubPath: opts.SubPath,
		}
	}

	if opts.LocalPath != "" {
		targetImage := opts.Image
		if strings.HasPrefix(opts.Image, "_") {
			riffBuildConfig, err := c.Core().ConfigMaps(function.Namespace).Get("riff-build", metav1.GetOptions{})
			if err != nil {
				return err
			}
			targetImage, err = buildv1alpha1.ResolveDefaultImage(function, riffBuildConfig.Data["default-image-prefix"])
			if err != nil {
				return err
			}
		}
		builders, err := c.Core().ConfigMaps("riff-system").Get("builders", metav1.GetOptions{})
		if err != nil {
			return err
		}
		builder := builders.Data["riff-function"]
		if builder == "" {
			return fmt.Errorf("unknown builder for %q", "riff-function")
		}
		err = c.Pack.Build(ctx, pack.BuildOptions{
			Image:   targetImage,
			AppDir:  opts.LocalPath,
			Builder: builder,
			Env: map[string]string{
				"RIFF":          "true",
				"RIFF_ARTIFACT": opts.Artifact,
				"RIFF_HANDLER":  opts.Handler,
				"RIFF_OVERRIDE": opts.Invoker,
			},
			Publish: true,
		})
		if err != nil {
			return err
		}
	}

	if opts.DryRun {
		cli.DryRunResource(ctx, function, function.GetGroupVersionKind())
	} else {
		var err error
		function, err = c.Build().Functions(opts.Namespace).Create(function)
		if err != nil {
			return err
		}
	}
	c.Successf("Created function %q\n", function.Name)
	if opts.Tail {
		// err guarded by Validate()
		timeout, _ := time.ParseDuration(opts.WaitTimeout)
		err := race.Run(ctx, timeout,
			func(ctx context.Context) error {
				return k8s.WaitUntilReady(ctx, c.Build().RESTClient(), "functions", function)
			},
			func(ctx context.Context) error {
				return c.Kail.FunctionLogs(ctx, function, cli.TailSinceCreateDefault, c.Stdout)
			},
		)
		if err == context.DeadlineExceeded {
			c.Errorf("Timeout after %q waiting for %q to become ready\n", opts.WaitTimeout, opts.Name)
			c.Infof("To view status run: %s function list %s %s\n", c.Name, cli.NamespaceFlagName, opts.Namespace)
			c.Infof("To continue watching logs run: %s function tail %s %s %s\n", c.Name, opts.Name, cli.NamespaceFlagName, opts.Namespace)
			err = cli.SilenceError(err)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (opts *FunctionCreateOptions) IsDryRun() bool {
	return opts.DryRun
}

func NewFunctionCreateCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &FunctionCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a function from source",
		Long: strings.TrimSpace(`
Create a function from source using the function Cloud Native Buildpack builder.

Function source can be specified either as a Git repository or as a local
directory. Builds from Git are run in the cluster while builds from a local
directory are run inside a local Docker daemon and are orchestrated by this
command (in the future, builds from local source may also be run in the
cluster).

In addition to the source code, functions are defined by these properties:

- invoker - language runtime that should host the function, the invoker is often
    auto-detected, but may need to be specified in cases of ambiguity.
- artifact - file in the source that contains the function.
- handler - invoker specific, typically the method or class within the artifact.

These values can be versioned with the source code in a riff.toml file, or
specified here to override the source. Versioning with the source is preferred
as changed can be deployed as a unit. Overriding is necessary when deploying
multiple functions from a single code base.

The riff.toml file takes the form:

    override = "<invoker name>"
	artifact = "<path to artifact>"
	handler = "<function handler>"

`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s function create my-func %s registry.example.com/image %s https://example.com/my-func.git", c.Name, cli.ImageFlagName, cli.GitRepoFlagName),
			fmt.Sprintf("%s function create my-func %s registry.example.com/image %s ./my-func", c.Name, cli.ImageFlagName, cli.LocalPathFlagName),
		}, "\n"),
		Args: cli.Args(
			cli.NameArg(&opts.Name),
		),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().StringVar(&opts.Image, cli.StripDash(cli.ImageFlagName), "_", "`repository` where the built images are pushed")
	cmd.Flags().StringVar(&opts.CacheSize, cli.StripDash(cli.CacheSizeFlagName), "", "`size` of persistent volume to cache resources between builds")
	cmd.Flags().StringVar(&opts.Artifact, cli.StripDash(cli.ArtifactFlagName), "", "`file` containing the function within the build workspace (detected by default)")
	cmd.Flags().StringVar(&opts.Handler, cli.StripDash(cli.HandlerFlagName), "", "`name` of the method or class to invoke, depends on the invoker (detected by default)")
	cmd.Flags().StringVar(&opts.Invoker, cli.StripDash(cli.InvokerFlagName), "", "language runtime invoker `name` (detected by default)")
	cmd.Flags().StringVar(&opts.LocalPath, cli.StripDash(cli.LocalPathFlagName), "", "path to `directory` containing source code on the local machine")
	cmd.Flags().StringVar(&opts.GitRepo, cli.StripDash(cli.GitRepoFlagName), "", "git `url` to remote source code")
	cmd.Flags().StringVar(&opts.GitRevision, cli.StripDash(cli.GitRevisionFlagName), "master", "`refspec` within the git repo to checkout")
	cmd.Flags().StringVar(&opts.SubPath, cli.StripDash(cli.SubPathFlagName), "", "path to `directory` within the git repo to checkout")
	cmd.Flags().BoolVar(&opts.Tail, cli.StripDash(cli.TailFlagName), false, "watch build logs")
	cmd.Flags().StringVar(&opts.WaitTimeout, cli.StripDash(cli.WaitTimeoutFlagName), "10m", "`duration` to wait for the function to become ready when watching logs")
	cmd.Flags().BoolVar(&opts.DryRun, cli.StripDash(cli.DryRunFlagName), false, "print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr")

	return cmd
}
