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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/projectriff/cli/pkg/k8s"
	kailtesting "github.com/projectriff/cli/pkg/testing/kail"
	"github.com/stretchr/testify/mock"
	cachetesting "k8s.io/client-go/tools/cache/testing"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/streaming/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	streamv1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestPulsarProviderCreateOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.PulsarProviderCreateOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidResourceOptionsFieldError.Also(
				cli.ErrMissingField(cli.ServiceURLFlagName),
			),
		},
		{
			Name: "valid provider",
			Options: &commands.PulsarProviderCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ServiceURL:      "pulsar://localhost:6650",
			},
			ShouldValidate: true,
		},
		{
			Name: "dry run",
			Options: &commands.PulsarProviderCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ServiceURL:      "pulsar://localhost:6650",
				DryRun:          true,
			},
			ShouldValidate: true,
		},
		{
			Name: "dry run, tail",
			Options: &commands.PulsarProviderCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ServiceURL:      "pulsar://localhost:6650",
				DryRun:          true,
				Tail:            true,
			},
			ExpectFieldErrors: cli.ErrMultipleOneOf(cli.DryRunFlagName, cli.TailFlagName),
		},
		{
			Name: "invalid timeout",
			Options: &commands.PulsarProviderCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ServiceURL:      "pulsar://localhost:6650",
				WaitTimeout:     -4 * time.Second,
			},
			ExpectFieldErrors: cli.ErrInvalidValue(-4*time.Second, cli.WaitTimeoutFlagName),
		},
	}

	table.Run(t)
}

func TestPulsarProviderCreateCommand(t *testing.T) {
	defaultNamespace := "default"
	pulsarProviderName := "my-pulsar-provider"
	serviceURL := "pulsar://localhost:6650"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "pulsar provider",
			Args: []string{pulsarProviderName, cli.ServiceURLFlagName, serviceURL},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      pulsarProviderName,
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: serviceURL,
					},
				},
			},
			ExpectOutput: `
Created pulsar provider "my-pulsar-provider"
`,
		},
		{
			Name: "dry run",
			Args: []string{pulsarProviderName, cli.ServiceURLFlagName, serviceURL, cli.DryRunFlagName},
			ExpectOutput: `
---
apiVersion: streaming.projectriff.io/v1alpha1
kind: PulsarProvider
metadata:
  creationTimestamp: null
  name: my-pulsar-provider
  namespace: default
spec:
  serviceURL: pulsar://localhost:6650
status: {}

Created pulsar provider "my-pulsar-provider"
`,
		},
		{
			Name: "error existing provider",
			Args: []string{pulsarProviderName, cli.ServiceURLFlagName, serviceURL},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      pulsarProviderName,
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      pulsarProviderName,
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: serviceURL,
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error during create",
			Args: []string{pulsarProviderName, cli.ServiceURLFlagName, serviceURL},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("create", "pulsarproviders"),
			},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      pulsarProviderName,
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: serviceURL,
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "tail logs",
			Args: []string{"franz", cli.ServiceURLFlagName, "some-host", cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("PulsarProviderLogs", mock.Anything, &streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "franz",
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: "some-host",
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
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "franz",
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: "some-host",
					},
				},
			},
			ExpectOutput: `
Created pulsar provider "franz"
...log output...
PulsarProvider "franz" is ready
`,
		},
		{
			Name: "tail timeout",
			Args: []string{"franz", cli.ServiceURLFlagName, "some-host", cli.TailFlagName, cli.WaitTimeoutFlagName, "7ms"},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("PulsarProviderLogs", mock.Anything, &streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "franz",
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: "some-host",
					},
				}, cli.TailSinceCreateDefault, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprintf(c.Stdout, "...log output...\n")
					// wait for context to be cancelled
					<-args[0].(context.Context).Done()
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
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "franz",
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: "some-host",
					},
				},
			},
			ExpectOutput: `
Created pulsar provider "franz"
...log output...
Timeout after "7ms" waiting for "franz" to become ready
To view status run: riff streaming pulsar-provider list --namespace default
`,
			ShouldError: true,
			Verify: func(t *testing.T, output string, err error) {
				if actual := err; !errors.Is(err, cli.SilentError) {
					t.Errorf("expected error to be silent, actual %#v", actual)
				}
			},
		},
		{
			Name: "tail error",
			Args: []string{"franz", cli.ServiceURLFlagName, "some-host", cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("PulsarProviderLogs", mock.Anything, &streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "franz",
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: "some-host",
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
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      "franz",
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: "some-host",
					},
				},
			},
			ShouldError: true,
			ExpectOutput: `
Created pulsar provider "franz"
`,
		},
	}

	table.Run(t, commands.NewPulsarProviderCreateCommand)
}
