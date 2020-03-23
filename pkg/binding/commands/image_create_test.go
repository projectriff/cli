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

	bindingsv1alpha1 "github.com/projectriff/bindings/pkg/apis/bindings/v1alpha1"
	"github.com/projectriff/cli/pkg/binding/commands"
	"github.com/projectriff/cli/pkg/cli"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"knative.dev/pkg/tracker"
)

func TestImageBindingCreateOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "empty resource",
			Options: &commands.ImageCreateOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidResourceOptionsFieldError.Also(
				cli.ErrInvalidValue("", cli.SubjectFlagName),
				cli.ErrInvalidValue("", cli.ProviderFlagName),
			),
		},
		{
			Name: "valid resource",
			Options: &commands.ImageCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Subject:         "deployment:my-deployment",
				Provider:        "function:my-function:user-container",
			},
			ShouldValidate: true,
		},
		{
			Name: "invalid subject",
			Options: &commands.ImageCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Subject:         "foo",
				Provider:        "function:my-function:user-container",
			},
			ExpectFieldErrors: cli.ErrInvalidValue("foo", cli.SubjectFlagName),
		},
		{
			Name: "invalid providers",
			Options: &commands.ImageCreateOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				Subject:         "deployment:my-deployment",
				Provider:        "foo",
			},
			ExpectFieldErrors: cli.ErrInvalidValue("foo", cli.ProviderFlagName),
		},
	}

	table.Run(t)
}

func TestImageBindingCreateCommand(t *testing.T) {
	defaultNamespace := "default"
	imageBindingName := "my-image-binding"
	functionName := "my-function"
	deploymentName := "my-deployment"
	containerName := "user-container"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "create",
			Args: []string{imageBindingName, cli.SubjectFlagName, "deployments.apps:my-deployment", cli.ProviderFlagName, "functions.build.projectriff.io:my-function:user-container"},
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config) (context.Context, error) {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				addTestDiscoveryResources(discovery)
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, config *cli.Config) error {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				discovery.Resources = []*metav1.APIResourceList{}
				return nil
			},
			ExpectCreates: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      imageBindingName,
					},
					Spec: bindingsv1alpha1.ImageBindingSpec{
						Subject: &tracker.Reference{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Namespace:  defaultNamespace,
							Name:       deploymentName,
						},
						Providers: []bindingsv1alpha1.ImageProvider{
							{
								ImageableRef: &tracker.Reference{
									APIVersion: "build.projectriff.io/v1alpha1",
									Kind:       "Function",
									Namespace:  defaultNamespace,
									Name:       functionName,
								},
								ContainerName: containerName,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created image binding "my-image-binding"
`,
		},
		{
			Name: "create, dry run",
			Args: []string{imageBindingName, cli.SubjectFlagName, "deployments.apps:my-deployment", cli.ProviderFlagName, "functions.build.projectriff.io:my-function:user-container", cli.DryRunFlagName},
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config) (context.Context, error) {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				addTestDiscoveryResources(discovery)
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, config *cli.Config) error {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				discovery.Resources = []*metav1.APIResourceList{}
				return nil
			},
			ExpectOutput: `
---
apiVersion: bindings.projectriff.io/v1alpha1
kind: ImageBinding
metadata:
  creationTimestamp: null
  name: my-image-binding
  namespace: default
spec:
  providers:
  - containerName: user-container
    imageableRef:
      apiVersion: build.projectriff.io/v1alpha1
      kind: Function
      name: my-function
      namespace: default
  subject:
    apiVersion: apps/v1
    kind: Deployment
    name: my-deployment
    namespace: default
status: {}

Created image binding "my-image-binding"
`,
		},
		{
			Name: "create, unknown subject",
			Args: []string{imageBindingName, cli.SubjectFlagName, "foo:my-foo", cli.ProviderFlagName, "functions.build.projectriff.io:my-function:user-container"},
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config) (context.Context, error) {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				addTestDiscoveryResources(discovery)
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, config *cli.Config) error {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				discovery.Resources = []*metav1.APIResourceList{}
				return nil
			},
			ShouldError: true,
		},
		{
			Name: "error existing image binding",
			Args: []string{imageBindingName, cli.SubjectFlagName, "deployments.apps:my-deployment", cli.ProviderFlagName, "functions.build.projectriff.io:my-function:user-container"},
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config) (context.Context, error) {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				addTestDiscoveryResources(discovery)
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, config *cli.Config) error {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				discovery.Resources = []*metav1.APIResourceList{}
				return nil
			},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      imageBindingName,
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      imageBindingName,
					},
					Spec: bindingsv1alpha1.ImageBindingSpec{
						Subject: &tracker.Reference{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Namespace:  defaultNamespace,
							Name:       deploymentName,
						},
						Providers: []bindingsv1alpha1.ImageProvider{
							{
								ImageableRef: &tracker.Reference{
									APIVersion: "build.projectriff.io/v1alpha1",
									Kind:       "Function",
									Namespace:  defaultNamespace,
									Name:       functionName,
								},
								ContainerName: containerName,
							},
						},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error during create",
			Args: []string{imageBindingName, cli.SubjectFlagName, "deployments.apps:my-deployment", cli.ProviderFlagName, "functions.build.projectriff.io:my-function:user-container"},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("create", "imagebindings"),
			},
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config) (context.Context, error) {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				addTestDiscoveryResources(discovery)
				return ctx, nil
			},
			CleanUp: func(t *testing.T, ctx context.Context, config *cli.Config) error {
				discovery := config.Client.Discovery().(*fakediscovery.FakeDiscovery)
				discovery.Resources = []*metav1.APIResourceList{}
				return nil
			},
			ExpectCreates: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      imageBindingName,
					},
					Spec: bindingsv1alpha1.ImageBindingSpec{
						Subject: &tracker.Reference{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Namespace:  defaultNamespace,
							Name:       deploymentName,
						},
						Providers: []bindingsv1alpha1.ImageProvider{
							{
								ImageableRef: &tracker.Reference{
									APIVersion: "build.projectriff.io/v1alpha1",
									Kind:       "Function",
									Namespace:  defaultNamespace,
									Name:       functionName,
								},
								ContainerName: containerName,
							},
						},
					},
				},
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewImageCreateCommand)
}

func addTestDiscoveryResources(discovery *fakediscovery.FakeDiscovery) {
	discovery.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "apps/v1",
			APIResources: []metav1.APIResource{
				{
					Name: "deployments",
					Kind: "Deployment",
				},
			},
		},
		{
			GroupVersion: "build.projectriff.io/v1alpha1",
			APIResources: []metav1.APIResource{
				{
					Name: "functions",
					Kind: "Function",
				},
			},
		},
	}
}
