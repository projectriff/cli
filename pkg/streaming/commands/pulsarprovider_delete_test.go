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

func TestPulsarProviderDeleteOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid delete",
			Options: &commands.PulsarProviderDeleteOptions{
				DeleteOptions: rifftesting.InvalidDeleteOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidDeleteOptionsFieldError,
		},
		{
			Name: "valid delete",
			Options: &commands.PulsarProviderDeleteOptions{
				DeleteOptions: rifftesting.ValidDeleteOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestPulsarProviderDeleteCommand(t *testing.T) {
	pulsarProviderName := "test-pulsar-provider"
	pulsarProviderOtherName := "test-other-pulsar-provider"
	defaultNamespace := "default"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "delete all pulsar providers",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "pulsarproviders",
				Namespace: defaultNamespace,
			}},
			ExpectOutput: `
Deleted pulsar providers in namespace "default"
`,
		},
		{
			Name: "delete all pulsar providers error",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete-collection", "pulsarproviders"),
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "pulsarproviders",
				Namespace: defaultNamespace,
			}},
			ShouldError: true,
		},
		{
			Name: "delete pulsar providers",
			Args: []string{pulsarProviderName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "pulsarproviders",
				Namespace: defaultNamespace,
				Name:      pulsarProviderName,
			}},
			ExpectOutput: `
Deleted pulsar provider "test-pulsar-provider"
`,
		},
		{
			Name: "delete pulsar provider",
			Args: []string{pulsarProviderName, pulsarProviderOtherName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderOtherName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "pulsarproviders",
				Namespace: defaultNamespace,
				Name:      pulsarProviderName,
			}, {
				Group:     "streaming.projectriff.io",
				Resource:  "pulsarproviders",
				Namespace: defaultNamespace,
				Name:      pulsarProviderOtherName,
			}},
			ExpectOutput: `
Deleted pulsar provider "test-pulsar-provider"
Deleted pulsar provider "test-other-pulsar-provider"
`,
		},
		{
			Name: "stream does not exist",
			Args: []string{pulsarProviderName},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "pulsarproviders",
				Namespace: defaultNamespace,
				Name:      pulsarProviderName,
			}},
			ShouldError: true,
		},
		{
			Name: "delete error",
			Args: []string{pulsarProviderName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete", "pulsarproviders"),
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "pulsarproviders",
				Namespace: defaultNamespace,
				Name:      pulsarProviderName,
			}},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewPulsarProviderDeleteCommand)
}
