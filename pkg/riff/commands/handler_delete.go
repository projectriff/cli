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

type HandlerDeleteOptions struct {
	cli.DeleteOptions
}

var (
	_ cli.Validatable = (*HandlerDeleteOptions)(nil)
	_ cli.Executable  = (*HandlerDeleteOptions)(nil)
)

func (opts *HandlerDeleteOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.DeleteOptions.Validate(ctx))

	return errs
}

func (opts *HandlerDeleteOptions) Exec(ctx context.Context, c *cli.Config) error {
	client := c.Request().Handlers(opts.Namespace)

	if opts.All {
		if err := client.DeleteCollection(nil, metav1.ListOptions{}); err != nil {
			return err
		}
		c.Successf("Deleted handlers in namespace %q\n", opts.Namespace)
		return nil
	}

	for _, name := range opts.Names {
		if err := client.Delete(name, nil); err != nil {
			return err
		}
		c.Successf("Deleted handler %q\n", name)
	}

	return nil
}

func NewHandlerDeleteCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &HandlerDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete handler(s)",
		Long: strings.TrimSpace(`
Delete one or more handlers by name or all handlers within a namespace.

New HTTP requests addressed to the handler will fail. A new handler created with
the same name will start to receive new HTTP requests addressed to the same
handler.
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s handler delete my-handler", c.Name),
			fmt.Sprintf("%s handler delete %s ", c.Name, cli.AllFlagName),
		}, "\n"),
		Args: cli.Args(
			cli.NamesArg(&opts.Names),
		),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().BoolVar(&opts.All, cli.StripDash(cli.AllFlagName), false, "delete all handlers within the namespace")

	return cmd
}
