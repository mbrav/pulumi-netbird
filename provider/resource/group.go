package resource

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// TEST: InputDiff: false

// Group represents a resource for managing NetBird groups.
type Group struct{}

// Annotate adds a description to the Group resource type.
func (g *Group) Annotate(a infer.Annotator) {
	a.Describe(&g, "A NetBird group, which represents a collection of peers.")
}

// GroupArgs defines input fields for creating or updating a group.
type GroupArgs struct {
	Name      string      `pulumi:"name"`
	Peers     *[]string   `pulumi:"peers,optional"`
	Resources *[]Resource `pulumi:"resources,optional"`
}

// Annotate provides documentation for GroupArgs fields.
func (g *GroupArgs) Annotate(a infer.Annotator) {
	a.Describe(&g.Name, "The name of the NetBird group.")
	a.Describe(&g.Peers, "An optional list of peer IDs to associate with this group.")
	a.Describe(&g.Resources, "An optional list of resources to associate with this group.")
}

// GroupState represents the output state of a group resource.
type GroupState struct {
	Name      string      `pulumi:"name"`
	Peers     *[]string   `pulumi:"peers,optional"`
	Resources *[]Resource `pulumi:"resources,optional"`
}

// Annotate provides documentation for GroupState fields.
func (g *GroupState) Annotate(a infer.Annotator) {
	a.Describe(&g.Name, "The name of the NetBird group.")
	a.Describe(&g.Peers, "An optional list of peer IDs associated with this group.")
	a.Describe(&g.Resources, "An optional list of resources to associate with this group.")
}

// Create creates a new NetBird group.
func (*Group) Create(ctx context.Context, req infer.CreateRequest[GroupArgs]) (infer.CreateResponse[GroupState], error) {
	p.GetLogger(ctx).Debugf("Create:Group name=%s, peers=%v", req.Inputs.Name, req.Inputs.Peers)

	if req.DryRun {
		return infer.CreateResponse[GroupState]{
			ID: "preview",
			Output: GroupState{
				Name:      req.Inputs.Name,
				Peers:     req.Inputs.Peers,
				Resources: req.Inputs.Resources,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[GroupState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	group, err := client.Groups.Create(ctx, nbapi.GroupRequest{
		Name:      req.Inputs.Name,
		Peers:     req.Inputs.Peers,
		Resources: toAPIResourceList(req.Inputs.Resources),
	})
	if err != nil {
		return infer.CreateResponse[GroupState]{}, fmt.Errorf("creating group failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:GroupAPI name=%s, id=%s", group.Name, group.Id)

	peerIDs := make([]string, len(group.Peers))
	for i, peer := range group.Peers {
		peerIDs[i] = peer.Id
	}

	slices.Sort(peerIDs)

	return infer.CreateResponse[GroupState]{
		ID: group.Id,
		Output: GroupState{
			Name:      group.Name,
			Peers:     &peerIDs,
			Resources: fromAPIResourceList(&group.Resources),
		},
	}, nil
}

// Read fetches the current state of a group resource from NetBird.
func (*Group) Read(ctx context.Context, req infer.ReadRequest[GroupArgs, GroupState]) (infer.ReadResponse[GroupArgs, GroupState], error) {
	p.GetLogger(ctx).Debugf("Read:GroupArgs[%s] name=%s", req.ID, req.Inputs.Name)
	p.GetLogger(ctx).Debugf("Read:GroupState[%s] name=%s, id=%s", req.ID, req.State.Name, req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[GroupArgs, GroupState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	group, err := client.Groups.Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[GroupArgs, GroupState]{}, fmt.Errorf("reading group failed: %w", err)
	}

	peerIDs := make([]string, len(group.Peers))
	for i, peer := range group.Peers {
		peerIDs[i] = peer.Id
	}
	// Always sort peer IDs before comparison or use
	slices.Sort(peerIDs)

	return infer.ReadResponse[GroupArgs, GroupState]{
		ID: req.ID,
		Inputs: GroupArgs{
			Name:      group.Name,
			Peers:     &peerIDs,
			Resources: req.Inputs.Resources,
		},
		State: GroupState{
			Name:      group.Name,
			Peers:     &peerIDs,
			Resources: fromAPIResourceList(&group.Resources),
		},
	}, nil
}

// Update updates the state of the group if needed.
func (*Group) Update(ctx context.Context, req infer.UpdateRequest[GroupArgs, GroupState]) (infer.UpdateResponse[GroupState], error) {
	p.GetLogger(ctx).Debugf("Update:Group[%s] name=%s", req.ID, req.Inputs.Name)

	if req.DryRun {
		return infer.UpdateResponse[GroupState]{
			Output: GroupState{
				Name:      req.Inputs.Name,
				Peers:     req.Inputs.Peers,
				Resources: req.Inputs.Resources,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[GroupState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.Groups.Update(ctx, req.ID, nbapi.GroupRequest{
		Name:      req.Inputs.Name,
		Peers:     req.Inputs.Peers,
		Resources: toAPIResourceList(req.Inputs.Resources),
	})
	if err != nil {
		return infer.UpdateResponse[GroupState]{}, fmt.Errorf("updating group failed: %w", err)
	}

	peerIDs := make([]string, len(updated.Peers))
	for i, peer := range updated.Peers {
		peerIDs[i] = peer.Id
	}

	return infer.UpdateResponse[GroupState]{
		Output: GroupState{
			Name:      updated.Name,
			Peers:     &peerIDs,
			Resources: fromAPIResourceList(&updated.Resources),
		},
	}, nil
}

// Delete removes a group from NetBird.
func (*Group) Delete(ctx context.Context, req infer.DeleteRequest[GroupState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:Group[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Groups.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting group failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*Group) Diff(ctx context.Context, req infer.DiffRequest[GroupArgs, GroupState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Group[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	// Name is reflected in state — normal comparison
	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	// Peers: compare if both are non-nil
	if req.Inputs.Peers != nil && req.State.Peers != nil {
		inPeers := slices.Clone(*req.Inputs.Peers)
		stPeers := slices.Clone(*req.State.Peers)
		slices.Sort(inPeers)
		slices.Sort(stPeers)

		if !slices.Equal(inPeers, stPeers) {
			diff["peers"] = p.PropertyDiff{
				InputDiff: false,
				Kind:      p.Update,
			}
		}
	} else if (req.Inputs.Peers == nil) != (req.State.Peers == nil) {
		diff["peers"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	// Resources: input-only field — always tracked as an input diff
	if req.Inputs.Resources != nil {
		diff["resources"] = p.PropertyDiff{
			InputDiff: true,
			Kind:      p.Update,
		}
	}

	p.GetLogger(ctx).Debugf("Diff:Group[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*Group) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[GroupArgs], error) {
	p.GetLogger(ctx).Debugf("Check:Group old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())
	args, failures, err := infer.DefaultCheck[GroupArgs](ctx, req.NewInputs)

	return infer.CheckResponse[GroupArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*Group) WireDependencies(f infer.FieldSelector, args *GroupArgs, state *GroupState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.Peers).DependsOn(f.InputField(&args.Peers))
}
