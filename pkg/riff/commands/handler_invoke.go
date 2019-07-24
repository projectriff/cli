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
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/system/pkg/apis/request"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type HandlerInvokeOptions struct {
	cli.ResourceOptions

	ContentTypeJSON bool
	ContentTypeText bool
	Path            string
	BareArgs        []string
}

var (
	_ cli.Validatable = (*HandlerInvokeOptions)(nil)
	_ cli.Executable  = (*HandlerInvokeOptions)(nil)
)

func (opts *HandlerInvokeOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	if opts.ContentTypeJSON && opts.ContentTypeText {
		errs = errs.Also(cli.ErrMultipleOneOf(cli.JSONFlagName, cli.TextFlagName))
	}

	return errs
}

func (opts *HandlerInvokeOptions) Exec(ctx context.Context, c *cli.Config) error {
	handler, err := c.Request().Handlers(opts.Namespace).Get(opts.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if !handler.Status.IsReady() || handler.Status.ServiceName == "" {
		return fmt.Errorf("handler %q is not ready", opts.Name)
	}

	port, err := opts.findAvailablePort()
	if err != nil {
		return err
	}

	pod, err := opts.findReadyPod(ctx, c)
	if err != nil {
		return err
	}

	stopChan := make(chan struct{}, 1)
	defer close(stopChan)
	err = opts.forwardPort(ctx, c, pod, port, stopChan)
	if err != nil {
		return err
	}

	return opts.invoke(ctx, c, port)
}

func NewHandlerInvokeCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &HandlerInvokeOptions{}

	cmd := &cobra.Command{
		Use:    "invoke",
		Hidden: true,
		Short:  "invoke an http request handler using curl",
		Long: strings.TrimSpace(`
This command is not supported and may be removed in the future.
`),
		Example: strings.Join([]string{
			fmt.Sprintf("%s handler invoke my-handler", c.Name),
			fmt.Sprintf("%s handler invoke my-handler --text -- -d 'hello' -w '\\n'", c.Name),
			fmt.Sprintf("%s handler invoke my-handler /request/path", c.Name),
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

func (opts *HandlerInvokeOptions) findAvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return port, err
}

func (opts *HandlerInvokeOptions) findReadyPod(ctx context.Context, c *cli.Config) (*corev1.Pod, error) {
	pods, err := c.Core().Pods(opts.Namespace).List(metav1.ListOptions{
		FieldSelector: "status.phase=Running",
		LabelSelector: fmt.Sprintf("%s=%s", request.HandlerLabelKey, opts.Name),
	})
	if err != nil {
		return nil, err
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no running pods found")
	}
	return &pods.Items[0], nil
}

func (opts *HandlerInvokeOptions) forwardPort(ctx context.Context, c *cli.Config, pod *corev1.Pod, port int, stopChan chan struct{}) error {
	url := c.Core().RESTClient().
		Post().
		Resource("pods").
		Namespace(pod.Namespace).
		Name(pod.Name).
		SubResource("portforward").
		URL()

	transport, upgrader, err := spdy.RoundTripperFor(c.KubeRestConfig())
	if err != nil {
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", url)

	ports := []string{
		fmt.Sprintf("%d:%d", port, 8080),
	}
	readyChan := make(chan struct{}, 1)
	out := &bytes.Buffer{}
	fw, err := portforward.New(dialer, ports, stopChan, readyChan, out, out)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		errChan <- fw.ForwardPorts()
	}()

	select {
	case err = <-errChan:
		c.Stderr.Write(out.Bytes())
		return fmt.Errorf("forwarding ports: %v", err)
	case <-fw.Ready:
		return nil
	}
}

func (opts *HandlerInvokeOptions) invoke(ctx context.Context, c *cli.Config, port int) error {
	curlArgs := []string{fmt.Sprintf("http://localhost:%d", port)}
	if opts.ContentTypeJSON {
		curlArgs = append(curlArgs, "-H", "Content-Type: application/json")
	}
	if opts.ContentTypeText {
		curlArgs = append(curlArgs, "-H", "Content-Type: text/plain")
	}
	curlArgs = append(curlArgs, opts.BareArgs...)

	curl := c.Exec(ctx, "curl", curlArgs...)

	curl.Stdin = c.Stdin
	curl.Stdout = c.Stdout
	curl.Stderr = c.Stderr

	return curl.Run()
}
