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

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/k8s"
	"github.com/projectriff/cli/pkg/streaming/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	kailtesting "github.com/projectriff/cli/pkg/testing/kail"
	streamv1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	cachetesting "k8s.io/client-go/tools/cache/testing"
)

func TestProcessorCreateOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidResourceOptionsFieldError.Also(
				cli.ErrMissingField(cli.FunctionRefFlagName),
				cli.ErrMissingField(cli.InputFlagName),
			),
		},
		{
			Name: "with inputs but no outputs",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
				Inputs:          []string{"input1", "input2"},
			},
			ShouldValidate: true,
		},
		{
			Name: "with inputs and outputs",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
				Inputs:          []string{"input1", "input2"},
				Outputs:         []string{"output1", "output2"},
			},
			ShouldValidate: true,
		},
		{
			Name: "with tail",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
				Inputs:          []string{"input"},
				Tail:            true,
				WaitTimeout:     "10m",
			},
			ShouldValidate: true,
		},
		{
			Name: "with tail, missing timeout",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
				Inputs:          []string{"input"},
				Tail:            true,
			},
			ExpectFieldErrors: cli.ErrMissingField(cli.WaitTimeoutFlagName),
		},
		{
			Name: "with tail, invalid timeout",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
				Inputs:          []string{"input"},
				Tail:            true,
				WaitTimeout:     "d",
			},
			ExpectFieldErrors: cli.ErrInvalidValue("d", cli.WaitTimeoutFlagName),
		},
		{
			Name: "dry run",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
				Inputs:          []string{"input"},
				DryRun:          true,
			},
			ShouldValidate: true,
		},
		{
			Name: "dry run, tail",
			Options: &commands.ProcessorCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				FunctionRef:     "my-function",
				Inputs:          []string{"input"},
				Tail:            true,
				WaitTimeout:     "10m",
				DryRun:          true,
			},
			ExpectFieldErrors: cli.ErrMultipleOneOf(cli.DryRunFlagName, cli.TailFlagName),
		},
	}

	table.Run(t)
}

func TestProcessorCreateCommand(t *testing.T) {
	defaultNamespace := "default"
	processorName := "my-processor"
	functionRef := "my-func"
	inputName := "input"
	outputName := "output"
	inputNameOther := "otherinput"
	outputNameOther := "otheroutput"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "create with single input",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
					},
				},
			},
			ExpectOutput: `
Created processor "my-processor"
`,
		},
		{
			Name: "dry run",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName, cli.DryRunFlagName},
			ExpectOutput: `
---
apiVersion: streaming.projectriff.io/v1alpha1
kind: Processor
metadata:
  creationTimestamp: null
  name: my-processor
  namespace: default
spec:
  functionRef: my-func
  inputs:
  - input
  outputs: []
status: {}

Created processor "my-processor"
`,
		},
		{
			Name: "create with multiple inputs",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName, cli.InputFlagName, inputNameOther},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName, inputNameOther},
					},
				},
			},
			ExpectOutput: `
Created processor "my-processor"
`,
		},
		{
			Name: "create with single output",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName, cli.InputFlagName, inputNameOther, cli.OutputFlagName, outputName},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName, inputNameOther},
						Outputs:     []string{outputName},
					},
				},
			},
			ExpectOutput: `
Created processor "my-processor"
`,
		},
		{
			Name: "create with multiple outputs",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName, cli.InputFlagName, inputNameOther, cli.OutputFlagName, outputName, cli.OutputFlagName, outputNameOther},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName, inputNameOther},
						Outputs:     []string{outputName, outputNameOther},
					},
				},
			},
			ExpectOutput: `
Created processor "my-processor"
`,
		},
		{
			Name: "error existing processor",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error during create",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("create", "processors"),
			},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "tail logs",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName, cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("StreamingProcessorLogs", mock.Anything, &streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
						Outputs:     []string{},
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
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
						Outputs:     []string{},
					},
				},
			},
			ExpectOutput: `
Created processor "my-processor"
...log output...
Processor "my-processor" is ready
`,
		},
		{
			Name: "tail timeout",
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName, cli.TailFlagName, cli.WaitTimeoutFlagName, "5ms"},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("StreamingProcessorLogs", mock.Anything, &streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
						Outputs:     []string{},
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
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
						Outputs:     []string{},
					},
				},
			},
			ExpectOutput: `
Created processor "my-processor"
...log output...
Timeout after "5ms" waiting for "my-processor" to become ready
To view status run: riff processor list --namespace default
To continue watching logs run: riff processor tail my-processor --namespace default
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
			Args: []string{processorName, cli.FunctionRefFlagName, functionRef, cli.InputFlagName, inputName, cli.TailFlagName},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				lw := cachetesting.NewFakeControllerSource()
				ctx = k8s.WithListerWatcher(ctx, lw)

				kail := &kailtesting.Logger{}
				c.Kail = kail
				kail.On("StreamingProcessorLogs", mock.Anything, &streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
						Outputs:     []string{},
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
				&streamv1alpha1.Processor{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      processorName,
					},
					Spec: streamv1alpha1.ProcessorSpec{
						FunctionRef: functionRef,
						Inputs:      []string{inputName},
						Outputs:     []string{},
					},
				},
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewProcessorCreateCommand)
}
