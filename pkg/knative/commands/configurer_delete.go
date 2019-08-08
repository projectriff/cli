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
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigurerDeleteOptions struct {
	cli.DeleteOptions
}

var (
	_ cli.Validatable = (*ConfigurerDeleteOptions)(nil)
	_ cli.Executable  = (*ConfigurerDeleteOptions)(nil)
)

func (opts *ConfigurerDeleteOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.DeleteOptions.Validate(ctx))

	return errs
}

func (opts *ConfigurerDeleteOptions) Exec(ctx context.Context, c *cli.Config) error {
	client := c.KnativeRuntime().Configurers(opts.Namespace)

	if opts.All {
		if err := client.DeleteCollection(nil, metav1.ListOptions{}); err != nil {
			return err
		}
		c.Successf("Deleted configurers in namespace %q\n", opts.Namespace)
		return nil
	}

	for _, name := range opts.Names {
		if err := client.Delete(name, nil); err != nil {
			return err
		}
		c.Successf("Deleted configurer %q\n", name)
	}

	return nil
}

func NewConfigurerDeleteCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &ConfigurerDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete configurer(s)",
		Long: strings.TrimSpace(`
Delete one or more configurers by name or all configurers within a namespace.

New HTTP requests addressed to the configurer will fail. A new configurer created with
the same name will start to receive new HTTP requests addressed to the same
configurer.
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s knative configurer delete my-configurer", c.Name),
			fmt.Sprintf("%s knative configurer delete %s", c.Name, cli.AllFlagName),
		}, "\n"),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.Args(cmd,
		cli.NamesArg(&opts.Names),
	)

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().BoolVar(&opts.All, cli.StripDash(cli.AllFlagName), false, "delete all configurers within the namespace")

	return cmd
}
