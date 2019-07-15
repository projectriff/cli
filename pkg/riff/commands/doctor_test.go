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
	table := rifftesting.CommandTable{
		{
			Name: "not installed",
			Args: []string{},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				passAccessReview(),
			},
			ShouldError: true,
			ExpectOutput: `
NAMESPACE         STATUS
istio-system      missing
knative-build     missing
knative-serving   missing
riff-system       missing

RESOURCE                            READ      WRITE
configmaps                          allowed   allowed
secrets                             allowed   allowed
applications.build.projectriff.io   missing   missing
functions.build.projectriff.io      missing   missing
handlers.request.projectriff.io     missing   missing
processors.stream.projectriff.io    missing   missing
streams.stream.projectriff.io       missing   missing

Installation is not OK
`,
		},
		{
			Name: "installed",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", verbs...),
				selfSubjectAccessReviewRequests("default", "request.projectriff.io", "handlers", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "processors", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "streams", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				passAccessReview(),
			},
			ExpectOutput: `
NAMESPACE         STATUS
istio-system      ok
knative-build     ok
knative-serving   ok
riff-system       ok

RESOURCE                            READ      WRITE
configmaps                          allowed   allowed
secrets                             allowed   allowed
applications.build.projectriff.io   allowed   allowed
functions.build.projectriff.io      allowed   allowed
handlers.request.projectriff.io     allowed   allowed
processors.stream.projectriff.io    allowed   allowed
streams.stream.projectriff.io       allowed   allowed

Installation is OK
`,
		},
		{
			Name: "read-only access",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", verbs...),
				selfSubjectAccessReviewRequests("default", "request.projectriff.io", "handlers", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "processors", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "streams", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				denyAccessReviewOn("*", "create"),
				denyAccessReviewOn("*", "update"),
				denyAccessReviewOn("*", "delete"),
				denyAccessReviewOn("*", "patch"),
				passAccessReview(),
			},
			ShouldError: true,
			ExpectOutput: `
NAMESPACE         STATUS
istio-system      ok
knative-build     ok
knative-serving   ok
riff-system       ok

RESOURCE                            READ      WRITE
configmaps                          allowed   denied
secrets                             allowed   denied
applications.build.projectriff.io   allowed   denied
functions.build.projectriff.io      allowed   denied
handlers.request.projectriff.io     allowed   denied
processors.stream.projectriff.io    allowed   denied
streams.stream.projectriff.io       allowed   denied

Installation is not OK
`,
		},
		{
			Name: "no-watch access",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", verbs...),
				selfSubjectAccessReviewRequests("default", "request.projectriff.io", "handlers", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "processors", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "streams", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				denyAccessReviewOn("*", "watch"),
				passAccessReview(),
			},
			ShouldError: true,
			ExpectOutput: `
NAMESPACE         STATUS
istio-system      ok
knative-build     ok
knative-serving   ok
riff-system       ok

RESOURCE                            READ    WRITE
configmaps                          mixed   allowed
secrets                             mixed   allowed
applications.build.projectriff.io   mixed   allowed
functions.build.projectriff.io      mixed   allowed
handlers.request.projectriff.io     mixed   allowed
processors.stream.projectriff.io    mixed   allowed
streams.stream.projectriff.io       mixed   allowed

Installation is not OK
`,
		},
		{
			Name: "error listing namespaces",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
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
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", verbs...),
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
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
			},
			ExpectCreates: selfSubjectAccessReviewRequests("default", "core", "configmaps", "get"),
			WithReactors: []rifftesting.ReactionFunc{
				failAccessReview(),
			},
			ShouldError: true,
		},
		{
			Name: "error evaluating selfsubjectaccessreview",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
			},
			ExpectCreates: selfSubjectAccessReviewRequests("default", "core", "configmaps", "get"),
			WithReactors: []rifftesting.ReactionFunc{
				failAccessReviewEvaluationOn("*", "*"),
			},
			ShouldError: true,
		},
		{
			Name: "error evaluating selfsubjectaccessreview",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "applications.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "functions.build.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "handlers.request.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "processors.stream.projectriff.io"}},
				&apiextensionsv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: "streams.stream.projectriff.io"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "applications", verbs...),
				selfSubjectAccessReviewRequests("default", "build.projectriff.io", "functions", verbs...),
				selfSubjectAccessReviewRequests("default", "request.projectriff.io", "handlers", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "processors", verbs...),
				selfSubjectAccessReviewRequests("default", "stream.projectriff.io", "streams", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				unknownAccessReviewOn("*", "*"),
			},
			ShouldError: true,
			ExpectOutput: `
NAMESPACE         STATUS
istio-system      ok
knative-build     ok
knative-serving   ok
riff-system       ok

RESOURCE                            READ      WRITE
configmaps                          unknown   unknown
secrets                             unknown   unknown
applications.build.projectriff.io   unknown   unknown
functions.build.projectriff.io      unknown   unknown
handlers.request.projectriff.io     unknown   unknown
processors.stream.projectriff.io    unknown   unknown
streams.stream.projectriff.io       unknown   unknown

Installation is not OK
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

func selfSubjectAccessReviewRequests(namespace, group, resource string, verbs ...string) []runtime.Object {
	result := make([]runtime.Object, len(verbs))
	for i, verb := range verbs {
		result[i] = &authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authorizationv1.ResourceAttributes{
					Namespace: namespace,
					Group:     group,
					Verb:      verb,
					Resource:  resource,
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
