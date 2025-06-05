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

// NetworkResource represents a Pulumi resource for NetBird network resources.
type NetworkResource struct{}

// Annotate adds a description annotation for the NetworkResource type for generated SDKs.
func (net *NetworkResource) Annotate(a infer.Annotator) {
	a.Describe(net, "A NetBird network resource, such as a CIDR range assigned to the network.")
}

// NetworkResourceArgs represents the input arguments for creating or updating a network resource.
type NetworkResourceArgs struct {
	Name        string   `pulumi:"name"`
	Description *string  `pulumi:"description,optional"`
	NetworkID   string   `pulumi:"network_id"`
	Address     string   `pulumi:"address"`
	Enabled     bool     `pulumi:"enabled"`
	GroupIDs    []string `pulumi:"group_ids"`
}

// Annotate provides documentation for NetworkResourceArgs fields.
func (a *NetworkResourceArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&a.Name, "Name of the network resource.")
	annotator.Describe(&a.Description, "Optional description of the resource.")
	annotator.Describe(&a.NetworkID, "ID of the network this resource belongs to.")
	annotator.Describe(&a.Address, "CIDR or IP address block assigned to the resource.")
	annotator.Describe(&a.Enabled, "Whether the resource is enabled.")
	annotator.Describe(&a.GroupIDs, "List of group IDs associated with this resource.")
}

// NetworkResourceState represents the state of a network resource.
type NetworkResourceState struct {
	Name        string   `pulumi:"name"`
	Description *string  `pulumi:"description,optional"`
	NetworkID   string   `pulumi:"network_id"`
	Address     string   `pulumi:"address"`
	Enabled     bool     `pulumi:"enabled"`
	GroupIDs    []string `pulumi:"group_ids"`
}

// Annotate provides documentation for NetworkResourceState fields.
func (s *NetworkResourceState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&s.Name, "Name of the network resource.")
	annotator.Describe(&s.Description, "Optional description of the resource.")
	annotator.Describe(&s.NetworkID, "ID of the network this resource belongs to.")
	annotator.Describe(&s.Address, "CIDR or IP address block assigned to the resource.")
	annotator.Describe(&s.Enabled, "Whether the resource is enabled.")
	annotator.Describe(&s.GroupIDs, "List of group IDs associated with this resource.")
}

