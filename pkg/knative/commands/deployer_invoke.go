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

type DeployerInvokeOptions struct {
	cli.ResourceOptions

	ContentTypeJSON bool
	ContentTypeText bool
	Path            string
	BareArgs        []string
}

var (
	_ cli.Validatable = (*DeployerInvokeOptions)(nil)
	_ cli.Executable  = (*DeployerInvokeOptions)(nil)
)

func (opts *DeployerInvokeOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	if opts.ContentTypeJSON && opts.ContentTypeText {
		errs = errs.Also(cli.ErrMultipleOneOf(cli.JSONFlagName, cli.TextFlagName))
	}

	return errs
}

func (opts *DeployerInvokeOptions) Exec(ctx context.Context, c *cli.Config) error {
	deployer, err := c.KnativeRuntime().Deployers(opts.Namespace).Get(opts.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if !deployer.Status.IsReady() || deployer.Status.URL == nil || deployer.Status.URL.Host == "" {
		return fmt.Errorf("deployer %q is not ready", opts.Name)
	}

	ingress, err := opts.ingressServiceHost(c)
	if err != nil {
		return err
	}

	curlArgs := []string{ingress + opts.Path, "-H", fmt.Sprintf("Host: %s", deployer.Status.URL.Host)}
	if opts.ContentTypeJSON {
		curlArgs = append(curlArgs, "-H", "Content-Type: application/json")
	}
	if opts.ContentTypeText {
		curlArgs = append(curlArgs, "-H", "Content-Type: text/plain")
	}
	curlArgs = append(curlArgs, opts.BareArgs...)

	curl := c.Exec(context.Background(), "curl", curlArgs...)

	curl.Stdin = c.Stdin
	curl.Stdout = c.Stdout
	curl.Stderr = c.Stderr

	return curl.Run()
}

func NewDeployerInvokeCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &DeployerInvokeOptions{}

	cmd := &cobra.Command{
		Use:    "invoke",
		Hidden: true,
		Short:  "invoke an http request deployer using curl",
		Long: strings.TrimSpace(`
This command is not supported and may be removed in the future.
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s knative deployer invoke my-deployer", c.Name),
			fmt.Sprintf("%s knative deployer invoke my-deployer --text -- -d 'hello' -w '\\n'", c.Name),
			fmt.Sprintf("%s knative deployer invoke my-deployer /request/path", c.Name),
		}, "\n"),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.Args(cmd,
		cli.NameArg(&opts.Name),
		cli.Arg{
			Name:     "path",
			Arity:    1,
			Optional: true,
			Set: func(cmd *cobra.Command, args []string, offset int) error {
				if offset >= cmd.ArgsLenAtDash() && cmd.ArgsLenAtDash() != -1 {
					return cli.ErrIgnoreArg
				}
				opts.Path = args[offset]
				return nil
			},
		},
		cli.BareDoubleDashArgs(&opts.BareArgs),
	)

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().BoolVar(&opts.ContentTypeJSON, cli.StripDash(cli.JSONFlagName), false, "set the content type to application/json")
	cmd.Flags().BoolVar(&opts.ContentTypeText, cli.StripDash(cli.TextFlagName), false, "set the content type to text/plain")

	return cmd
}

func (opts *DeployerInvokeOptions) ingressServiceHost(c *cli.Config) (string, error) {
	// TODO allow setting ingress manually
	svc, err := c.Core().Services("istio-system").Get("istio-ingressgateway", metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	ingress := ""
	if svc.Spec.Type == "LoadBalancer" {
		ingresses := svc.Status.LoadBalancer.Ingress
		if len(ingresses) > 0 {
			ingress = ingresses[0].IP
			if ingress == "" {
				ingress = ingresses[0].Hostname
			}
		}
	}
	if ingress == "" {
		for _, port := range svc.Spec.Ports {
			if port.Name == "http" || port.Name == "http2" {
				config := c.KubeRestConfig()
				host := config.Host[0:strings.LastIndex(config.Host, ":")]
				host = strings.Replace(host, "https", "http", 1)
				ingress = fmt.Sprintf("%s:%d", host, port.NodePort)
				break
			}
		}
	}
	if ingress == "" {
		return "", fmt.Errorf("ingress not available")
	}

	return ingress, nil
}
