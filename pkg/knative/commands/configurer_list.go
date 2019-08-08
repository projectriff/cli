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
	knativev1alpha1 "github.com/projectriff/system/pkg/apis/knative/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ConfigurerListOptions struct {
	cli.ListOptions
}

var (
	_ cli.Validatable = (*ConfigurerListOptions)(nil)
	_ cli.Executable  = (*ConfigurerListOptions)(nil)
)

func (opts *ConfigurerListOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.ListOptions.Validate(ctx))

	return errs
}

func (opts *ConfigurerListOptions) Exec(ctx context.Context, c *cli.Config) error {
	configurers, err := c.KnativeRuntime().Configurers(opts.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(configurers.Items) == 0 {
		c.Infof("No configurers found.\n")
		return nil
	}

	tablePrinter := printers.NewTablePrinter(printers.PrintOptions{
		WithNamespace: opts.AllNamespaces,
	}).With(func(h printers.PrintHandler) {
		columns := opts.printColumns()
		h.TableHandler(columns, opts.printList)
		h.TableHandler(columns, opts.print)
	})

	configurers = configurers.DeepCopy()
	cli.SortByNamespaceAndName(configurers.Items)

	return tablePrinter.PrintObj(configurers, c.Stdout)
}

func NewConfigurerListCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &ConfigurerListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "table listing of configurers",
		Long: strings.TrimSpace(`
List configurers in a namespace or across all namespaces.

For detail regarding the status of a single configurer, run:

    ` + c.Name + ` knative configurer status <configurer-name>
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s knative configurer list", c.Name),
			fmt.Sprintf("%s knative configurer list %s", c.Name, cli.AllNamespacesFlagName),
		}, "\n"),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.AllNamespacesFlag(cmd, c, &opts.Namespace, &opts.AllNamespaces)

	return cmd
}

func (opts *ConfigurerListOptions) printList(configurers *knativev1alpha1.ConfigurerList, printOpts printers.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(configurers.Items))
	for i := range configurers.Items {
		r, err := opts.print(&configurers.Items[i], printOpts)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func (opts *ConfigurerListOptions) print(configurer *knativev1alpha1.Configurer, _ printers.PrintOptions) ([]metav1beta1.TableRow, error) {
	now := time.Now()
	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: configurer},
	}
	refType, refValue := opts.formatRef(configurer)
	host := ""
	if configurer.Status.URL != nil {
		host = configurer.Status.URL.Host
	}
	row.Cells = append(row.Cells,
		configurer.Name,
		refType,
		refValue,
		cli.FormatEmptyString(host),
		cli.FormatConditionStatus(configurer.Status.GetCondition(knativev1alpha1.ConfigurerConditionReady)),
		cli.FormatTimestampSince(configurer.CreationTimestamp, now),
	)
	return []metav1beta1.TableRow{row}, nil
}

func (opts *ConfigurerListOptions) printColumns() []metav1beta1.TableColumnDefinition {
	return []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Type", Type: "string"},
		{Name: "Ref", Type: "string"},
		{Name: "Host", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
	}
}

func (opts *ConfigurerListOptions) formatRef(configurer *knativev1alpha1.Configurer) (string, string) {
	if configurer.Spec.Build != nil {
		if configurer.Spec.Build.ApplicationRef != "" {
			return "application", configurer.Spec.Build.ApplicationRef
		}
		if configurer.Spec.Build.FunctionRef != "" {
			return "function", configurer.Spec.Build.FunctionRef
		}
		if configurer.Spec.Build.ContainerRef != "" {
			return "container", configurer.Spec.Build.ContainerRef
		}
	} else if configurer.Spec.Template != nil && configurer.Spec.Template.Containers[0].Image != "" {
		return "image", configurer.Spec.Template.Containers[0].Image
	}
	return cli.Swarnf("<unknown>"), cli.Swarnf("<unknown>")
}
