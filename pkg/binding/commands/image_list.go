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

	bindingsv1alpha1 "github.com/projectriff/bindings/pkg/apis/bindings/v1alpha1"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/cli/options"
	"github.com/projectriff/cli/pkg/cli/printers"
	"github.com/projectriff/system/pkg/apis"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	knapis "knative.dev/pkg/apis"
)

type ImageListOptions struct {
	options.ListOptions
}

var (
	_ cli.Validatable = (*ImageListOptions)(nil)
	_ cli.Executable  = (*ImageListOptions)(nil)
)

func (opts *ImageListOptions) Validate(ctx context.Context) cli.FieldErrors {
	errs := cli.FieldErrors{}

	errs = errs.Also(opts.ListOptions.Validate(ctx))

	return errs
}

func (opts *ImageListOptions) Exec(ctx context.Context, c *cli.Config) error {
	images, err := c.Bindings().ImageBindings(opts.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(images.Items) == 0 {
		c.Infof("No image bindings found.\n")
		return nil
	}

	tablePrinter := printers.NewTablePrinter(printers.PrintOptions{
		WithNamespace: opts.AllNamespaces,
	}).With(func(h printers.PrintHandler) {
		columns := opts.printColumns()
		h.TableHandler(columns, opts.printList)
		h.TableHandler(columns, opts.print)
	})

	images = images.DeepCopy()
	cli.SortByNamespaceAndName(images.Items)

	return tablePrinter.PrintObj(images, c.Stdout)
}

func NewImageListCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &ImageListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "table listing of image bindings",
		Long: strings.TrimSpace(`
List image bindings in a namespace or across all namespaces.

For detail regarding the status of a single image, run:

    ` + c.Name + ` binding image status <image-binding-name>
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s binding image list", c.Name),
			fmt.Sprintf("%s binding image list %s", c.Name, cli.AllNamespacesFlagName),
		}, "\n"),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.AllNamespacesFlag(cmd, c, &opts.Namespace, &opts.AllNamespaces)

	return cmd
}

func (opts *ImageListOptions) printList(images *bindingsv1alpha1.ImageBindingList, printOpts printers.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(images.Items))
	for i := range images.Items {
		r, err := opts.print(&images.Items[i], printOpts)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func (opts *ImageListOptions) print(image *bindingsv1alpha1.ImageBinding, _ printers.PrintOptions) ([]metav1beta1.TableRow, error) {
	now := time.Now()
	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: image},
	}
	condition := image.Status.GetCondition(bindingsv1alpha1.ImageBindingConditionReady)
	if condition == nil {
		condition = &knapis.Condition{}
	}
	// TODO use discovery client to get resource names for group/kind
	var subject, provider string
	if image.Spec.Subject != nil {
		subject = fmt.Sprintf("%ss.%s:%s", strings.ToLower(image.Spec.Subject.Kind), strings.Split(image.Spec.Subject.APIVersion, "/")[0], image.Spec.Subject.Name)
	}
	if image.Spec.Provider != nil {
		provider = fmt.Sprintf("%ss.%s:%s", strings.ToLower(image.Spec.Provider.Kind), strings.Split(image.Spec.Provider.APIVersion, "/")[0], image.Spec.Provider.Name)
	}
	row.Cells = append(row.Cells,
		image.Name,
		cli.FormatEmptyString(subject),
		cli.FormatEmptyString(provider),
		cli.FormatEmptyString(image.Spec.ContainerName),
		cli.FormatConditionStatus(&apis.Condition{
			Type:   apis.ConditionReady,
			Reason: condition.Reason,
			Status: condition.Status,
		}),
		cli.FormatTimestampSince(image.CreationTimestamp, now),
	)
	return []metav1beta1.TableRow{row}, nil
}

func (opts *ImageListOptions) printColumns() []metav1beta1.TableColumnDefinition {
	return []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Subject", Type: "string"},
		{Name: "Providers", Type: "string"},
		{Name: "Container Name", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
	}
}
