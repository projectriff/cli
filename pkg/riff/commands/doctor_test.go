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
			Name: "istio-system is missing",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: "knative-build"},
				},
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"},
				},
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: "riff-system"},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				passAccessReview(),
			},
			ExpectOutput: `
Namespace "istio-system"      Missing
Namespace "knative-build"     OK
Namespace "knative-serving"   OK
Namespace "riff-system"       OK

Installation is not healthy
`,
		},
		{
			Name: "multiple namespaces are missing",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: "knative-build"},
				},
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				passAccessReview(),
			},
			ExpectOutput: `
Namespace "istio-system"      Missing
Namespace "knative-build"     OK
Namespace "knative-serving"   OK
Namespace "riff-system"       Missing

Installation is not healthy
`,
		},
		{
			Name: "custom resource definitions are missing",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "istio-system"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-build"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "knative-serving"}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "riff-system"}},
			},
			ExpectCreates: merge(
				selfSubjectAccessReviewRequests("default", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("default", "core", "secrets", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				passAccessReview(),
			},
			ExpectOutput: `
Namespace "istio-system"      OK
Namespace "knative-build"     OK
Namespace "knative-serving"   OK
Namespace "riff-system"       OK

RESOURCE                            READ      WRITE
configmaps                          allowed   allowed   
secrets                             allowed   allowed   
applications.build.projectriff.io   missing   missing   
functions.build.projectriff.io      missing   missing   
handlers.request.projectriff.io     missing   missing   
processors.stream.projectriff.io    missing   missing   
streams.stream.projectriff.io       missing   missing   

Installation is not healthy
`,
		},
		{
			Name: "read and write KO for functions in specified namespace",
			Args: []string{cli.NamespaceFlagName, "foo"},
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
				selfSubjectAccessReviewRequests("foo", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("foo", "core", "secrets", verbs...),
				selfSubjectAccessReviewRequests("foo", "build.projectriff.io", "applications", verbs...),
				selfSubjectAccessReviewRequests("foo", "build.projectriff.io", "functions", verbs...),
				selfSubjectAccessReviewRequests("foo", "request.projectriff.io", "handlers", verbs...),
				selfSubjectAccessReviewRequests("foo", "stream.projectriff.io", "processors", verbs...),
				selfSubjectAccessReviewRequests("foo", "stream.projectriff.io", "streams", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				failAccessReviewOn("functions", "*"),
				passAccessReview(),
			},
			ExpectOutput: `
Namespace "istio-system"      OK
Namespace "knative-build"     OK
Namespace "knative-serving"   OK
Namespace "riff-system"       OK

RESOURCE                            READ      WRITE
configmaps                          allowed   allowed   
secrets                             allowed   allowed   
applications.build.projectriff.io   allowed   allowed   
functions.build.projectriff.io      denied    denied    
handlers.request.projectriff.io     allowed   allowed   
processors.stream.projectriff.io    allowed   allowed   
streams.stream.projectriff.io       allowed   allowed   

Installation is not healthy
`,
		},
		{
			Name: "read status mixed for handlers, write status mixed for streams in specified namespace",
			Args: []string{cli.NamespaceFlagName, "foo"},
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
				selfSubjectAccessReviewRequests("foo", "core", "configmaps", verbs...),
				selfSubjectAccessReviewRequests("foo", "core", "secrets", verbs...),
				selfSubjectAccessReviewRequests("foo", "build.projectriff.io", "applications", verbs...),
				selfSubjectAccessReviewRequests("foo", "build.projectriff.io", "functions", verbs...),
				selfSubjectAccessReviewRequests("foo", "request.projectriff.io", "handlers", verbs...),
				selfSubjectAccessReviewRequests("foo", "stream.projectriff.io", "processors", verbs...),
				selfSubjectAccessReviewRequests("foo", "stream.projectriff.io", "streams", verbs...),
			),
			WithReactors: []rifftesting.ReactionFunc{
				failAccessReviewOn("handlers", "get"),
				failAccessReviewOn("streams", "update"),
				passAccessReview(),
			},
			ExpectOutput: `
Namespace "istio-system"      OK
Namespace "knative-build"     OK
Namespace "knative-serving"   OK
Namespace "riff-system"       OK

RESOURCE                            READ      WRITE
configmaps                          allowed   allowed   
secrets                             allowed   allowed   
applications.build.projectriff.io   allowed   allowed   
functions.build.projectriff.io      allowed   allowed   
handlers.request.projectriff.io     mixed     allowed   
processors.stream.projectriff.io    allowed   allowed   
streams.stream.projectriff.io       allowed   mixed     

Installation is not healthy
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

func failAccessReviewOn(resource string, verb string) func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
	return func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
		if !action.Matches("create", "selfsubjectaccessreviews") {
			return false, nil, nil
		}
		creationAction, _ := action.(clientgotesting.CreateAction)
		review, _ := creationAction.GetObject().(*authorizationv1.SelfSubjectAccessReview)
		if review.Spec.ResourceAttributes.Resource != resource || (verb != "*" && review.Spec.ResourceAttributes.Verb != verb) {
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
