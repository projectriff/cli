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
	"github.com/projectriff/system/pkg/apis"
	"github.com/spf13/cobra"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ImageStatusOptions struct {
	options.ResourceOptions
}

var (
	_ cli.Validatable = (*ImageStatusOptions)(nil)
	_ cli.Executable  = (*ImageStatusOptions)(nil)
)

func (opts *ImageStatusOptions) Validate(ctx context.Context) cli.FieldErrors {
	errs := cli.FieldErrors{}

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	return errs
}

func (opts *ImageStatusOptions) Exec(ctx context.Context, c *cli.Config) error {
	image, err := c.Bindings().ImageBindings(opts.Namespace).Get(opts.Name, metav1.GetOptions{})
	if err != nil {
		if !apierrs.IsNotFound(err) {
			return err
		}
		c.Errorf("Image binding %q not found\n", fmt.Sprintf("%s/%s", opts.Namespace, opts.Name))
		return cli.SilenceError(err)
	}

	ready := image.Status.GetCondition(apis.ConditionReady)
	cli.PrintResourceStatus(c, image.Name, &apis.Condition{
		Type:    apis.ConditionReady,
		Message: ready.Message,
		Reason:  ready.Reason,
		Status:  ready.Status,
		LastTransitionTime: apis.VolatileTime{
			Inner: ready.LastTransitionTime.Inner,
		},
	})

	return nil
}

func NewImageStatusCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &ImageStatusOptions{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "show image binding status",
		Long: strings.TrimSpace(`
Display status details for an image binding.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
image roll out is processed.
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s binding image status my-imagebinding", c.Name),
		}, "\n"),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.Args(cmd,
		cli.NameArg(&opts.Name),
	)

	cli.NamespaceFlag(cmd, c, &opts.Namespace)

	return cmd
}
