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
	"fmt"
	"testing"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/riff/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	authorizationv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgotesting "k8s.io/client-go/testing"
)

func TestDoctorOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "valid",
			Options: &commands.DoctorOptions{
				Namespace: "default",
			},
			ShouldValidate: true,
		},
		{
			Name:             "invalid",
			Options:          &commands.DoctorOptions{},
			ExpectFieldError: cli.ErrMissingField(cli.NamespaceFlagName),
		},
	}

	table.Run(t)
}

func TestDoctorCommand(t *testing.T) {
	verbs := []string{"get", "list", "create", "update", "delete", "patch", "watch"}
	readVerbs := []string{"get", "list", "watch"}
	table := rifftesting.CommandTable{
		{
			Name: "not installed",
			Args: []string{},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "", readVerbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "log", readVerbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				passAccessReview(),
			},
			ExpectOutput: `
NAMESPACE     STATUS
riff-system   missing

RESOURCE                              READ      WRITE
configmaps                            allowed   allowed
secrets                               allowed   allowed
pods                                  allowed   n/a
pods/log                              allowed   n/a
applications.build.projectriff.io     missing   missing
containers.build.projectriff.io       missing   missing
functions.build.projectriff.io        missing   missing
deployers.core.projectriff.io         missing   missing
processors.streaming.projectriff.io   missing   missing
streams.streaming.projectriff.io      missing   missing
adapters.knative.projectriff.io       missing   missing
deployers.knative.projectriff.io      missing   missing
`,
		},
		{
			Name: "installed",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "", readVerbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "log", readVerbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "containers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core.projectriff.io", "deployers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "processors", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "streams", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "adapters", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "deployers", "", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				passAccessReview(),
			},
			ExpectOutput: `
NAMESPACE     STATUS
riff-system   ok

RESOURCE                              READ      WRITE
configmaps                            allowed   allowed
secrets                               allowed   allowed
pods                                  allowed   n/a
pods/log                              allowed   n/a
applications.build.projectriff.io     allowed   allowed
containers.build.projectriff.io       allowed   allowed
functions.build.projectriff.io        allowed   allowed
deployers.core.projectriff.io         allowed   allowed
processors.streaming.projectriff.io   allowed   allowed
streams.streaming.projectriff.io      allowed   allowed
adapters.knative.projectriff.io       allowed   allowed
deployers.knative.projectriff.io      allowed   allowed
`,
		},
		{
			Name: "read-only access",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "", readVerbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "log", readVerbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "containers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core.projectriff.io", "deployers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "processors", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "streams", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "adapters", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "deployers", "", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				denyAccessReviewOn("*", "create"),
				denyAccessReviewOn("*", "update"),
				denyAccessReviewOn("*", "delete"),
				denyAccessReviewOn("*", "patch"),
				passAccessReview(),
			},
			ExpectOutput: `
NAMESPACE     STATUS
riff-system   ok

RESOURCE                              READ      WRITE
configmaps                            allowed   denied
secrets                               allowed   denied
pods                                  allowed   n/a
pods/log                              allowed   n/a
applications.build.projectriff.io     allowed   denied
containers.build.projectriff.io       allowed   denied
functions.build.projectriff.io        allowed   denied
deployers.core.projectriff.io         allowed   denied
processors.streaming.projectriff.io   allowed   denied
streams.streaming.projectriff.io      allowed   denied
adapters.knative.projectriff.io       allowed   denied
deployers.knative.projectriff.io      allowed   denied
`,
		},
		{
			Name: "no-watch access",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "", readVerbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "log", readVerbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "containers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core.projectriff.io", "deployers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "processors", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "streams", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "adapters", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "deployers", "", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				denyAccessReviewOn("*", "watch"),
				passAccessReview(),
			},
			ExpectOutput: `
NAMESPACE     STATUS
riff-system   ok

RESOURCE                              READ    WRITE
configmaps                            mixed   allowed
secrets                               mixed   allowed
pods                                  mixed   n/a
pods/log                              mixed   n/a
applications.build.projectriff.io     mixed   allowed
containers.build.projectriff.io       mixed   allowed
functions.build.projectriff.io        mixed   allowed
deployers.core.projectriff.io         mixed   allowed
processors.streaming.projectriff.io   mixed   allowed
streams.streaming.projectriff.io      mixed   allowed
adapters.knative.projectriff.io       mixed   allowed
deployers.knative.projectriff.io      mixed   allowed
`,
		},
		{
			Name: "error listing namespaces",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("list", "namespaces"),
			},
			ShouldError: true,
		},
		{
			Name: "error getting crd",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "", readVerbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "log", readVerbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("get", "customresourcedefinitions"),
				passAccessReview(),
			},
			ShouldError: true,
		},
		{
			Name: "error creating selfsubjectaccessreview",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			ExpectCreates: selfSubjectAccessReviewRequests("default", "core", "configmaps", "", "get"),
			WithReactors: []rifftesting.ReactionFunc{
				failAccessReview(),
			},
			ShouldError: true,
		},
		{
			Name: "error evaluating selfsubjectaccessreview",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			ExpectCreates: selfSubjectAccessReviewRequests("default", "core", "configmaps", "", "get"),
			WithReactors: []rifftesting.ReactionFunc{
				failAccessReviewEvaluationOn("*", "*"),
			},
			ShouldError: true,
		},
		{
			Name: "error evaluating selfsubjectaccessreview",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "containers.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.core.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.streaming.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "adapters.knative.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "deployers.knative.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "", readVerbs...),
				selfSubjectAccessReviewRequests("default", "core", "pods", "log", readVerbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "containers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", "", verbs...),
				selfSubjectAccessReviewRequests("default", "core.projectriff.io", "deployers", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "processors", "", verbs...),
				selfSubjectAccessReviewRequests("default", "streaming.projectriff.io", "streams", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "adapters", "", verbs...),
				selfSubjectAccessReviewRequests("default", "knative.projectriff.io", "deployers", "", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				unknownAccessReviewOn("*", "*"),
			},
			ExpectOutput: `
NAMESPACE     STATUS
riff-system   ok

RESOURCE                              READ      WRITE
configmaps                            unknown   unknown
secrets                               unknown   unknown
pods                                  unknown   n/a
pods/log                              unknown   n/a
applications.build.projectriff.io     unknown   unknown
containers.build.projectriff.io       unknown   unknown
functions.build.projectriff.io        unknown   unknown
deployers.core.projectriff.io         unknown   unknown
processors.streaming.projectriff.io   unknown   unknown
streams.streaming.projectriff.io      unknown   unknown
adapters.knative.projectriff.io       unknown   unknown
deployers.knative.projectriff.io      unknown   unknown
`,
		},
	}

	table.Run(t, commands.NewDoctorCommand)
}

