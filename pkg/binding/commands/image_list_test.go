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
	"testing"

	"github.com/projectriff/cli/pkg/binding/commands"
	"github.com/projectriff/cli/pkg/cli"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	"github.com/projectriff/reconciler-runtime/apis"
	bindingsv1alpha1 "github.com/projectriff/system/pkg/apis/bindings/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestImageBindingListOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid list",
			Options: &commands.ImageListOptions{
				ListOptions: rifftesting.InvalidListOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidListOptionsFieldError,
		},
		{
			Name: "valid list",
			Options: &commands.ImageListOptions{
				ListOptions: rifftesting.ValidListOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestImageBindingListCommand(t *testing.T) {
	imageBindingName := "test-image-binding"
	imageBindingOtherName := "test-other-image-binding"
	defaultNamespace := "default"
	otherNamespace := "other-namespace"

	table := rifftesting.CommandTable{
		{
			Name: "invalid args",
			Args: []string{},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				// disable default namespace
				c.Client.(*rifftesting.FakeClient).Namespace = ""
				return ctx, nil
			},
			ShouldError: true,
		},
		{
			Name: "empty",
			Args: []string{},
			ExpectOutput: `
No image bindings found.
`,
		},
		{
			Name: "lists an item",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
NAME                 SUBJECT   PROVIDER   CONTAINER NAME   STATUS      AGE
test-image-binding   <empty>   <empty>    <empty>          <unknown>   <unknown>
`,
		},
		{
			Name: "filters by namespace",
			Args: []string{cli.NamespaceFlagName, otherNamespace},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
No image bindings found.
`,
		},
		{
			Name: "all namespace",
			Args: []string{cli.AllNamespacesFlagName},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
				},
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingOtherName,
						Namespace: otherNamespace,
					},
				},
			},
			ExpectOutput: `
NAMESPACE         NAME                       SUBJECT   PROVIDER   CONTAINER NAME   STATUS      AGE
default           test-image-binding         <empty>   <empty>    <empty>          <unknown>   <unknown>
other-namespace   test-other-image-binding   <empty>   <empty>    <empty>          <unknown>   <unknown>
`,
		},
		{
			Name: "table populates all columns",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
					Spec: bindingsv1alpha1.ImageBindingSpec{
						Subject: &bindingsv1alpha1.Reference{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Namespace:  "default",
							Name:       "my-deployment",
						},
						Provider: &bindingsv1alpha1.Reference{
							APIVersion: "build.projectriff.io/v1alpha1",
							Kind:       "Function",
							Namespace:  "default",
							Name:       "my-function",
						},
						ContainerName: "user-container",
					},
					Status: bindingsv1alpha1.ImageBindingStatus{
						Status: apis.Status{
							Conditions: apis.Conditions{
								{Type: apis.ConditionReady, Status: "True"},
							},
						},
					},
				},
			},
			ExpectOutput: `
NAME                 SUBJECT                          PROVIDER                                     CONTAINER NAME   STATUS   AGE
test-image-binding   deployments.apps:my-deployment   functions.build.projectriff.io:my-function   user-container   Ready    <unknown>
`,
		},
		{
			Name: "list error",
			Args: []string{},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("list", "imagebindings"),
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewImageListCommand)
}
