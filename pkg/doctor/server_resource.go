package doctor

import (
	"fmt"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/cli/printers"
	authv1 "k8s.io/api/authorization/v1"
)

type ServerResource struct {
	namespace string
	apiGroup  string
	group     string
	kind      string
	Custom    bool
}

func NewStandardResource(namespace string, apiGroup string, group string, kind string) ServerResource {
	return ServerResource{
		namespace: namespace,
		apiGroup:  apiGroup,
		group:     group,
		kind:      kind,
		Custom:    false,
	}
}

func NewCustomResource(namespace string, apiGroup string, group string, kind string) ServerResource {
	return ServerResource{
		namespace: namespace,
		apiGroup:  apiGroup,
		group:     group,
		kind:      kind,
		Custom:    true,
	}
}

type StringSet map[string]struct{}

func (ss StringSet) Contains(value string) bool {
	_, ok := ss[value]
	return ok
}

func (resource *ServerResource) AsReview(verb Verb) *authv1.SelfSubjectAccessReview {
	return &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Namespace: resource.namespace,
				Group:     resource.group,
				Verb:      verb.String(),
				Resource:  resource.kind,
			},
		},
	}
}

func (resource *ServerResource) NamespaceString() string {
	ns := resource.namespace
	if ns != "" {
		return ns
	}
	return "*"
}

func (resource *ServerResource) CrdName() string {
	return fmt.Sprintf("%s.%s", resource.kind, resource.group)
}

type AccessChecks struct {
	Resource ServerResource
	Verbs    []Verb
}

type Verb string

func (v *Verb) String() string {
	return string(*v)
}

func (v *Verb) IsRead() bool {
	verb := *v
	return verb == "get" || verb == "list" || verb == "watch"
}

func (v *Verb) IsWrite() bool {
	verb := *v
	return verb == "create" || verb == "update" || verb == "delete" || verb == "patch"
}

type AccessSummary struct {
	Statuses []Status
}

func (as *AccessSummary) IsHealthy() bool {
	for _, status := range as.Statuses {
		if status.ReadStatus != Allowed || status.WriteStatus != Allowed {
			return false
		}
	}
	return true
}

func (as *AccessSummary) Print(c *cli.Config) {
	printer := printers.GetNewTabWriter(c.Stdout)
	defer printer.Flush()
	fmt.Fprintf(printer, "\nNAMESPACE\tGROUP\tRESOURCE\tREAD STATUS\tWRITE STATUS\n")
	for _, status := range as.Statuses {
		resource := status.Resource
		fmt.Fprintf(printer, "%s\t%s\t%s\t%s\t%s\t\n",
			resource.NamespaceString(),
			resource.group,
			resource.kind,
			status.ReadStatus.String(),
			status.WriteStatus.String())
	}
}

type Status struct {
	Resource    ServerResource
	ReadStatus  AccessStatus
	WriteStatus AccessStatus
}

type AccessStatus int

const (
	AccessUndefined AccessStatus = iota
	Allowed                      /* right is granted */
	Denied                       /* right is denied */
	Mixed                        /* for the same resource, some rights are granted, some are denied */
	Missing                      /* resource not deployed */
)

func (as *AccessStatus) Combine(new *AccessStatus) AccessStatus {
	if *as == AccessUndefined {
		return *new
	}
	if *as != *new {
		return Mixed
	}
	if *as == Allowed {
		return Allowed
	}
	return Denied
}

func (as *AccessStatus) String() string {
	status := *as
	switch status {
	case Allowed:
		return cli.Ssuccessf("OK")
	case Mixed:
		return cli.Swarnf("MIXED")
	case Denied:
		return cli.Serrorf("KO")
	case Missing:
		return cli.Serrorf("MISSING RESOURCE")
	}
	panic(fmt.Sprintf("Unsupported value %v", status))
}
