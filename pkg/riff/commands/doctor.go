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

package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/cli/printers"
	"github.com/projectriff/cli/pkg/riff/resource"
	"github.com/spf13/cobra"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type DoctorOptions struct {
	cli.NamespaceOptions
}

var (
	_ cli.Validatable = (*DoctorOptions)(nil)
	_ cli.Executable  = (*DoctorOptions)(nil)
)

func (opts *DoctorOptions) Exec(ctx context.Context, c *cli.Config) error {
	requiredNamespaces := []string{
		"istio-system",
		"knative-build",
		"knative-serving",
		"riff-system",
	}
	installationOk, err := opts.checkNamespaces(c, requiredNamespaces)
	if err != nil || !installationOk {
		c.Errorf("\nInstallation is not healthy\n")
		return err
	}

	ns := opts.Namespace
	verbs := []resource.Verb{"get", "list", "create", "update", "delete", "patch", "watch"}
	checks := []resource.AccessChecks{
		{Resource: resource.NewStandardResource(ns, "v1", "core", "configmaps"), Verbs: verbs},
		{Resource: resource.NewStandardResource(ns, "v1", "core", "secrets"), Verbs: verbs},
		{Resource: resource.NewCustomResource(ns, "build.projectriff.io/v1alpha1", "build.projectriff.io", "applications"), Verbs: verbs},
		{Resource: resource.NewCustomResource(ns, "build.projectriff.io/v1alpha1", "build.projectriff.io", "functions"), Verbs: verbs},
		{Resource: resource.NewCustomResource(ns, "request.projectriff.io/v1alpha1", "request.projectriff.io", "handlers"), Verbs: verbs},
		{Resource: resource.NewCustomResource(ns, "stream.projectriff.io/v1alpha1", "stream.projectriff.io", "processors"), Verbs: verbs},
		{Resource: resource.NewCustomResource(ns, "stream.projectriff.io/v1alpha1", "stream.projectriff.io", "streams"), Verbs: verbs},
	}

	existenceSummary, err := checkCustomResourceExistence(c, checks)
	if err != nil {
		c.Errorf("\nAn error occurred while checking for CustomResourceDefinition existence\n")
		c.Errorf("\nInstallation is not healthy\n")
		return err
	}
	existenceSummary.Print(c)
	if !existenceSummary.IsHealthy() {
		c.Errorf("\nInstallation is not healthy\n")
		return nil
	}

	accessSummary, err := checkResourceAccesses(c, checks)
	if err != nil {
		c.Errorf("\nAn error occurred while checking for resource access\n")
		c.Errorf("\nInstallation is not healthy\n")
		return err
	}
	accessSummary.Print(c)
	if !accessSummary.IsHealthy() {
		c.Errorf("\nInstallation is not healthy\n")
	} else {
		c.Successf("\nInstallation is OK\n")
	}

	return nil
}

func NewDoctorCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &DoctorOptions{}

	cmd := &cobra.Command{
		Use:     "doctor",
		Aliases: []string{"doc"},
		Short:   "check " + c.Name + "'s requirements are installed",
		Long: strings.TrimSpace(`
    <todo>
    `),
		Example: "riff doctor",
		Args:    cli.Args(),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.AllNamespacesFlag(cmd, c, &opts.Namespace, &opts.AllNamespaces)

	return cmd
}

func (*DoctorOptions) checkNamespaces(c *cli.Config, requiredNamespaces []string) (bool, error) {
	namespaces, err := c.Core().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	foundNamespaces := sets.NewString()
	for _, namespace := range namespaces.Items {
		foundNamespaces.Insert(namespace.Name)
	}
	printer := printers.GetNewTabWriter(c.Stdout)
	defer printer.Flush()
	ok := true
	for _, namespace := range requiredNamespaces {
		var status string
		if foundNamespaces.Has(namespace) {
			status = cli.Ssuccessf("OK")
		} else {
			ok = false
			status = cli.Serrorf("Missing")
		}
		fmt.Fprintf(printer, "Namespace %q\t%s\n", namespace, status)
	}
	return ok, nil
}

func checkCustomResourceExistence(c *cli.Config, checks []resource.AccessChecks) (*resource.CrdSummary, error) {
	var aggregatedStatuses []resource.CrdStatus
	crds := c.ApiExtensions().CustomResourceDefinitions()
	for _, check := range checks {
		serverResource := check.Resource
		if !serverResource.Custom {
			continue
		}
		crdName := serverResource.CrdName()
		_, err := crds.Get(crdName, metav1.GetOptions{})
		if err == nil {
			aggregatedStatuses = append(aggregatedStatuses, resource.CrdStatus{Resource: serverResource, ExistenceStatus: resource.Exists})
			continue
		}
		if !errors.IsNotFound(err) {
			return nil, err
		}
		aggregatedStatuses = append(aggregatedStatuses, resource.CrdStatus{Resource: serverResource, ExistenceStatus: resource.NotFound})
	}
	return &resource.CrdSummary{Statuses: aggregatedStatuses}, nil
}

func checkResourceAccesses(c *cli.Config, checks []resource.AccessChecks) (*resource.AccessSummary, error) {
	aggregatedStatuses := make([]resource.Status, len(checks))
	for i, check := range checks {
		serverResource := check.Resource
		aggregatedStatus := resource.Status{Resource: serverResource, ReadStatus: resource.AccessUndefined, WriteStatus: resource.AccessUndefined}
		for _, verb := range check.Verbs {
			reviewRequest := serverResource.AsReview(verb)
			result, err := c.Auth().SelfSubjectAccessReviews().Create(reviewRequest)
			if err != nil {
				return nil, err
			}
			evaluationError := result.Status.EvaluationError
			if evaluationError != "" {
				return nil, fmt.Errorf(evaluationError)
			}
			status, err := determineAccessStatus(result)
			if err != nil {
				return nil, err
			}
			if verb.IsRead() {
				aggregatedStatus.ReadStatus = aggregatedStatus.ReadStatus.Combine(status)
			} else {
				aggregatedStatus.WriteStatus = aggregatedStatus.WriteStatus.Combine(status)
			}
		}
		aggregatedStatuses[i] = aggregatedStatus
	}
	return &resource.AccessSummary{Statuses: aggregatedStatuses}, nil
}

func determineAccessStatus(review *authv1.SelfSubjectAccessReview) (*resource.AccessStatus, error) {
	status := review.Status
	if status.Allowed {
		result := resource.Allowed
		return &result, nil
	}
	if status.Denied {
		result := resource.Denied
		return &result, nil
	}
	return nil, fmt.Errorf("unexpected state, review is neither allowed nor denied: %v", review)
}
