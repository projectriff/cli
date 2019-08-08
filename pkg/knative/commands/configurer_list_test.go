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

func TestConfigurerListOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid list",
			Options: &commands.ConfigurerListOptions{
				ListOptions: rifftesting.InvalidListOptions,
			},
			ExpectFieldError: rifftesting.InvalidListOptionsFieldError,
		},
		{
			Name: "valid list",
			Options: &commands.ConfigurerListOptions{
				ListOptions: rifftesting.ValidListOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestConfigurerListCommand(t *testing.T) {
	configurerName := "test-configurer"
	configurerOtherName := "test-other-configurer"
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
No configurers found.
`,
		},
		{
			Name: "lists an item",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
NAME              TYPE        REF         HOST      STATUS      AGE
test-configurer   <unknown>   <unknown>   <empty>   <unknown>   <unknown>
`,
		},
		{
			Name: "filters by namespace",
			Args: []string{cli.NamespaceFlagName, otherNamespace},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
No configurers found.
`,
		},
		{
			Name: "all namespace",
			Args: []string{cli.AllNamespacesFlagName},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerName,
						Namespace: defaultNamespace,
					},
				},
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerOtherName,
						Namespace: otherNamespace,
					},
				},
			},
			ExpectOutput: `
NAMESPACE         NAME                    TYPE        REF         HOST      STATUS      AGE
default           test-configurer         <unknown>   <unknown>   <empty>   <unknown>   <unknown>
other-namespace   test-other-configurer   <unknown>   <unknown>   <empty>   <unknown>   <unknown>
`,
		},
		{
			Name: "table populates all columns",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "img",
						Namespace: defaultNamespace,
					},
					Spec: knativev1alpha1.ConfigurerSpec{
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{
								{Image: "projectriff/upper"},
							},
						},
					},
					Status: knativev1alpha1.ConfigurerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: knativev1alpha1.ConfigurerConditionReady, Status: "True"},
							},
						},
						URL: &apis.URL{
							Host: "img.default.example.com",
						},
					},
				},
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "app",
						Namespace: defaultNamespace,
					},
					Spec: knativev1alpha1.ConfigurerSpec{
						Build: &knativev1alpha1.Build{ApplicationRef: "petclinic"},
					},
					Status: knativev1alpha1.ConfigurerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: knativev1alpha1.ConfigurerConditionReady, Status: "True"},
							},
						},
						URL: &apis.URL{
							Host: "app.default.example.com",
						},
					},
				},
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "func",
						Namespace: defaultNamespace,
					},
					Spec: knativev1alpha1.ConfigurerSpec{
						Build: &knativev1alpha1.Build{FunctionRef: "square"},
					},
					Status: knativev1alpha1.ConfigurerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: knativev1alpha1.ConfigurerConditionReady, Status: "True"},
							},
						},
						URL: &apis.URL{
							Host: "func.default.example.com",
						},
					},
				},
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "container",
						Namespace: defaultNamespace,
					},
					Spec: knativev1alpha1.ConfigurerSpec{
						Build: &knativev1alpha1.Build{ContainerRef: "busybox"},
					},
					Status: knativev1alpha1.ConfigurerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{Type: knativev1alpha1.ConfigurerConditionReady, Status: "True"},
							},
						},
						URL: &apis.URL{
							Host: "container.default.example.com",
						},
					},
				},
			},
			ExpectOutput: `
NAME        TYPE          REF                 HOST                            STATUS   AGE
app         application   petclinic           app.default.example.com         Ready    <unknown>
container   container     busybox             container.default.example.com   Ready    <unknown>
func        function      square              func.default.example.com        Ready    <unknown>
img         image         projectriff/upper   img.default.example.com         Ready    <unknown>
`,
		},
		{
			Name: "list error",
			Args: []string{},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("list", "configurers"),
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewConfigurerListCommand)
}