func merge(objectSets ...[]runtime.Object) []runtime.Object {
	var result []runtime.Object
	for _, objects := range objectSets {
		result = append(result, objects...)
	}
	return result
}

func selfSubjectAccessReviewRequests(namespace, group, resource string, subresource string, verbs ...string) []runtime.Object {
	result := make([]runtime.Object, len(verbs))
	for i, verb := range verbs {
		result[i] = &authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authorizationv1.ResourceAttributes{
					Namespace:   namespace,
					Group:       group,
					Verb:        verb,
					Resource:    resource,
					Subresource: subresource,
				},
			},
		}
	}
	return result
}

func failAccessReview() func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
	// fork of riffesting.InduceFailure that returns the review to avoid a panic in the fake client
	return func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
		if !action.Matches("create", "selfsubjectaccessreviews") {
			return false, nil, nil
		}
		creationAction, _ := action.(clientgotesting.CreateAction)
		review, _ := creationAction.GetObject().(*authorizationv1.SelfSubjectAccessReview)
		return true, review, fmt.Errorf("inducing failure for %s %s", action.GetVerb(), action.GetResource().Resource)
	}
}

func failAccessReviewEvaluationOn(resource string, verb string) func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
	return func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
		if !action.Matches("create", "selfsubjectaccessreviews") {
			return false, nil, nil
		}
		creationAction, _ := action.(clientgotesting.CreateAction)
		review, _ := creationAction.GetObject().(*authorizationv1.SelfSubjectAccessReview)
		if (resource != "*" && review.Spec.ResourceAttributes.Resource != resource) || (verb != "*" && review.Spec.ResourceAttributes.Verb != verb) {
			return false, nil, nil
		}
		review = review.DeepCopy()
		review.Status.EvaluationError = "induced failure"
		return true, review, nil
	}
}
func unknownAccessReviewOn(resource string, verb string) func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
	return func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
		if !action.Matches("create", "selfsubjectaccessreviews") {
			return false, nil, nil
		}
		creationAction, _ := action.(clientgotesting.CreateAction)
		review, _ := creationAction.GetObject().(*authorizationv1.SelfSubjectAccessReview)
		if (resource != "*" && review.Spec.ResourceAttributes.Resource != resource) || (verb != "*" && review.Spec.ResourceAttributes.Verb != verb) {
			return false, nil, nil
		}
		review = review.DeepCopy()
		return true, review, nil
	}
}

func denyAccessReviewOn(resource string, verb string) func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
	return func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
		if !action.Matches("create", "selfsubjectaccessreviews") {
			return false, nil, nil
		}
		creationAction, _ := action.(clientgotesting.CreateAction)
		review, _ := creationAction.GetObject().(*authorizationv1.SelfSubjectAccessReview)
		if (resource != "*" && review.Spec.ResourceAttributes.Resource != resource) || (verb != "*" && review.Spec.ResourceAttributes.Verb != verb) {
			return false, nil, nil
		}
		review = review.DeepCopy()
		review.Status.Denied = true
		return true, review, nil
	}
}

func passAccessReview() func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
	return func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
		if !action.Matches("create", "selfsubjectaccessreviews") {
			return false, nil, nil
		}
		creationAction, _ := action.(clientgotesting.CreateAction)
		review, _ := creationAction.GetObject().(*authorizationv1.SelfSubjectAccessReview)
		review = review.DeepCopy()
		review.Status.Allowed = true
		return true, review, nil
	}
}
