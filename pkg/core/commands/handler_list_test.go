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

	duckv1beta1 "github.com/knative/pkg/apis/duck/v1beta1"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/core/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	corev1alpha1 "github.com/projectriff/system/pkg/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestHandlerListOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid list",
			Options: &commands.HandlerListOptions{
				ListOptions: rifftesting.InvalidListOptions,
			},
			ExpectFieldError: rifftesting.InvalidListOptionsFieldError,
		},
		{
			Name: "valid list",
			Options: &commands.HandlerListOptions{
				ListOptions: rifftesting.ValidListOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestHandlerListCommand(t *testing.T) {
	handlerName := "test-handler"
	handlerOtherName := "test-other-handler"
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
No handlers found.
`,
		},
		{
			Name: "lists an item",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      handlerName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
NAME           TYPE        REF         SERVICE   STATUS      AGE
test-handler   <unknown>   <unknown>   <empty>   <unknown>   <unknown>
`,
		},
		{
			Name: "filters by namespace",
			Args: []string{cli.NamespaceFlagName, otherNamespace},
			GivenObjects: []runtime.Object{
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      handlerName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
No handlers found.
`,
		},
		{
			Name: "all namespace",
			Args: []string{cli.AllNamespacesFlagName},
			GivenObjects: []runtime.Object{
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      handlerName,
						Namespace: defaultNamespace,
					},
				},
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      handlerOtherName,
						Namespace: otherNamespace,
					},
				},
			},
			ExpectOutput: `
NAMESPACE         NAME                 TYPE        REF         SERVICE   STATUS      AGE
default           test-handler         <unknown>   <unknown>   <empty>   <unknown>   <unknown>
other-namespace   test-other-handler   <unknown>   <unknown>   <empty>   <unknown>   <unknown>
`,
		},
		{
			Name: "table populates all columns",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "img",
						Namespace: defaultNamespace,
					},
					Spec: corev1alpha1.HandlerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{
								{Image: "projectriff/upper"},
							},
						},
					},
					Status: corev1alpha1.HandlerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: corev1alpha1.HandlerConditionReady, Status: "True"},
							},
						},
						DeploymentName: "img-handler",
						ServiceName:    "img-handler",
					},
				},
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "app",
						Namespace: defaultNamespace,
					},
					Spec: corev1alpha1.HandlerSpec{
						Build: &corev1alpha1.Build{ApplicationRef: "petclinic"},
					},
					Status: corev1alpha1.HandlerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: corev1alpha1.HandlerConditionReady, Status: "True"},
							},
						},
						DeploymentName: "app-handler",
						ServiceName:    "app-handler",
					},
				},
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "func",
						Namespace: defaultNamespace,
					},
					Spec: corev1alpha1.HandlerSpec{
						Build: &corev1alpha1.Build{FunctionRef: "square"},
					},
					Status: corev1alpha1.HandlerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: corev1alpha1.HandlerConditionReady, Status: "True"},
							},
						},
						DeploymentName: "func-handler",
						ServiceName:    "func-handler",
					},
				},
				&corev1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "container",
						Namespace: defaultNamespace,
					},
					Spec: corev1alpha1.HandlerSpec{
						Build: &corev1alpha1.Build{ContainerRef: "busybox"},
					},
					Status: corev1alpha1.HandlerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: corev1alpha1.HandlerConditionReady, Status: "True"},
							},
						},
						DeploymentName: "container-handler",
						ServiceName:    "container-handler",
					},
				},
			},
			ExpectOutput: `
NAME        TYPE          REF                 SERVICE             STATUS   AGE
app         application   petclinic           app-handler         Ready    <unknown>
container   container     busybox             container-handler   Ready    <unknown>
func        function      square              func-handler        Ready    <unknown>
img         image         projectriff/upper   img-handler         Ready    <unknown>
`,
		},
		{
			Name: "list error",
			Args: []string{},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("list", "handlers"),
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewHandlerListCommand)
}
