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

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/cli/printers"
	"github.com/projectriff/cli/pkg/doctor"
	"github.com/spf13/cobra"
	authv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type DoctorOptions struct {
	Namespace string
}

var (
	_ cli.Validatable = (*DoctorOptions)(nil)
	_ cli.Executable  = (*DoctorOptions)(nil)
)

func (opts *DoctorOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := cli.EmptyFieldError

	if opts.Namespace == "" {
		errs = errs.Also(cli.ErrMissingField(cli.NamespaceFlagName))
	}

	return errs
}

func (opts *DoctorOptions) Exec(ctx context.Context, c *cli.Config) error {
	requiredNamespaces := []string{
		"istio-system",
		"knative-build",
		"knative-serving",
		"riff-system",
	}
	installationOk, err := opts.checkNamespaces(c, requiredNamespaces)
	if err != nil || !installationOk {
		c.Printf("\n")
		c.Errorf("Installation is not healthy\n")
		return err
	}

	ns := opts.Namespace
	verbs := []doctor.Verb{"get", "list", "create", "update", "delete", "patch", "watch"}
	checks := []doctor.AccessChecks{
		{Resource: doctor.NewStandardResource(ns, "v1", "core", "configmaps"), Verbs: verbs},
		{Resource: doctor.NewStandardResource(ns, "v1", "core", "secrets"), Verbs: verbs},
		{Resource: doctor.NewCustomResource(ns, "build.projectriff.io/v1alpha1", "build.projectriff.io", "applications"), Verbs: verbs},
		{Resource: doctor.NewCustomResource(ns, "build.projectriff.io/v1alpha1", "build.projectriff.io", "functions"), Verbs: verbs},
		{Resource: doctor.NewCustomResource(ns, "request.projectriff.io/v1alpha1", "request.projectriff.io", "handlers"), Verbs: verbs},
		{Resource: doctor.NewCustomResource(ns, "stream.projectriff.io/v1alpha1", "stream.projectriff.io", "processors"), Verbs: verbs},
		{Resource: doctor.NewCustomResource(ns, "stream.projectriff.io/v1alpha1", "stream.projectriff.io", "streams"), Verbs: verbs},
	}

	accessSummary, err := opts.checkResourceAccesses(c, checks)
	if err != nil {
		c.Printf("\n")
		c.Errorf("An error occurred while checking for resource access\n")
		c.Printf("\n")
		c.Errorf("Installation is not healthy\n")
		return err
	}
	c.Printf("\n")
	accessSummary.Print(c)
	c.Printf("\n")
	if !accessSummary.IsHealthy() {
		c.Errorf("Installation is not healthy\n")
	} else {
		c.Successf("Installation is OK\n")
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
Check that ` + c.Name + ` is installed.

The doctor checks that necessary system components are installed and the user
has access to resources in the namespace.

The doctor is not a tool for monitoring the health of the cluster.
`),
		Example: "riff doctor",
		Args:    cli.Args(),
		PreRunE: cli.ValidateOptions(ctx, opts),
		RunE:    cli.ExecOptions(ctx, c, opts),
	}

	cli.NamespaceFlag(cmd, c, &opts.Namespace)

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

func (opts *DoctorOptions) checkResourceAccesses(c *cli.Config, checks []doctor.AccessChecks) (*doctor.AccessSummary, error) {
	crds := c.ApiExtensions().CustomResourceDefinitions()
	aggregatedStatuses := make([]doctor.Status, len(checks))
	for i, check := range checks {
		serverResource := check.Resource
		aggregatedStatus := doctor.Status{Resource: serverResource, ReadStatus: doctor.AccessUndefined, WriteStatus: doctor.AccessUndefined}
		if serverResource.Custom {
			missing, err := opts.isCustomResourceMissing(crds, serverResource.CrdName())
			if err != nil {
				return nil, err
			}
			if missing {
				aggregatedStatus.ReadStatus = doctor.Missing
				aggregatedStatus.WriteStatus = doctor.Missing
				aggregatedStatuses[i] = aggregatedStatus
				continue
			}
		}
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
			status, err := opts.determineAccessStatus(result)
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
	return &doctor.AccessSummary{Statuses: aggregatedStatuses}, nil
}

func (*DoctorOptions) determineAccessStatus(review *authv1.SelfSubjectAccessReview) (*doctor.AccessStatus, error) {
	status := review.Status
	if status.Allowed {
		result := doctor.Allowed
		return &result, nil
	}
	if status.Denied {
		result := doctor.Denied
		return &result, nil
	}
	return nil, fmt.Errorf("unexpected state, review is neither allowed nor denied: %v", review)
}

func (*DoctorOptions) isCustomResourceMissing(crds v1beta1.CustomResourceDefinitionInterface, crdName string) (bool, error) {
	_, err := crds.Get(crdName, metav1.GetOptions{})
	if err == nil {
		return false, nil
	}
	if errors.IsNotFound(err) {
		return true, nil
	}
	return false, err
}
