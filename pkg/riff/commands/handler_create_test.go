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
	"context"
	"fmt"
	"testing"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/k8s"
	"github.com/projectriff/cli/pkg/riff/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	kailtesting "github.com/projectriff/cli/pkg/testing/kail"
	requestv1alpha1 "github.com/projectriff/system/pkg/apis/request/v1alpha1"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	cachetesting "k8s.io/client-go/tools/cache/testing"
)

func TestHandlerCreateOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldError: rifftesting.InvalidResourceOptionsFieldError.Also(
				cli.ErrMissingOneOf(cli.ApplicationRefFlagName, cli.ContainerRefFlagName, cli.FunctionRefFlagName, cli.ImageFlagName),
			),
		},
		{
			Name: "from application",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ApplicationRef:  "my-application",
			},
			ShouldValidate: true,
		},
		{
			Name: "from container",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ContainerRef:    "my-container",
			},
			ShouldValidate: true,
		},
		{
			Name: "from function",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
			},
			ShouldValidate: true,
		},
		{
			Name: "from image",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
			},
			ShouldValidate: true,
		},
		{
			Name: "from application, container, funcation and image",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ApplicationRef:  "my-application",
				ContainerRef:    "my-container",
				FunctionRef:     "my-function",
				Image:           "example.com/repo:tag",
			},
			ExpectFieldError: cli.ErrMultipleOneOf(cli.ApplicationRefFlagName, cli.ContainerRefFlagName, cli.FunctionRefFlagName, cli.ImageFlagName),
		},
		{
			Name: "with env",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				Env:             []string{"VAR1=foo", "VAR2=bar"},
			},
			ShouldValidate: true,
		},
		{
			Name: "with invalid env",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				Env:             []string{"=foo"},
			},
			ExpectFieldError: cli.ErrInvalidArrayValue("=foo", cli.EnvFlagName, 0),
		},
		{
			Name: "with envfrom secret",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				EnvFrom:         []string{"VAR1=secretKeyRef:name:key"},
			},
			ShouldValidate: true,
		},
		{
			Name: "with envfrom configmap",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				EnvFrom:         []string{"VAR1=configMapKeyRef:name:key"},
			},
			ShouldValidate: true,
		},
		{
			Name: "with invalid envfrom",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				EnvFrom:         []string{"VAR1=someOtherKeyRef:name:key"},
			},
			ExpectFieldError: cli.ErrInvalidArrayValue("VAR1=someOtherKeyRef:name:key", cli.EnvFromFlagName, 0),
		},
		{
			Name: "with tail",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				Tail:            true,
				WaitTimeout:     "10m",
			},
			ShouldValidate: true,
		},
		{
			Name: "with tail, missing timeout",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				Tail:            true,
			},
			ExpectFieldError: cli.ErrMissingField(cli.WaitTimeoutFlagName),
		},
		{
			Name: "with tail, invalid timeout",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				Tail:            true,
				WaitTimeout:     "d",
			},
			ExpectFieldError: cli.ErrInvalidValue("d", cli.WaitTimeoutFlagName),
		},
		{
			Name: "dry run",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				DryRun:          true,
			},
			ShouldValidate: true,
		},
		{
			Name: "dry run, tail",
			Options: &commands.HandlerCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				Tail:            true,
				WaitTimeout:     "10m",
				DryRun:          true,
			},
			ExpectFieldError: cli.ErrMultipleOneOf(cli.DryRunFlagName, cli.TailFlagName),
		},
	}

	table.Run(t)
}