// Create creates a new NetBird network resource.
func (*NetworkResource) Create(ctx context.Context, req infer.CreateRequest[NetworkResourceArgs]) (infer.CreateResponse[NetworkResourceState], error) {
	p.GetLogger(ctx).Debugf("Create:NetworkResource name=%s, description=%s net_id=%s", req.Inputs.Name, strPtr(req.Inputs.Description), req.Inputs.NetworkID)

	// Always sort Input group IDS
	slices.Sort(req.Inputs.GroupIDs)

	if req.DryRun {
		return infer.CreateResponse[NetworkResourceState]{
			ID: "preview",
			Output: NetworkResourceState{
				Name:        req.Inputs.Name,
				Description: req.Inputs.Description,
				NetworkID:   req.Inputs.NetworkID,
				Address:     req.Inputs.Address,
				Enabled:     req.Inputs.Enabled,
				GroupIDs:    req.Inputs.GroupIDs,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[NetworkResourceState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	net, err := client.Networks.Resources(req.Inputs.NetworkID).Create(ctx, nbapi.NetworkResourceRequest{
		Name:        req.Inputs.Name,
		Description: req.Inputs.Description,
		Address:     req.Inputs.Address,
		Enabled:     req.Inputs.Enabled,
		Groups:      req.Inputs.GroupIDs,
	})
	if err != nil {
		return infer.CreateResponse[NetworkResourceState]{}, fmt.Errorf("creating network failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:NetworkResourceAPI name=%s, id=%s, net_id=%s", net.Name, net.Id, req.Inputs.NetworkID)

	return infer.CreateResponse[NetworkResourceState]{
		ID: net.Id,
		Output: NetworkResourceState{
			Name:        net.Name,
			Description: net.Description,
			NetworkID:   req.Inputs.NetworkID,
			Address:     net.Address,
			Enabled:     net.Enabled,
			GroupIDs:    getNetworkResourceGroupIDs(net),
		},
	}, nil
}

// Read fetches the current state of a network resource from NetBird.
func (*NetworkResource) Read(ctx context.Context, req infer.ReadRequest[NetworkResourceArgs, NetworkResourceState]) (infer.ReadResponse[NetworkResourceArgs, NetworkResourceState], error) {
	p.GetLogger(ctx).Debugf("Read:NetworkReourceArgs[%s] name=%s", req.ID, req.Inputs.Name)
	p.GetLogger(ctx).Debugf("Read:NetworkResourceState[%s] name=%s, netd_id=%s", req.ID, req.State.Name, req.State.NetworkID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[NetworkResourceArgs, NetworkResourceState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	net, err := client.Networks.Resources(req.State.NetworkID).Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[NetworkResourceArgs, NetworkResourceState]{}, fmt.Errorf("reading network failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Read:NetworkResourceAPI[%s] name=%s", net.Id, net.Name)

	return infer.ReadResponse[NetworkResourceArgs, NetworkResourceState]{
		ID: req.ID,
		Inputs: NetworkResourceArgs{
			Name:        net.Name,
			Description: net.Description,
			NetworkID:   req.State.NetworkID,
			Address:     net.Address,
			Enabled:     net.Enabled,
			GroupIDs:    getNetworkResourceGroupIDs(net),
		},
		State: NetworkResourceState{
			Name:        net.Name,
			Description: net.Description,
			NetworkID:   req.State.NetworkID,
			Address:     net.Address,
			Enabled:     net.Enabled,
			GroupIDs:    getNetworkResourceGroupIDs(net),
		},
	}, nil
}

// Update updates the state of the NetBird network resource if needed.
func (*NetworkResource) Update(ctx context.Context, req infer.UpdateRequest[NetworkResourceArgs, NetworkResourceState]) (infer.UpdateResponse[NetworkResourceState], error) {
	p.GetLogger(ctx).Debugf("Update:NetworkResource[%s] name=%s", req.ID, req.Inputs.Name)

	// Always sort group IDs before comparison or use
	slices.Sort(req.Inputs.GroupIDs)

	if req.DryRun {
		return infer.UpdateResponse[NetworkResourceState]{
			Output: NetworkResourceState{
				Name:        req.Inputs.Name,
				Description: req.Inputs.Description,
				NetworkID:   req.Inputs.NetworkID,
				Address:     req.Inputs.Address,
				Enabled:     req.Inputs.Enabled,
				GroupIDs:    req.Inputs.GroupIDs,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[NetworkResourceState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	net, err := client.Networks.Resources(req.Inputs.NetworkID).Update(ctx, req.ID, nbapi.NetworkResourceRequest{
		Name:        req.Inputs.Name,
		Description: req.Inputs.Description,
		Address:     req.Inputs.Address,
		Enabled:     req.Inputs.Enabled,
		Groups:      req.Inputs.GroupIDs,
	})
	if err != nil {
		return infer.UpdateResponse[NetworkResourceState]{}, fmt.Errorf("updating network resource failed: %w", err)
	}

	return infer.UpdateResponse[NetworkResourceState]{
		Output: NetworkResourceState{
			Name:        net.Name,
			Description: net.Description,
			NetworkID:   req.Inputs.NetworkID,
			Address:     net.Address,
			Enabled:     net.Enabled,
			GroupIDs:    getNetworkResourceGroupIDs(net),
		},
	}, nil
}

// Delete removes a network resource from NetBird.
func (*NetworkResource) Delete(ctx context.Context, req infer.DeleteRequest[NetworkResourceState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:NetworkResource[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Networks.Resources(req.State.NetworkID).Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting network resource failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*NetworkResource) Diff(ctx context.Context, req infer.DiffRequest[NetworkResourceArgs, NetworkResourceState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:NetworkResource[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if !equalPtr(req.Inputs.Description, req.State.Description) {
		diff["description"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Address != req.State.Address {
		diff["address"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.GroupIDs != nil && req.State.GroupIDs != nil {
		slices.Sort(req.Inputs.GroupIDs)
		slices.Sort(req.State.GroupIDs)

		if !slices.Equal(req.Inputs.GroupIDs, req.State.GroupIDs) {
			diff["group_ids"] = p.PropertyDiff{
				InputDiff: false,
				Kind:      p.Update,
			}

			p.GetLogger(ctx).Debugf("Diff:NetworkResource group_ids input=%s output=%s", req.Inputs.GroupIDs, req.State.GroupIDs)
		}
	}

	p.GetLogger(ctx).Debugf("Diff:NetworkResource[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*NetworkResource) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[NetworkResourceArgs], error) {
	p.GetLogger(ctx).Debugf("Check:NetworkResource old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())
	args, failures, err := infer.DefaultCheck[NetworkResourceArgs](ctx, req.NewInputs)

	return infer.CheckResponse[NetworkResourceArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*NetworkResource) WireDependencies(field infer.FieldSelector, args *NetworkResourceArgs, state *NetworkResourceState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Description).DependsOn(field.InputField(&args.Description))
	field.OutputField(&state.NetworkID).DependsOn(field.InputField(&args.NetworkID))
	field.OutputField(&state.Address).DependsOn(field.InputField(&args.Address))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.GroupIDs).DependsOn(field.InputField(&args.GroupIDs))
}

// Extract and sort group IDs.
func getNetworkResourceGroupIDs(net *nbapi.NetworkResource) []string {
	groupIDs := make([]string, 0, len(net.Groups))
	for _, g := range net.Groups {
		groupIDs = append(groupIDs, g.Id)
	}

	slices.Sort(groupIDs)

	return groupIDs
}
