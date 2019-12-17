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
	"github.com/projectriff/cli/pkg/cli/options"
	"github.com/projectriff/cli/pkg/k8s"
	"github.com/projectriff/cli/pkg/race"
	streamv1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PulsarProviderCreateOptions struct {
	options.ResourceOptions

	ServiceURL string

	DryRun bool

	Tail        bool
	WaitTimeout time.Duration
}

var (
	_ cli.Validatable = (*PulsarProviderCreateOptions)(nil)
	_ cli.Executable  = (*PulsarProviderCreateOptions)(nil)
	_ cli.DryRunable  = (*PulsarProviderCreateOptions)(nil)
)

func (opts *PulsarProviderCreateOptions) Validate(ctx context.Context) cli.FieldErrors {
	errs := cli.FieldErrors{}

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	if opts.ServiceURL == "" {
		errs = errs.Also(cli.ErrMissingField(cli.ServiceURLFlagName))
	}
	if opts.DryRun && opts.Tail {
		errs = errs.Also(cli.ErrMultipleOneOf(cli.DryRunFlagName, cli.TailFlagName))
	}
	if opts.WaitTimeout < 0 {
		errs = errs.Also(cli.ErrInvalidValue(opts.WaitTimeout, cli.WaitTimeoutFlagName))
	}

	return errs
}

func (opts *PulsarProviderCreateOptions) Exec(ctx context.Context, c *cli.Config) error {
	provider := &streamv1alpha1.PulsarProvider{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      opts.Name,
		},
		Spec: streamv1alpha1.PulsarProviderSpec{
			ServiceURL: opts.ServiceURL,
		},
	}

	if opts.DryRun {
		cli.DryRunResource(ctx, provider, provider.GetGroupVersionKind())
	} else {
		var err error
		provider, err = c.StreamingRuntime().PulsarProviders(opts.Namespace).Create(provider)
		if err != nil {
			return err
		}
	}
	c.Successf("Created pulsar provider %q\n", provider.Name)
	if opts.Tail {
		err := race.Run(ctx, opts.WaitTimeout,
			func(ctx context.Context) error {
				return k8s.WaitUntilReady(ctx, c.StreamingRuntime().RESTClient(), "pulsarproviders", provider)
			},
			func(ctx context.Context) error {
				return c.Kail.PulsarProviderLogs(ctx, provider, cli.TailSinceCreateDefault, c.Stdout)
			},
		)
		if err == context.DeadlineExceeded {
			c.Errorf("Timeout after %q waiting for %q to become ready\n", opts.WaitTimeout, opts.Name)
			c.Infof("To view status run: %s streaming pulsar-provider list %s %s\n", c.Name, cli.NamespaceFlagName, opts.Namespace)
			err = cli.SilenceError(err)
		}
		if err != nil {
			return err
		}
		c.Successf("PulsarProvider %q is ready\n", provider.Name)
	}

	return nil
}

func (opts *PulsarProviderCreateOptions) IsDryRun() bool {
	return opts.DryRun
}

func NewPulsarProviderCreateCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &PulsarProviderCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a pulsar provider of messages",
		Long: strings.TrimSpace(`
<todo>
`),
		Example: fmt.Sprintf("%s streaming pulsar-provider create my-pulsar-provider %s pulsar://localhost:6650", c.Name, cli.ServiceURLFlagName),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.Args(cmd,
		cli.NameArg(&opts.Name),
	)

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().StringVar(&opts.ServiceURL, cli.StripDash(cli.ServiceURLFlagName), "", "`url` of the pulsar service")
	cmd.Flags().BoolVar(&opts.DryRun, cli.StripDash(cli.DryRunFlagName), false, "print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr")
	cmd.Flags().BoolVar(&opts.Tail, cli.StripDash(cli.TailFlagName), false, "watch creation progress")
	cmd.Flags().DurationVar(&opts.WaitTimeout, cli.StripDash(cli.WaitTimeoutFlagName), time.Minute*1, "`duration` to wait for the provider to become ready when watching progress")

	return cmd
}
