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

	bindingsv1alpha1 "github.com/projectriff/bindings/pkg/apis/bindings/v1alpha1"
	"github.com/projectriff/cli/pkg/binding/commands"
	"github.com/projectriff/cli/pkg/cli"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestImageBindingDeleteOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid delete",
			Options: &commands.ImageDeleteOptions{
				DeleteOptions: rifftesting.InvalidDeleteOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidDeleteOptionsFieldError,
		},
		{
			Name: "valid delete",
			Options: &commands.ImageDeleteOptions{
				DeleteOptions: rifftesting.ValidDeleteOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestImageBindingDeleteCommand(t *testing.T) {
	imageBindingName := "test-image-binding"
	imageBindingOtherName := "test-other-image-binding"
	defaultNamespace := "default"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "delete all image bindings",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "bindings.projectriff.io",
				Resource:  "imagebindings",
				Namespace: defaultNamespace,
			}},
			ExpectOutput: `
Deleted image bindings in namespace "default"
`,
		},
		{
			Name: "delete all image bindings error",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete-collection", "imagebindings"),
			},
			ExpectDeleteCollections: []rifftesting.DeleteCollectionRef{{
				Group:     "bindings.projectriff.io",
				Resource:  "imagebindings",
				Namespace: defaultNamespace,
			}},
			ShouldError: true,
		},
		{
			Name: "delete image bindings",
			Args: []string{imageBindingName},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "bindings.projectriff.io",
				Resource:  "imagebindings",
				Namespace: defaultNamespace,
				Name:      imageBindingName,
			}},
			ExpectOutput: `
Deleted image binding "test-image-binding"
`,
		},
		{
			Name: "delete image bindings",
			Args: []string{imageBindingName, imageBindingOtherName},
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
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "bindings.projectriff.io",
				Resource:  "imagebindings",
				Namespace: defaultNamespace,
				Name:      imageBindingName,
			}, {
				Group:     "bindings.projectriff.io",
				Resource:  "imagebindings",
				Namespace: defaultNamespace,
				Name:      imageBindingOtherName,
			}},
			ExpectOutput: `
Deleted image binding "test-image-binding"
Deleted image binding "test-other-image-binding"
`,
		},
		{
			Name: "image binding does not exist",
			Args: []string{imageBindingName},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "bindings.projectriff.io",
				Resource:  "imagebindings",
				Namespace: defaultNamespace,
				Name:      imageBindingName,
			}},
			ShouldError: true,
		},
		{
			Name: "delete error",
			Args: []string{imageBindingName},
			GivenObjects: []runtime.Object{
				&bindingsv1alpha1.ImageBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      imageBindingName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("delete", "imagebindings"),
			},
			ExpectDeletes: []rifftesting.DeleteRef{{
				Group:     "bindings.projectriff.io",
				Resource:  "imagebindings",
				Namespace: defaultNamespace,
				Name:      imageBindingName,
			}},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewImageDeleteCommand)
}