func TestHandlerCreateCommand(t *testing.T) {
	defaultNamespace := "default"
	handlerName := "my-handler"
	image := "registry.example.com/repo@sha256:deadbeefdeadbeefdeadbeefdeadbeef"
	applicationRef := "my-app"
	containerRef := "my-container"
	functionRef := "my-func"
	envName := "MY_VAR"
	envValue := "my-value"
	envVar := fmt.Sprintf("%s=%s", envName, envValue)
	envNameOther := "MY_VAR_OTHER"
	envValueOther := "my-value-other"
	envVarOther := fmt.Sprintf("%s=%s", envNameOther, envValueOther)
	envVarFromConfigMap := "MY_VAR_FROM_CONFIGMAP=configMapKeyRef:my-configmap:my-key"
	envVarFromSecret := "MY_VAR_FROM_SECRET=secretKeyRef:my-secret:my-key"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "create from image",
			Args: []string{handlerName, cli.ImageFlagName, image},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{
								{Image: image},
							},
						},
					},
				},
			},
			ExpectOutput: `
Created handler "my-handler"
`,
		},
		{
			Name: "create from application ref",
			Args: []string{handlerName, cli.ApplicationRefFlagName, applicationRef},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Build: &requestv1alpha1.Build{
							ApplicationRef: applicationRef,
						},
					},
				},
			},
			ExpectOutput: `
Created handler "my-handler"
`,
		},
		{
			Name: "create from container ref",
			Args: []string{handlerName, cli.ContainerRefFlagName, containerRef},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Build: &requestv1alpha1.Build{
							ContainerRef: containerRef,
						},
					},
				},
			},
			ExpectOutput: `
Created handler "my-handler"
`,
		},
		{
			Name: "create from function ref",
			Args: []string{handlerName, cli.FunctionRefFlagName, functionRef},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Build: &requestv1alpha1.Build{
							FunctionRef: functionRef,
						},
					},
				},
			},
			ExpectOutput: `
Created handler "my-handler"
`,
		},
		{
			Name: "dry run",
			Args: []string{handlerName, cli.ImageFlagName, image, cli.DryRunFlagName},
			ExpectOutput: `
---
apiVersion: request.projectriff.io/v1alpha1
kind: Handler
metadata:
  creationTimestamp: null
  name: my-handler
  namespace: default
spec:
  template:
    containers:
    - image: registry.example.com/repo@sha256:deadbeefdeadbeefdeadbeefdeadbeef
      name: ""
      resources: {}
status: {}

Created handler "my-handler"
`,
		},
		{
			Name: "create from image with env and env-from",
			Args: []string{handlerName, cli.ImageFlagName, image, cli.EnvFlagName, envVar, cli.EnvFlagName, envVarOther, cli.EnvFromFlagName, envVarFromConfigMap, cli.EnvFromFlagName, envVarFromSecret},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Image: image,
									Env: []corev1.EnvVar{
										{Name: envName, Value: envValue},
										{Name: envNameOther, Value: envValueOther},
										{
											Name: "MY_VAR_FROM_CONFIGMAP",
											ValueFrom: &corev1.EnvVarSource{
												ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
													LocalObjectReference: corev1.LocalObjectReference{
														Name: "my-configmap",
													},
													Key: "my-key",
												},
											},
										},
										{
											Name: "MY_VAR_FROM_SECRET",
											ValueFrom: &corev1.EnvVarSource{
												SecretKeyRef: &corev1.SecretKeySelector{
													LocalObjectReference: corev1.LocalObjectReference{
														Name: "my-secret",
													},
													Key: "my-key",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			ExpectOutput: `
Created handler "my-handler"
`,
		},
		{
			Name: "error existing handler",
			Args: []string{handlerName, cli.ImageFlagName, image},
			GivenObjects: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{
								{Image: image},
							},
						},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error during create",
			Args: []string{handlerName, cli.ImageFlagName, image},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("create", "handlers"),
			},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{
								{Image: image},
							},
						},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "tail logs",
			Args: []string{handlerName, cli.ImageFlagName, image, cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("HandlerLogs", mock.Anything, &requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{{Image: image}},
						},
					},
				}, cli.TailSinceCreateDefault, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...log output...\n")
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				if lw, ok := k8s.GetListerWatcher(ctx, nil, "", nil).(*cachetesting.FakeControllerSource); ok {
					lw.Shutdown()
				}

				kail := c.Kail.(*kailtesting.Logger)
				kail.AssertExpectations(t)
				return nil
			},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{{Image: image}},
						},
					},
				},
			},
			ExpectOutput: `
Created handler "my-handler"
...log output...
`,
		},
		{
			Name: "tail timeout",
			Args: []string{handlerName, cli.ImageFlagName, image, cli.TailFlagName, cli.WaitTimeoutFlagName, "5ms"},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("HandlerLogs", mock.Anything, &requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{{Image: image}},
						},
					},
				}, cli.TailSinceCreateDefault, mock.Anything).Return(k8s.ErrWaitTimeout).Run(func(args mock.Arguments) {
					ctx := args[0].(context.Context)
					fmt.Fprintf(c.Stdout, "...log output...\n")
					// wait for context to be cancelled
					<-ctx.Done()
				})
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				if lw, ok := k8s.GetListerWatcher(ctx, nil, "", nil).(*cachetesting.FakeControllerSource); ok {
					lw.Shutdown()
				}

				kail := c.Kail.(*kailtesting.Logger)
				kail.AssertExpectations(t)
				return nil
			},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{{Image: image}},
						},
					},
				},
			},
			ExpectOutput: `
Created handler "my-handler"
...log output...
Timeout after "5ms" waiting for "my-handler" to become ready
To view status run: riff handler list --namespace default
To continue watching logs run: riff handler tail my-handler --namespace default
`,
			ShouldError: true,
			Verify: func(t *testing.T, output string, err error) {
				if actual := err; !cli.IsSilent(err) {
					t.Errorf("expected error to be silent, actual %#v", actual)
				}
			},
		},
		{
			Name: "tail error",
			Args: []string{handlerName, cli.ImageFlagName, image, cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("HandlerLogs", mock.Anything, &requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{{Image: image}},
						},
					},
				}, cli.TailSinceCreateDefault, mock.Anything).Return(fmt.Errorf("kail error"))
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, c *cli.Config) error {
				if lw, ok := k8s.GetListerWatcher(ctx, nil, "", nil).(*cachetesting.FakeControllerSource); ok {
					lw.Shutdown()
				}

				kail := c.Kail.(*kailtesting.Logger)
				kail.AssertExpectations(t)
				return nil
			},
			ExpectCreates: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      handlerName,
					},
					Spec: requestv1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{{Image: image}},
						},
					},
				},
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewHandlerCreateCommand)
}
