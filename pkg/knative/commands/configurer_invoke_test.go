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

package commands_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/knative/pkg/apis"
	duckv1beta1 "github.com/knative/pkg/apis/duck/v1beta1"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/knative/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	knativev1alpha1 "github.com/projectriff/system/pkg/apis/knative/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestConfigurerInvokeOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.ConfigurerInvokeOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldError: rifftesting.InvalidResourceOptionsFieldError,
		},
		{
			Name: "valid resource",
			Options: &commands.ConfigurerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
			},
			ShouldValidate: true,
		},
		{
			Name: "json content type",
			Options: &commands.ConfigurerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ContentTypeJSON: true,
			},
			ShouldValidate: true,
		},
		{
			Name: "text content type",
			Options: &commands.ConfigurerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ContentTypeText: true,
			},
			ShouldValidate: true,
		},
		{
			Name: "multiple content types",
			Options: &commands.ConfigurerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ContentTypeJSON: true,
				ContentTypeText: true,
			},
			ExpectFieldError: cli.ErrMultipleOneOf(cli.JSONFlagName, cli.TextFlagName),
		},
	}

	table.Run(t)
}

func TestConfigurerInvokeCommand(t *testing.T) {
	t.Parallel()

	configurerName := "test-configurer"
	defaultNamespace := "default"

	configurer := &knativev1alpha1.Configurer{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: defaultNamespace,
			Name:      configurerName,
		},
		Status: knativev1alpha1.ConfigurerStatus{
			Status: duckv1beta1.Status{
				Conditions: duckv1beta1.Conditions{
					{Type: knativev1alpha1.ConfigurerConditionReady, Status: "True"},
				},
			},
			URL: &apis.URL{
				Host: fmt.Sprintf("%s.example.com", configurerName),
			},
		},
	}

	ingressService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "istio-system",
			Name:      "istio-ingressgateway",
		},
		Spec: corev1.ServiceSpec{
			Type: "LoadBalancer",
		},
		Status: corev1.ServiceStatus{
			LoadBalancer: corev1.LoadBalancerStatus{
				Ingress: []corev1.LoadBalancerIngress{
					{Hostname: "localhost"},
				},
			},
		},
	}

	table := rifftesting.CommandTable{
		{
			Name:       "ingress loadbalancer hostname",
			Args:       []string{configurerName},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "istio-system",
						Name:      "istio-ingressgateway",
					},
					Spec: corev1.ServiceSpec{
						Type: "LoadBalancer",
					},
					Status: corev1.ServiceStatus{
						LoadBalancer: corev1.LoadBalancerStatus{
							Ingress: []corev1.LoadBalancerIngress{
								{Hostname: "localhost"},
							},
						},
					},
				},
			},
			Verify: func(t *testing.T, output string, err error) {
				for _, expected := range []string{
					"curl localhost -H 'Host: test-configurer.example.com'\n",
				} {
					if !strings.Contains(output, expected) {
						t.Errorf("expected command output to contain %q, actually %q", expected, output)
					}
				}
			},
		},
		{
			Name:       "ingress loadbalancer ip",
			Args:       []string{configurerName},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "istio-system",
						Name:      "istio-ingressgateway",
					},
					Spec: corev1.ServiceSpec{
						Type: "LoadBalancer",
					},
					Status: corev1.ServiceStatus{
						LoadBalancer: corev1.LoadBalancerStatus{
							Ingress: []corev1.LoadBalancerIngress{
								{IP: "127.0.0.1"},
							},
						},
					},
				},
			},
			Verify: func(t *testing.T, output string, err error) {
				for _, expected := range []string{
					"curl 127.0.0.1 -H 'Host: test-configurer.example.com'\n",
				} {
					if !strings.Contains(output, expected) {
						t.Errorf("expected command output to contain %q, actually %q", expected, output)
					}
				}
			},
		},
		{
			Name:       "ingress nodeport",
			Args:       []string{configurerName},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "istio-system",
						Name:      "istio-ingressgateway",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{Name: "http2", NodePort: 54321},
						},
					},
				},
			},
			Verify: func(t *testing.T, output string, err error) {
				for _, expected := range []string{
					"curl http://localhost:54321 -H 'Host: test-configurer.example.com'\n",
				} {
					if !strings.Contains(output, expected) {
						t.Errorf("expected command output to contain %q, actually %q", expected, output)
					}
				}
			},
		},
		{
			Name:       "request path",
			Args:       []string{configurerName, "/path"},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
				ingressService,
			},
			Verify: func(t *testing.T, output string, err error) {
				for _, expected := range []string{
					"curl localhost/path -H 'Host: test-configurer.example.com'\n",
				} {
					if !strings.Contains(output, expected) {
						t.Errorf("expected command output to contain %q, actually %q", expected, output)
					}
				}
			},
		},
		{
			Name:       "content type json",
			Args:       []string{configurerName, cli.JSONFlagName},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
				ingressService,
			},
			Verify: func(t *testing.T, output string, err error) {
				for _, expected := range []string{
					"curl localhost -H 'Host: test-configurer.example.com' -H 'Content-Type: application/json'\n",
				} {
					if !strings.Contains(output, expected) {
						t.Errorf("expected command output to contain %q, actually %q", expected, output)
					}
				}
			},
		},
		{
			Name:       "content type text",
			Args:       []string{configurerName, cli.TextFlagName},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
				ingressService,
			},
			Verify: func(t *testing.T, output string, err error) {
				for _, expected := range []string{
					"curl localhost -H 'Host: test-configurer.example.com' -H 'Content-Type: text/plain'\n",
				} {
					if !strings.Contains(output, expected) {
						t.Errorf("expected command output to contain %q, actually %q", expected, output)
					}
				}
			},
		},
		{
			Name:       "pass extra args to curl",
			Args:       []string{configurerName, "--", "-w", "\n"},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
				ingressService,
			},
			ExpectOutput: `
Command executed: curl localhost -H 'Host: test-configurer.example.com' -w '` + "\n" + `'
`,
		},
		{
			Name: "unknown ingress",
			Args: []string{configurerName},
			GivenObjects: []runtime.Object{
				configurer,
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "istio-system",
						Name:      "istio-ingressgateway",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name:       "missing ingress",
			Args:       []string{configurerName},
			ExecHelper: "ConfigurerInvoke",
			GivenObjects: []runtime.Object{
				configurer,
			},
			ShouldError: true,
		},
		{
			Name: "configurer not ready",
			Args: []string{configurerName},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      configurerName,
					},
					Status: knativev1alpha1.ConfigurerStatus{
						URL: &apis.URL{
							Host: fmt.Sprintf("%s.example.com", configurerName),
						},
					},
				},
				ingressService,
			},
			ShouldError: true,
		},
		{
			Name: "configurer missing domain",
			Args: []string{configurerName},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      configurerName,
					},
					Status: knativev1alpha1.ConfigurerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: knativev1alpha1.ConfigurerConditionReady, Status: "True"},
							},
						},
					},
				},
				ingressService,
			},
			ShouldError: true,
		},
		{
			Name: "missing configurer",
			Args: []string{configurerName},
			GivenObjects: []runtime.Object{
				ingressService,
			},
			ShouldError: true,
		},
		{
			Name:       "curl error",
			Args:       []string{configurerName},
			ExecHelper: "ConfigurerInvokeError",
			GivenObjects: []runtime.Object{
				configurer,
				ingressService,
			},
			ExpectOutput: `
Command executed: curl localhost -H 'Host: test-configurer.example.com'
`,
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewConfigurerInvokeCommand)
}

func TestHelperProcess_ConfigurerInvoke(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stderr, "Command executed: %s\n", shellquote.Join(argsAfterBareDoubleDash(os.Args)...))
	os.Exit(0)
}

func TestHelperProcess_ConfigurerInvokeError(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stderr, "Command executed: %s\n", shellquote.Join(argsAfterBareDoubleDash(os.Args)...))
	os.Exit(1)
}

func argsAfterBareDoubleDash(args []string) []string {
	for i, arg := range args {
		if arg == "--" {
			return args[i+1:]
		}
	}
	return []string{}
}
