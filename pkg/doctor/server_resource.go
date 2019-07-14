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

package doctor

import (
	"fmt"
	"io"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/cli/printers"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewReview(ns string, gr schema.GroupResource, verb string) *authv1.SelfSubjectAccessReview {
	return &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Namespace: ns,
				Group:     gr.Group,
				Resource:  gr.Resource,
				Verb:      verb,
			},
		},
	}
}

type AccessChecks struct {
	Resource schema.GroupResource
	Verbs    []string
}

func IsRead(verb string) bool {
	return verb == "get" || verb == "list" || verb == "watch"
}

type AccessSummary struct {
	Statuses []Status
}

func (as *AccessSummary) IsHealthy() bool {
	for _, status := range as.Statuses {
		if status.ReadStatus != AccessAllowed || status.WriteStatus != AccessAllowed {
			return false
		}
	}
	return true
}

func (as *AccessSummary) Fprint(out io.Writer) {
	printer := printers.GetNewTabWriter(out)
	defer printer.Flush()
	fmt.Fprintf(printer, "RESOURCE\tREAD\tWRITE\n")
	for _, status := range as.Statuses {
		resource := status.Resource.Resource
		if status.Resource.Group != "core" {
			resource = fmt.Sprintf("%s.%s", resource, status.Resource.Group)
		}
		fmt.Fprintf(printer, "%s\t%s\t%s\n", resource, status.ReadStatus.String(), status.WriteStatus.String())
	}
}

type Status struct {
	Resource    schema.GroupResource
	ReadStatus  AccessStatus
	WriteStatus AccessStatus
}

type AccessStatus int

const (
	AccessUndefined AccessStatus = iota
	AccessAllowed                /* right is granted */
	AccessDenied                 /* right is denied */
	AccessMixed                  /* for the same resource, some rights are granted, some are denied */
	AccessMissing                /* resource not deployed */
)

func (as AccessStatus) Combine(new AccessStatus) AccessStatus {
	if as == AccessUndefined {
		return new
	}
	if as != new {
		return AccessMixed
	}
	if as == AccessAllowed {
		return AccessAllowed
	}
	return AccessDenied
}

func (as *AccessStatus) String() string {
	status := *as
	switch status {
	case AccessAllowed:
		return cli.Ssuccessf("allowed")
	case AccessMixed:
		return cli.Swarnf("mixed")
	case AccessDenied:
		return cli.Swarnf("denied")
	case AccessMissing:
		return cli.Serrorf("missing")
	}
	panic(fmt.Sprintf("Unsupported value %v", status))
}
