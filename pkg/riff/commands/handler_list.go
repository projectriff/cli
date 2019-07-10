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
	"github.com/projectriff/cli/pkg/cli/printers"
	requestv1alpha1 "github.com/projectriff/system/pkg/apis/request/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type HandlerListOptions struct {
	cli.ListOptions
}

var (
	_ cli.Validatable = (*HandlerListOptions)(nil)
	_ cli.Executable  = (*HandlerListOptions)(nil)
)

func (opts *HandlerListOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.ListOptions.Validate(ctx))

	return errs
}

func (opts *HandlerListOptions) Exec(ctx context.Context, c *cli.Config) error {
	handlers, err := c.Request().Handlers(opts.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(handlers.Items) == 0 {
		c.Infof("No handlers found.\n")
		return nil
	}

	tablePrinter := printers.NewTablePrinter(printers.PrintOptions{
		WithNamespace: opts.AllNamespaces,
	}).With(func(h printers.PrintHandler) {
		columns := opts.printColumns()
		h.TableHandler(columns, opts.printList)
		h.TableHandler(columns, opts.print)
	})

	handlers = handlers.DeepCopy()
	cli.SortByNamespaceAndName(handlers.Items)

	return tablePrinter.PrintObj(handlers, c.Stdout)
}

func NewHandlerListCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &HandlerListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "table listing of handlers",
		Long: strings.TrimSpace(`
List handlers in a namespace or across all namespaces.

For detail regarding the status of a single handler, run:

	` + c.Name + ` handler status <handler-name>
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s handler list", c.Name),
			fmt.Sprintf("%s handler list %s", c.Name, cli.AllNamespacesFlagName),
		}, "\n"),
		Args:    cli.Args(),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.AllNamespacesFlag(cmd, c, &opts.Namespace, &opts.AllNamespaces)

	return cmd
}

func (opts *HandlerListOptions) printList(handlers *requestv1alpha1.HandlerList, printOpts printers.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(handlers.Items))
	for i := range handlers.Items {
		r, err := opts.print(&handlers.Items[i], printOpts)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func (opts *HandlerListOptions) print(handler *requestv1alpha1.Handler, _ printers.PrintOptions) ([]metav1beta1.TableRow, error) {
	now := time.Now()
	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: handler},
	}
	refType, refValue := opts.formatRef(handler)
	host := ""
	if handler.Status.URL != nil {
		host = handler.Status.URL.Host
	}
	row.Cells = append(row.Cells,
		handler.Name,
		refType,
		refValue,
		cli.FormatEmptyString(host),
		cli.FormatConditionStatus(handler.Status.GetCondition(requestv1alpha1.HandlerConditionReady)),
		cli.FormatTimestampSince(handler.CreationTimestamp, now),
	)
	return []metav1beta1.TableRow{row}, nil
}

func (opts *HandlerListOptions) printColumns() []metav1beta1.TableColumnDefinition {
	return []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Type", Type: "string"},
		{Name: "Ref", Type: "string"},
		{Name: "Host", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
	}
}

func (opts *HandlerListOptions) formatRef(handler *requestv1alpha1.Handler) (string, string) {
	if handler.Spec.Build != nil {
		if handler.Spec.Build.ApplicationRef != "" {
			return "application", handler.Spec.Build.ApplicationRef
		}
		if handler.Spec.Build.FunctionRef != "" {
			return "function", handler.Spec.Build.FunctionRef
		}
	} else if handler.Spec.Template != nil && handler.Spec.Template.Containers[0].Image != "" {
		return "image", handler.Spec.Template.Containers[0].Image
	}
	return cli.Swarnf("<unknown>"), cli.Swarnf("<unknown>")
}
