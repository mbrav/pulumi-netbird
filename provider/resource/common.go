// Package resource provides the NetBird resource types
package resource

import (
	"slices"
	"strings"

	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
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

// equalResourcesPtr compares two *[]Resource slices by value, treating nil and empty as equal.
func equalResourcesPtr(resourcesA, resourcesB *[]Resource) bool {
	aLen := 0
	if resourcesA != nil {
		aLen = len(*resourcesA)
	}

	bLen := 0
	if resourcesB != nil {
		bLen = len(*resourcesB)
	}

	if aLen == 0 && bLen == 0 {
		return true
	}

	if resourcesA == nil || resourcesB == nil {
		return false
	}

	if aLen != bLen {
		return false
	}

	aSorted := sortedResources(*resourcesA)
	bSorted := sortedResources(*resourcesB)

	for i := range aSorted {
		if !equalResourcePtr(&aSorted[i], &bSorted[i]) {
			return false
		}
	}

	return true
}

func sortedResources(resources []Resource) []Resource {
	sorted := slices.Clone(resources)
	slices.SortFunc(sorted, compareResources)

	return sorted
}

func compareResources(resA, resB Resource) int {
	if typeCompare := strings.Compare(string(resA.Type), string(resB.Type)); typeCompare != 0 {
		return typeCompare
	}

	return strings.Compare(resA.ID, resB.ID)
}

func equalResourcePtr(resourceA, resourceB *Resource) bool {
	if resourceA == nil && resourceB == nil {
		return true
	}

	if resourceA == nil || resourceB == nil {
		return false
	}

	return resourceA.Type == resourceB.Type && resourceA.ID == resourceB.ID
}
