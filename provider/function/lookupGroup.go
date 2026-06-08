package function

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// LookupGroup looks up an existing NetBird group by name.
type LookupGroup struct{}

// Annotate describes the function.
func (f *LookupGroup) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up an existing NetBird group by name and return its ID, peer list, and resource list.")
}

// LookupGroupArgs are the inputs for LookupGroup.
type LookupGroupArgs struct {
	Name string `pulumi:"name"`
}

// Annotate provides field descriptions for LookupGroupArgs.
func (a *LookupGroupArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "The name of the group to look up.")
}

// ResourceRef is a reference to a NetBird resource (domain, host, or subnet).
type ResourceRef struct {
	ID   string `pulumi:"id"`
	Type string `pulumi:"type"`
}

// Annotate provides field descriptions for ResourceRef.
func (r *ResourceRef) Annotate(ann infer.Annotator) {
	ann.Describe(&r.ID, "The unique identifier of the resource.")
	ann.Describe(&r.Type, "The type of resource: 'domain', 'host', or 'subnet'.")
}

// LookupGroupResult is the output of LookupGroup.
type LookupGroupResult struct {
	ID             string        `pulumi:"groupId"`
	Name           string        `pulumi:"name"`
	PeersCount     int           `pulumi:"peersCount"`
	ResourcesCount int           `pulumi:"resourcesCount"`
	Peers          []string      `pulumi:"peers"`
	Resources      []ResourceRef `pulumi:"resources"`
}

// Annotate provides field descriptions for LookupGroupResult.
func (r *LookupGroupResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.ID, "The NetBird group ID.")
	ann.Describe(&r.Name, "The group name.")
	ann.Describe(&r.PeersCount, "Number of peers in the group.")
	ann.Describe(&r.ResourcesCount, "Number of resources in the group.")
	ann.Describe(&r.Peers, "IDs of peers belonging to the group.")
	ann.Describe(&r.Resources, "Resources associated with the group.")
}

// Invoke looks up a group by name.
func (f *LookupGroup) Invoke(ctx context.Context, req infer.FunctionRequest[LookupGroupArgs]) (infer.FunctionResponse[LookupGroupResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupGroupResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	groups, err := client.Groups.List(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupGroupResult]{}, fmt.Errorf("listing groups failed: %w", err)
	}

	for _, group := range groups {
		if group.Name != req.Input.Name {
			continue
		}

		peers := make([]string, len(group.Peers))
		for i, p := range group.Peers {
			peers[i] = p.Id
		}

		resources := make([]ResourceRef, len(group.Resources))
		for i, r := range group.Resources {
			resources[i] = ResourceRef{
				ID:   r.Id,
				Type: string(r.Type),
			}
		}

		return infer.FunctionResponse[LookupGroupResult]{
			Output: LookupGroupResult{
				ID:             group.Id,
				Name:           group.Name,
				PeersCount:     group.PeersCount,
				ResourcesCount: group.ResourcesCount,
				Peers:          peers,
				Resources:      resources,
			},
		}, nil
	}

	return infer.FunctionResponse[LookupGroupResult]{}, fmt.Errorf("group %q not found", req.Input.Name)
}
