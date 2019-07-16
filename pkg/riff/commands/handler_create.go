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

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/k8s"
	"github.com/projectriff/cli/pkg/parsers"
	"github.com/projectriff/cli/pkg/race"
	"github.com/projectriff/cli/pkg/validation"
	requestv1alpha1 "github.com/projectriff/system/pkg/apis/request/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HandlerCreateOptions struct {
	cli.ResourceOptions

	Image          string
	ApplicationRef string
	FunctionRef    string

	Env     []string
	EnvFrom []string

	Tail        bool
	WaitTimeout string

	DryRun bool
}

var (
	_ cli.Validatable = (*HandlerCreateOptions)(nil)
	_ cli.Executable  = (*HandlerCreateOptions)(nil)
	_ cli.DryRunable  = (*HandlerCreateOptions)(nil)
)

func (opts *HandlerCreateOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.ResourceOptions.Validate((ctx)))

	// application-ref, build-ref and image are mutually exclusive
	used := []string{}
	unused := []string{}

	if opts.ApplicationRef != "" {
		used = append(used, cli.ApplicationRefFlagName)
	} else {
		unused = append(unused, cli.ApplicationRefFlagName)
	}

	if opts.FunctionRef != "" {
		used = append(used, cli.FunctionRefFlagName)
	} else {
		unused = append(unused, cli.FunctionRefFlagName)
	}

	if opts.Image != "" {
		used = append(used, cli.ImageFlagName)
	} else {
		unused = append(unused, cli.ImageFlagName)
	}

	if len(used) == 0 {
		errs = errs.Also(cli.ErrMissingOneOf(unused...))
	} else if len(used) > 1 {
		errs = errs.Also(cli.ErrMultipleOneOf(used...))
	}

	errs = errs.Also(validation.EnvVars(opts.Env, cli.EnvFlagName))
	errs = errs.Also(validation.EnvVarFroms(opts.EnvFrom, cli.EnvFromFlagName))

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

func (opts *HandlerCreateOptions) Exec(ctx context.Context, c *cli.Config) error {
	handler := &requestv1alpha1.Handler{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      opts.Name,
		},
		Spec: requestv1alpha1.HandlerSpec{
			Template: &corev1.PodSpec{
				Containers: []corev1.Container{{}},
			},
		},
	}

	if opts.ApplicationRef != "" {
		handler.Spec.Build = &requestv1alpha1.Build{
			ApplicationRef: opts.ApplicationRef,
		}
	}
	if opts.FunctionRef != "" {
		handler.Spec.Build = &requestv1alpha1.Build{
			FunctionRef: opts.FunctionRef,
		}
	}
	if opts.Image != "" {
		handler.Spec.Template.Containers[0].Image = opts.Image
	}

	for _, env := range opts.Env {
		if handler.Spec.Template.Containers[0].Env == nil {
			handler.Spec.Template.Containers[0].Env = []corev1.EnvVar{}
		}
		handler.Spec.Template.Containers[0].Env = append(handler.Spec.Template.Containers[0].Env, parsers.EnvVar(env))
	}
	for _, env := range opts.EnvFrom {
		if handler.Spec.Template.Containers[0].Env == nil {
			handler.Spec.Template.Containers[0].Env = []corev1.EnvVar{}
		}
		handler.Spec.Template.Containers[0].Env = append(handler.Spec.Template.Containers[0].Env, parsers.EnvVarFrom(env))
	}

	if opts.DryRun {
		cli.DryRunResource(ctx, handler, handler.GetGroupVersionKind())
	} else {
		var err error
		handler, err = c.Request().Handlers(opts.Namespace).Create(handler)
		if err != nil {
			return err
		}
	}
	c.Successf("Created handler %q\n", handler.Name)
	if opts.Tail {
		// err guarded by Validate()
		timeout, _ := time.ParseDuration(opts.WaitTimeout)
		err := race.Run(ctx, timeout,
			func(ctx context.Context) error {
				return k8s.WaitUntilReady(ctx, c.Request().RESTClient(), "handlers", handler)
			},
			func(ctx context.Context) error {
				return c.Kail.HandlerLogs(ctx, handler, cli.TailSinceCreateDefault, c.Stdout)
			},
		)
		if err == context.DeadlineExceeded {
			c.Errorf("Timeout after %q waiting for %q to become ready\n", opts.WaitTimeout, opts.Name)
			c.Infof("To view status run: %s handler list %s %s\n", c.Name, cli.NamespaceFlagName, opts.Namespace)
			c.Infof("To continue watching logs run: %s handler tail %s %s %s\n", c.Name, opts.Name, cli.NamespaceFlagName, opts.Namespace)
			err = cli.SilenceError(err)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (opts *HandlerCreateOptions) IsDryRun() bool {
	return opts.DryRun
}

func NewHandlerCreateCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &HandlerCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a handler to map HTTP requests to an application, function or image",
		Long: strings.TrimSpace(`
Create an HTTP request handler.

There are three way to create a handler:
- from an application reference
- from a function reference
- from an image

Application and function references are resolved within the same namespace as
the handler. As the build produces new images, the image will roll out
automatically.

Image based handlers must be updated manually to roll out new images.

The runtime environment can be configured by ` + cli.EnvFlagName + ` for static key-value pairs
and ` + cli.EnvFromFlagName + ` to map values from a ConfigMap or Secret.
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s handler create my-app-handler %s my-app", c.Name, cli.ApplicationRefFlagName),
			fmt.Sprintf("%s handler create my-func-handler %s my-func", c.Name, cli.FunctionRefFlagName),
			fmt.Sprintf("%s handler create my-image-handler %s registry.example.com/my-image:latest", c.Name, cli.ImageFlagName),
		}, "\n"),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.Args(cmd,
		cli.NameArg(&opts.Name),
	)

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().StringVar(&opts.Image, cli.StripDash(cli.ImageFlagName), "", "container `image` to deploy")
	cmd.Flags().StringVar(&opts.ApplicationRef, cli.StripDash(cli.ApplicationRefFlagName), "", "`name` of application to deploy")
	cmd.Flags().StringVar(&opts.FunctionRef, cli.StripDash(cli.FunctionRefFlagName), "", "`name` of function to deploy")
	cmd.Flags().StringArrayVar(&opts.Env, cli.StripDash(cli.EnvFlagName), []string{}, fmt.Sprintf("environment `variable` defined as a key value pair separated by an equals sign, example %q (may be set multiple times)", fmt.Sprintf("%s MY_VAR=my-value", cli.EnvFlagName)))
	cmd.Flags().StringArrayVar(&opts.EnvFrom, cli.StripDash(cli.EnvFromFlagName), []string{}, fmt.Sprintf("environment `variable` from a config map or secret, example %q, %q (may be set multiple times)", fmt.Sprintf("%s MY_SECRET_VALUE=secretKeyRef:my-secret-name:key-in-secret", cli.EnvFromFlagName), fmt.Sprintf("%s MY_CONFIG_MAP_VALUE=configMapKeyRef:my-config-map-name:key-in-config-map", cli.EnvFromFlagName)))
	cmd.Flags().BoolVar(&opts.Tail, cli.StripDash(cli.TailFlagName), false, "watch handler logs")
	cmd.Flags().StringVar(&opts.WaitTimeout, cli.StripDash(cli.WaitTimeoutFlagName), "10m", "`duration` to wait for the handler to become ready when watching logs")
	cmd.Flags().BoolVar(&opts.DryRun, cli.StripDash(cli.DryRunFlagName), false, "print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr")

	return cmd
}
