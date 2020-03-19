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

	bindingsv1alpha1 "github.com/projectriff/bindings/pkg/apis/bindings/v1alpha1"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/cli/options"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/tracker"
)

type ImageCreateOptions struct {
	options.ResourceOptions

	Subject   string
	Providers []string

	DryRun bool
}

var (
	_ cli.Validatable = (*ImageCreateOptions)(nil)
	_ cli.Executable  = (*ImageCreateOptions)(nil)
	_ cli.DryRunable  = (*ImageCreateOptions)(nil)
)

func (opts *ImageCreateOptions) Validate(ctx context.Context) cli.FieldErrors {
	errs := cli.FieldErrors{}

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	if chunks := strings.Split(opts.Subject, ":"); len(chunks) != 2 {
		errs = errs.Also(cli.ErrInvalidValue(opts.Subject, cli.SubjectFlagName))
	}

	if len(opts.Providers) == 0 {
		errs = errs.Also(cli.ErrMissingField(cli.ProvidersFlagName))
	}
	for i := range opts.Providers {
		if chunks := strings.Split(opts.Providers[i], ":"); len(chunks) != 3 {
			errs = errs.Also(cli.ErrInvalidValue(opts.Providers[i], cli.CurrentField).ViaFieldIndex(cli.ProvidersFlagName, i))
		}
	}

	return errs
}

func (opts *ImageCreateOptions) Exec(ctx context.Context, c *cli.Config) error {
	image := &bindingsv1alpha1.ImageBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      opts.Name,
		},
		Spec: bindingsv1alpha1.ImageBindingSpec{
			Providers: make([]bindingsv1alpha1.ImageProvider, len(opts.Providers)),
		},
	}

	resources, err := c.Discovery().ServerResources()
	if err != nil {
		return err
	}

	for i := range opts.Providers {
		apiVersion, kind, name, container, err := opts.ResolveObjectRef(resources, opts.Providers[i])
		if err != nil {
			return err
		}
		image.Spec.Providers[i] = bindingsv1alpha1.ImageProvider{
			ImageableRef: &tracker.Reference{
				APIVersion: apiVersion,
				Kind:       kind,
				Namespace:  opts.Namespace,
				Name:       name,
			},
			ContainerName: container,
		}
	}
	apiVersion, kind, name, _, err := opts.ResolveObjectRef(resources, opts.Subject)
	if err != nil {
		return err
	}
	image.Spec.Subject = &tracker.Reference{
		APIVersion: apiVersion,
		Kind:       kind,
		Namespace:  opts.Namespace,
		Name:       name,
	}

	if opts.DryRun {
		cli.DryRunResource(ctx, image, image.GetGroupVersionKind())
	} else {
		var err error
		image, err = c.Bindings().ImageBindings(opts.Namespace).Create(image)
		if err != nil {
			return err
		}
	}
	c.Successf("Created image binding %q\n", image.Name)
	return nil
}

func (opts *ImageCreateOptions) ResolveObjectRef(resources []*metav1.APIResourceList, ref string) (apiVersion, kind, name, container string, err error) {
	chunks := strings.Split(ref, ":")

	name = chunks[1]
	if len(chunks) > 2 {
		container = chunks[2]
	}

	// tease out apiVersion and kind
	target := fmt.Sprintf("%s/", chunks[0])
	for _, rl := range resources {
		for _, r := range rl.APIResources {
			fullname := fmt.Sprintf("%s.%s", r.Name, rl.GroupVersion)
			if strings.HasPrefix(fullname, target) {
				apiVersion = rl.GroupVersion
				kind = r.Kind
				return
			}
		}
	}

	err = fmt.Errorf("the server doesn't have a resource type %q", target)
	return
}

func (opts *ImageCreateOptions) IsDryRun() bool {
	return opts.DryRun
}

func NewImageCreateCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &ImageCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a image to deploy a workload",
		Long: strings.TrimSpace(`
Create an image binding.

<todo>
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s binding image create my-image-binding", c.Name),
		}, "\n"),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.Args(cmd,
		cli.NameArg(&opts.Name),
	)

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().StringVar(&opts.Subject, cli.StripDash(cli.SubjectFlagName), "", "subject `object reference` to inject images into")
	cmd.Flags().StringArrayVar(&opts.Providers, cli.StripDash(cli.ProvidersFlagName), []string{}, "provider `object reference` to get images from per container")
	cmd.Flags().BoolVar(&opts.DryRun, cli.StripDash(cli.DryRunFlagName), false, "print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr")

	return cmd
}
