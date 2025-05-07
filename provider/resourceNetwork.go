package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// Network represents a resource for managing NetBird networks.
type Network struct{}

// Annotate adds a description to the Network resource type.
func (n *Network) Annotate(a infer.Annotator) {
	a.Describe(&n, "A NetBird network.")
}

// NetworkArgs defines input fields for creating or updating a network.
type NetworkArgs struct {
	Name        string  `pulumi:"name"`
	Description *string `pulumi:"description,optional"`
}

// Annotate provides documentation for NetworkArgs fields.
func (n *NetworkArgs) Annotate(a infer.Annotator) {
	a.Describe(&n.Name, "The name of the NetBird network.")
	a.Describe(&n.Description, "An optional description of the network.")
}

// NetworkState represents the output state of a network resource.
type NetworkState struct {
	Name        string  `pulumi:"name"`
	Description *string `pulumi:"description,optional"`
	NbID        string  `pulumi:"nbId"`
}

// Annotate provides documentation for NetworkState fields.
func (n *NetworkState) Annotate(a infer.Annotator) {
	a.Describe(&n.Name, "The name of the NetBird network.")
	a.Describe(&n.Description, "An optional description of the network.")
	a.Describe(&n.NbID, "The internal NetBird network ID.")
}

// Create creates a new NetBird network.
func (*Network) Create(ctx context.Context, req infer.CreateRequest[NetworkArgs]) (infer.CreateResponse[NetworkState], error) {
	p.GetLogger(ctx).Debugf("Create:NetworkState name=%s, description=%s", req.Inputs.Name, strPtr(req.Inputs.Description))
	// Return preview state without performing real actions.
	if req.DryRun {
		return infer.CreateResponse[NetworkState]{
			Output: NetworkState{
				Name:        req.Inputs.Name,
				Description: req.Inputs.Description,
				NbID:        "unknown",
			},
		}, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[NetworkState]{}, err
	}

	net, err := client.Networks.Create(ctx, nbapi.NetworkRequest{
		Name:        req.Inputs.Name,
		Description: req.Inputs.Description,
	})
	if err != nil {
		return infer.CreateResponse[NetworkState]{}, fmt.Errorf("creating network failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:NetworkAPI name=%s, description=%s id=%s", net.Name, strPtr(net.Description), net.Id)

	// Return the created network's ID and state
	return infer.CreateResponse[NetworkState]{
		ID: net.Id,
		Output: NetworkState{
			Name:        net.Name,
			Description: net.Description,
			NbID:        net.Id,
		},
	}, nil
}

// Read fetches the current state of a network resource from NetBird.
func (*Network) Read(ctx context.Context, req infer.ReadRequest[NetworkArgs, NetworkState]) (infer.ReadResponse[NetworkArgs, NetworkState], error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState]{}, err
	}

	net, err := client.Networks.Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState]{}, fmt.Errorf("reading network failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Read:NetworkAPI name=%s, description=%s id=%s", net.Name, strPtr(net.Description), net.Id)

	return infer.ReadResponse[NetworkArgs, NetworkState]{
		ID: req.ID,
		Inputs: NetworkArgs{
			Name:        net.Name,
			Description: net.Description,
		},
		State: NetworkState{
			Name:        net.Name,
			Description: net.Description,
			NbID:        req.ID,
		},
	}, nil
}

// Diff compares the desired and actual state of a network resource.
// func (*Network) Diff(ctx context.Context, req infer.DiffRequest[NetworkArgs, NetworkState]) (infer.DiffResponse, error) {
// 	diff := map[string]p.PropertyDiff{}
//
// 	// Compare Name input vs state
// 	if req.Inputs.Name != req.State.Name {
// 		diff["name"] = p.PropertyDiff{Kind: p.Update}
// 	}
//
// 	// Compare Description input vs state (handling nils)
// 	oldDesc := ""
// 	if req.State.Description != nil {
// 		oldDesc = *req.State.Description
// 	}
// 	newDesc := ""
// 	if req.Inputs.Description != nil {
// 		newDesc = *req.Inputs.Description
// 	}
// 	if oldDesc != newDesc {
// 		diff["description"] = p.PropertyDiff{Kind: p.Update}
// 	}
//
// 	p.GetLogger(ctx).Debugf("Read:Network description old:%s new:%s", newDesc, oldDesc)
//
// 	return infer.DiffResponse{
// 		HasChanges:   len(diff) > 0,
// 		DetailedDiff: diff,
// 	}, nil
// }

func (*Network) Diff(ctx context.Context, req infer.DiffRequest[NetworkArgs, NetworkState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Network state=%+v, input=%+v", req.State, req.Inputs)
	return infer.DiffResponse{}, nil
}

// Update modifies an existing network resource.
func (*Network) Update(ctx context.Context, req infer.UpdateRequest[NetworkArgs, NetworkState]) (infer.UpdateResponse[NetworkState], error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[NetworkState]{}, err
	}

	net, err := client.Networks.Update(ctx, req.ID, nbapi.NetworkRequest{
		Name:        req.Inputs.Name,
		Description: req.Inputs.Description,
	})
	if err != nil {
		return infer.UpdateResponse[NetworkState]{}, fmt.Errorf("updating network failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Update:NetworkAPI name=%s, description=%s id=%s", net.Name, strPtr(net.Description), net.Id)

	return infer.UpdateResponse[NetworkState]{
		Output: NetworkState{
			Name:        net.Name,
			Description: net.Description,
		},
	}, nil
}

// Delete removes a network resource from NetBird.
func (*Network) Delete(ctx context.Context, req infer.DeleteRequest[NetworkState]) (infer.DeleteResponse, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	p.GetLogger(ctx).Debugf("Delete:NetworkState name=%s, description=%s id=%s", req.State.Name, strPtr(req.State.Description), req.State.NbID)

	if err := client.Networks.Delete(ctx, req.ID); err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting network failed: %w", err)
	}
	return infer.DeleteResponse{}, nil
}

// WireDependencies links output fields to their corresponding input fields for proper dependency tracking.
func (*Network) WireDependencies(f infer.FieldSelector, args *NetworkArgs, state *NetworkState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.Description).DependsOn(f.InputField(&args.Description))
}
