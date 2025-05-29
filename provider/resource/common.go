// Package resource provides the NetBird resource types
package resource

import (
	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Resource represents a single NetBird resource used in a rule (e.g., domain, host, subnet).
type Resource struct {
	ID   string `pulumi:"id"`   // The unique ID of the resource
	Type Type   `pulumi:"type"` // The type of the resource (domain, host, subnet)
}

// Annotate adds descriptive annotations to the Resource fields for use in generated SDKs.
func (r *Resource) Annotate(a infer.Annotator) {
	a.Describe(&r.ID, "The unique identifier of the resource.")
	a.Describe(&r.Type, "The type of resource: 'domain', 'host', or 'subnet'.")
}

// Type defines the allowed resource types for a policy rule.
type Type string

// ResourceTypeDomain, ResourceTypeHost, and ResourceTypeSyyubnet represent different types of network resources.
const (
	ResourceTypeDomain Type = Type(nbapi.ResourceTypeDomain)
	ResourceTypeHost   Type = Type(nbapi.ResourceTypeHost)
	ResourceTypeSubnet Type = Type(nbapi.ResourceTypeSubnet)
)

// Values returns the list of supported ResourceType values for Pulumi enum generation.
func (Type) Values() []infer.EnumValue[Type] {
	return []infer.EnumValue[Type]{
		{Name: "Domain", Value: ResourceTypeDomain, Description: "A domain resource (e.g., example.com)."},
		{Name: "Host", Value: ResourceTypeHost, Description: "A host resource (e.g., peer or device)."},
		{Name: "Subnet", Value: ResourceTypeSubnet, Description: "A subnet resource (e.g., 192.168.0.0/24)."},
	}
}

// Converts a single Resource to nbapi.Resource.
func toAPIResource(resource *Resource) *nbapi.Resource {
	if resource == nil {
		return nil
	}

	return &nbapi.Resource{
		Id:   resource.ID,
		Type: nbapi.ResourceType(resource.Type),
	}
}

// Converts a slice of *Resource to a pointer to a slice of nbapi.Resource.
// Returns nil if the input is nil.
func toAPIResourceList(resources *[]Resource) *[]nbapi.Resource {
	if resources == nil {
		return nil
	}

	converted := make([]nbapi.Resource, len(*resources))
	for i, r := range *resources {
		converted[i] = *toAPIResource(&r)
	}

	return &converted
}

// Converts a single nbapi.Resource to Resource.
func fromAPIResource(apiResource *nbapi.Resource) *Resource {
	if apiResource == nil {
		return nil
	}

	return &Resource{
		ID:   apiResource.Id,
		Type: Type(apiResource.Type),
	}
}

// Converts a slice of nbapi.Resource to a pointer to a slice of Resource.
// Returns nil if the input is nil.
func fromAPIResourceList(apiResources *[]nbapi.Resource) *[]Resource {
	if apiResources == nil {
		return nil
	}

	converted := make([]Resource, len(*apiResources))
	for i, r := range *apiResources {
		converted[i] = *fromAPIResource(&r)
	}

	return &converted
}

// Refactored equalResourcePtr to explicitly check for both pointers being nil
// before comparing their fields. This ensures correct equality checks and
// prevents potential nil pointer dereference issues.
func equalResourcePtr(resourceA, resourceB *Resource) bool {
	if resourceA == nil && resourceB == nil {
		return true
	}

	if resourceA == nil || resourceB == nil {
		return false
	}

	return resourceA.Type == resourceB.Type && resourceA.ID == resourceB.ID
}
