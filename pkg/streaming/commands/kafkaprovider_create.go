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

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/cli/options"
	streamv1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KafkaProviderCreateOptions struct {
	options.ResourceOptions

	BootstrapServers string

	DryRun bool
}

var (
	_ cli.Validatable = (*KafkaProviderCreateOptions)(nil)
	_ cli.Executable  = (*KafkaProviderCreateOptions)(nil)
	_ cli.DryRunable  = (*KafkaProviderCreateOptions)(nil)
)

func (opts *KafkaProviderCreateOptions) Validate(ctx context.Context) cli.FieldErrors {
	errs := cli.FieldErrors{}

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	if opts.BootstrapServers == "" {
		errs = errs.Also(cli.ErrMissingField(cli.BootstrapServersFlagName))
	}

	return errs
}

func (opts *KafkaProviderCreateOptions) Exec(ctx context.Context, c *cli.Config) error {
	provider := &streamv1alpha1.KafkaProvider{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      opts.Name,
		},
		Spec: streamv1alpha1.KafkaProviderSpec{
			BootstrapServers: opts.BootstrapServers,
		},
	}

	if opts.DryRun {
		cli.DryRunResource(ctx, provider, provider.GetGroupVersionKind())
	} else {
		var err error
		provider, err = c.StreamingRuntime().KafkaProviders(opts.Namespace).Create(provider)
		if err != nil {
			return err
		}
	}
	c.Successf("Created kafka provider %q\n", provider.Name)
	return nil
}

func (opts *KafkaProviderCreateOptions) IsDryRun() bool {
	return opts.DryRun
}

func NewKafkaProviderCreateCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &KafkaProviderCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a kafka provider of messages",
		Long: strings.TrimSpace(`
<todo>
`),
		Example: fmt.Sprintf("%s streaming kafka-provider create my-kafka-provider %s kafka.local:9092", c.Name, cli.BootstrapServersFlagName),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.Args(cmd,
		cli.NameArg(&opts.Name),
	)

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().StringVar(&opts.BootstrapServers, cli.StripDash(cli.BootstrapServersFlagName), "", "`address` of the kafka broker")
	cmd.Flags().BoolVar(&opts.DryRun, cli.StripDash(cli.DryRunFlagName), false, "print kubernetes resources to stdout rather than apply them to the cluster, messages normally on stdout will be sent to stderr")

	return cmd
}
