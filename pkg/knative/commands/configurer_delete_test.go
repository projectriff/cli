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
	"testing"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/knative/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	knativev1alpha1 "github.com/projectriff/system/pkg/apis/knative/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestConfigurerDeleteOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid delete",
			Options: &commands.ConfigurerDeleteOptions{
				DeleteOptions: rifftesting.InvalidDeleteOptions,
			},
			ExpectFieldError: rifftesting.InvalidDeleteOptionsFieldError,
		},
		{
			Name: "valid delete",
			Options: &commands.ConfigurerDeleteOptions{
				DeleteOptions: rifftesting.ValidDeleteOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestConfigurerDeleteCommand(t *testing.T) {
	configurerName := "test-configurer"
	configurerOtherName := "test-other-configurer"
	defaultNamespace := "default"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "delete all configurers",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "knative.projectriff.io",
				Resource:  "configurers",
				Namespace: defaultNamespace,
			}},
			ExpectOutput: `
Deleted configurers in namespace "default"
`,
		},
		{
			Name: "delete all configurers error",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete-collection", "configurers"),
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "knative.projectriff.io",
				Resource:  "configurers",
				Namespace: defaultNamespace,
			}},
			ShouldError: true,
		},
		{
			Name: "delete configurer",
			Args: []string{configurerName},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "knative.projectriff.io",
				Resource:  "configurers",
				Namespace: defaultNamespace,
				Name:      configurerName,
			}},
			ExpectOutput: `
Deleted configurer "test-configurer"
`,
		},
		{
			Name: "delete configurers",
			Args: []string{configurerName, configurerOtherName},
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
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "knative.projectriff.io",
				Resource:  "configurers",
				Namespace: defaultNamespace,
				Name:      configurerName,
			}, {
				Group:     "knative.projectriff.io",
				Resource:  "configurers",
				Namespace: defaultNamespace,
				Name:      configurerOtherName,
			}},
			ExpectOutput: `
Deleted configurer "test-configurer"
Deleted configurer "test-other-configurer"
`,
		},
		{
			Name: "configurer does not exist",
			Args: []string{configurerName},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "knative.projectriff.io",
				Resource:  "configurers",
				Namespace: defaultNamespace,
				Name:      configurerName,
			}},
			ShouldError: true,
		},
		{
			Name: "delete error",
			Args: []string{configurerName},
			GivenObjects: []runtime.Object{
				&knativev1alpha1.Configurer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configurerName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete", "configurers"),
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "knative.projectriff.io",
				Resource:  "configurers",
				Namespace: defaultNamespace,
				Name:      configurerName,
			}},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewConfigurerDeleteCommand)
}
