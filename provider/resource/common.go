package resource

import (
	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Resource represents a single NetBird resource used in a rule (e.g., domain, host, subnet).
type Resource struct {
	Id   string       `pulumi:"id"`   // The unique ID of the resource
	Type ResourceType `pulumi:"type"` // The type of the resource (domain, host, subnet)
}

// Annotation for Resource for generated SDKs.
func (r *Resource) Annotate(a infer.Annotator) {
	a.Describe(&r.Id, "The unique identifier of the resource.")
	a.Describe(&r.Type, "The type of resource: 'domain', 'host', or 'subnet'.")
}

// ResourceType defines the allowed resource types for a policy rule.
type ResourceType string

// Enum constants for resource types.
const (
	ResourceTypeDomain ResourceType = ResourceType(nbapi.ResourceTypeDomain)
	ResourceTypeHost   ResourceType = ResourceType(nbapi.ResourceTypeHost)
	ResourceTypeSubnet ResourceType = ResourceType(nbapi.ResourceTypeSubnet)
)

// Values returns the list of supported ResourceType values for Pulumi enum generation.
func (ResourceType) Values() []infer.EnumValue[ResourceType] {
	return []infer.EnumValue[ResourceType]{
		{Name: "Domain", Value: ResourceTypeDomain, Description: "A domain resource (e.g., example.com)."},
		{Name: "Host", Value: ResourceTypeHost, Description: "A host resource (e.g., peer or device)."},
		{Name: "Subnet", Value: ResourceTypeSubnet, Description: "A subnet resource (e.g., 192.168.0.0/24)."},
	}
}

// Converts a single Resource to nbapi.Resource.
func toAPIResource(in *Resource) *nbapi.Resource {
	if in == nil {
		return nil
	}

	return &nbapi.Resource{
		Id:   in.Id,
		Type: nbapi.ResourceType(in.Type),
	}
}

// Converts a single nbapi.Resource to Resource.
func fromAPIResource(in *nbapi.Resource) *Resource {
	if in == nil {
		return nil
	}

	return &Resource{
		Id:   in.Id,
		Type: ResourceType(in.Type),
	}
}

// Refactored equalResourcePtr to explicitly check for both pointers being nil
// before comparing their fields. This ensures correct equality checks and
// prevents potential nil pointer dereference issues.
func equalResourcePtr(a, b *Resource) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a.Type == b.Type && a.Id == b.Id
}
