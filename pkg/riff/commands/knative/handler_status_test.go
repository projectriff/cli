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

package commands_knative_test

import (
	"testing"
	"time"

	knapis "github.com/knative/pkg/apis"
	duckv1beta1 "github.com/knative/pkg/apis/duck/v1beta1"
	commands_knative "github.com/projectriff/cli/pkg/riff/commands/knative"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	requestv1alpha1 "github.com/projectriff/system/pkg/apis/knative/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestHandlerStatusOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands_knative.HandlerStatusOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldError: rifftesting.InvalidResourceOptionsFieldError,
		},
		{
			Name: "valid resource",
			Options: &commands_knative.HandlerStatusOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestHandlerStatusCommand(t *testing.T) {
	defaultNamespace := "default"
	handlerName := "my-handler"

	table := rifftesting.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "show status",
			Args: []string{handlerName},
			GivenObjects: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      handlerName,
						Namespace: defaultNamespace,
					},
					Status: requestv1alpha1.HandlerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{
									Type:    knapis.ConditionReady,
									Status:  corev1.ConditionFalse,
									Reason:  "OopsieDoodle",
									Message: "a hopefully informative message about what went wrong",
									LastTransitionTime: knapis.VolatileTime{
										Inner: metav1.Time{
											Time: time.Date(2019, 6, 29, 01, 44, 05, 0, time.UTC),
										},
									},
								},
							},
						},
					},
				},
			},
			ExpectOutput: `
# my-handler: OopsieDoodle
---
lastTransitionTime: "2019-06-29T01:44:05Z"
message: a hopefully informative message about what went wrong
reason: OopsieDoodle
status: "False"
type: Ready
`,
		},
		{
			Name: "not found",
			Args: []string{handlerName},
			ExpectOutput: `
Handler "default/my-handler" not found
`,
			ShouldError: true,
		},
		{
			Name: "get error",
			Args: []string{handlerName},
			GivenObjects: []runtime.Object{
				&requestv1alpha1.Handler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      handlerName,
						Namespace: defaultNamespace,
					},
					Status: requestv1alpha1.HandlerStatus{
						Status: duckv1beta1.Status{
							Conditions: duckv1beta1.Conditions{
								{
									Type:    knapis.ConditionReady,
									Status:  corev1.ConditionFalse,
									Reason:  "OopsieDoodle",
									Message: "a hopefully informative message about what went wrong",
									LastTransitionTime: knapis.VolatileTime{
										Inner: metav1.Time{
											Time: time.Date(2019, 6, 29, 01, 44, 05, 0, time.UTC),
										},
									},
								},
							},
						},
					},
				},
			},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("get", "handlers"),
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands_knative.NewHandlerStatusCommand)
}
