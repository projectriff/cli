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

func TestKafkaProviderCreateOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.KafkaProviderCreateOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidResourceOptionsFieldError.Also(
				cli.ErrMissingField(cli.BootstrapServersFlagName),
			),
		},
		{
			Name: "valid provider",
			Options: &commands.KafkaProviderCreateOptions{
				ResourceOptions:  rifftesting.ValidResourceOptions,
				BootstrapServers: "localhost:9092",
			},
			ShouldValidate: true,
		},
		{
			Name: "dry run",
			Options: &commands.KafkaProviderCreateOptions{
				ResourceOptions:  rifftesting.ValidResourceOptions,
				BootstrapServers: "localhost:9092",
				DryRun:           true,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestKafkaProviderCreateCommand(t *testing.T) {
	defaultNamespace := "default"
	kafkaProviderName := "my-kafka-provider"
	bootstrapServers := "localhost:9092"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "kafka provider",
			Args: []string{kafkaProviderName, cli.BootstrapServersFlagName, bootstrapServers},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      kafkaProviderName,
					},
					Spec: streamv1alpha1.KafkaProviderSpec{
						BootstrapServers: bootstrapServers,
					},
				},
			},
			ExpectOutput: `
Created kafka provider "my-kafka-provider"
`,
		},
		{
			Name: "dry run",
			Args: []string{kafkaProviderName, cli.BootstrapServersFlagName, bootstrapServers, cli.DryRunFlagName},
			ExpectOutput: `
---
apiVersion: streaming.projectriff.io/v1alpha1
kind: KafkaProvider
metadata:
  creationTimestamp: null
  name: my-kafka-provider
  namespace: default
spec:
  bootstrapServers: localhost:9092
status: {}

Created kafka provider "my-kafka-provider"
`,
		},
		{
			Name: "error existing provider",
			Args: []string{kafkaProviderName, cli.BootstrapServersFlagName, bootstrapServers},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      kafkaProviderName,
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      kafkaProviderName,
					},
					Spec: streamv1alpha1.KafkaProviderSpec{
						BootstrapServers: bootstrapServers,
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error during create",
			Args: []string{kafkaProviderName, cli.BootstrapServersFlagName, bootstrapServers},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("create", "kafkaproviders"),
			},
			ExpectCreates: []runtime.Object{
				&streamv1alpha1.KafkaProvider{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      kafkaProviderName,
					},
					Spec: streamv1alpha1.KafkaProviderSpec{
						BootstrapServers: bootstrapServers,
					},
				},
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewKafkaProviderCreateCommand)
}
