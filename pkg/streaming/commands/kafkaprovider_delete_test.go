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

func TestKafkaProviderDeleteOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid delete",
			Options: &commands.KafkaProviderDeleteOptions{
				DeleteOptions: rifftesting.InvalidDeleteOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidDeleteOptionsFieldError,
		},
		{
			Name: "valid delete",
			Options: &commands.KafkaProviderDeleteOptions{
				DeleteOptions: rifftesting.ValidDeleteOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestKafkaProviderDeleteCommand(t *testing.T) {
	kafkaProviderName := "test-kafka-provider"
	kafkaProviderOtherName := "test-other-kafka-provider"
	defaultNamespace := "default"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "delete all kafka providers",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "kafkaproviders",
				Namespace: defaultNamespace,
			}},
			ExpectOutput: `
Deleted kafka providers in namespace "default"
`,
		},
		{
			Name: "delete all kafka providers error",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete-collection", "kafkaproviders"),
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "kafkaproviders",
				Namespace: defaultNamespace,
			}},
			ShouldError: true,
		},
		{
			Name: "delete kafka providers",
			Args: []string{kafkaProviderName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "kafkaproviders",
				Namespace: defaultNamespace,
				Name:      kafkaProviderName,
			}},
			ExpectOutput: `
Deleted kafka provider "test-kafka-provider"
`,
		},
		{
			Name: "delete kafka provider",
			Args: []string{kafkaProviderName, kafkaProviderOtherName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaProviderName,
						Namespace: defaultNamespace,
					},
				},
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaProviderOtherName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "kafkaproviders",
				Namespace: defaultNamespace,
				Name:      kafkaProviderName,
			}, {
				Group:     "streaming.projectriff.io",
				Resource:  "kafkaproviders",
				Namespace: defaultNamespace,
				Name:      kafkaProviderOtherName,
			}},
			ExpectOutput: `
Deleted kafka provider "test-kafka-provider"
Deleted kafka provider "test-other-kafka-provider"
`,
		},
		{
			Name: "stream does not exist",
			Args: []string{kafkaProviderName},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "kafkaproviders",
				Namespace: defaultNamespace,
				Name:      kafkaProviderName,
			}},
			ShouldError: true,
		},
		{
			Name: "delete error",
			Args: []string{kafkaProviderName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete", "kafkaproviders"),
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "streaming.projectriff.io",
				Resource:  "kafkaproviders",
				Namespace: defaultNamespace,
				Name:      kafkaProviderName,
			}},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewKafkaProviderDeleteCommand)
}
