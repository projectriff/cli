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
	"github.com/projectriff/cli/pkg/streaming/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	streamv1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestInMemoryProviderDeleteOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid delete",
			Options: &commands.InMemoryProviderDeleteOptions{
				DeleteOptions: rifftesting.InvalidDeleteOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidDeleteOptionsFieldError,
		},
		{
			Name: "valid delete",
			Options: &commands.InMemoryProviderDeleteOptions{
				DeleteOptions: rifftesting.ValidDeleteOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestInMemoryProviderDeleteCommand(t *testing.T) {
	inmemoryProviderName := "test-inmemory-provider"
	inmemoryProviderOtherName := "test-other-inmemory-provider"
	defaultNamespace := "default"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "delete all in-memory providers",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.InMemoryProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inmemoryProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "inmemoryproviders",
				Namespace: defaultNamespace,
			}},
			ExpectOutput: `
Deleted in-memory providers in namespace "default"
`,
		},
		{
			Name: "delete all in-memory providers error",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.InMemoryProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inmemoryProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete-collection", "inmemoryproviders"),
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "inmemoryproviders",
				Namespace: defaultNamespace,
			}},
			ShouldError: true,
		},
		{
			Name: "delete in-memory providers",
			Args: []string{inmemoryProviderName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.InMemoryProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inmemoryProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "inmemoryproviders",
				Namespace: defaultNamespace,
				Name:      inmemoryProviderName,
			}},
			ExpectOutput: `
Deleted in-memory provider "test-inmemory-provider"
`,
		},
		{
			Name: "delete in-memory provider",
			Args: []string{inmemoryProviderName, inmemoryProviderOtherName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.InMemoryProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inmemoryProviderName,
						Namespace: defaultNamespace,
					},
				},
				&streamv1alpha1.InMemoryProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inmemoryProviderOtherName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "inmemoryproviders",
				Namespace: defaultNamespace,
				Name:      inmemoryProviderName,
			}, {
				Group:     "streaming.projectriff.io",
				Resource:  "inmemoryproviders",
				Namespace: defaultNamespace,
				Name:      inmemoryProviderOtherName,
			}},
			ExpectOutput: `
Deleted in-memory provider "test-inmemory-provider"
Deleted in-memory provider "test-other-inmemory-provider"
`,
		},
		{
			Name: "stream does not exist",
			Args: []string{inmemoryProviderName},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "inmemoryproviders",
				Namespace: defaultNamespace,
				Name:      inmemoryProviderName,
			}},
			ShouldError: true,
		},
		{
			Name: "delete error",
			Args: []string{inmemoryProviderName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.InMemoryProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inmemoryProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete", "inmemoryproviders"),
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "inmemoryproviders",
				Namespace: defaultNamespace,
				Name:      inmemoryProviderName,
			}},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewInMemoryProviderDeleteCommand)
}
