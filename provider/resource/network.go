package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
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
}

// Annotate provides documentation for NetworkState fields.
func (n *NetworkState) Annotate(a infer.Annotator) {
	a.Describe(&n.Name, "The name of the NetBird network.")
	a.Describe(&n.Description, "An optional description of the network.")
}

// Create creates a new NetBird network.
func (*Network) Create(ctx context.Context, req infer.CreateRequest[NetworkArgs]) (infer.CreateResponse[NetworkState], error) {
	p.GetLogger(ctx).Debugf("Create:Network name=%s, description=%s", req.Inputs.Name, strPtr(req.Inputs.Description))

	if req.DryRun {
		return infer.CreateResponse[NetworkState]{
			Output: NetworkState{
				Name:        req.Inputs.Name,
				Description: req.Inputs.Description,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
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

	p.GetLogger(ctx).Debugf("Create:NetworkAPI name=%s, id=%s", net.Name, net.Id)

	return infer.CreateResponse[NetworkState]{
		ID: net.Id,
		Output: NetworkState{
			Name:        net.Name,
			Description: net.Description,
		},
	}, nil
}

// Read fetches the current state of a network from NetBird.
func (*Network) Read(ctx context.Context, req infer.ReadRequest[NetworkArgs, NetworkState]) (infer.ReadResponse[NetworkArgs, NetworkState], error) {
	p.GetLogger(ctx).Debugf("Read:NetworkArgs[%s] name=%s", req.ID, req.Inputs.Name)
	p.GetLogger(ctx).Debugf("Read:NetworkState[%s] name=%s, id=%s", req.ID, req.State.Name, req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState]{}, err
	}

	net, err := client.Networks.Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState]{}, fmt.Errorf("reading network failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Read:NetworkAPI[%s] name=%s", net.Id, net.Name)

	return infer.ReadResponse[NetworkArgs, NetworkState]{
		ID: req.ID,
		Inputs: NetworkArgs{
			Name:        net.Name,
			Description: net.Description,
		},
		State: NetworkState{
			Name:        net.Name,
			Description: net.Description,
		},
	}, nil
}

// Update updates the state of the network if needed.
func (*Network) Update(ctx context.Context, req infer.UpdateRequest[NetworkArgs, NetworkState]) (infer.UpdateResponse[NetworkState], error) {
	p.GetLogger(ctx).Debugf("Update:Network[%s] name=%s", req.ID, req.Inputs.Name)

	if req.DryRun {
		return infer.UpdateResponse[NetworkState]{
			Output: NetworkState{
				Name:        req.Inputs.Name,
				Description: req.Inputs.Description,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[NetworkState]{}, err
	}

	_, err = client.Networks.Update(ctx, req.ID, nbapi.NetworkRequest{
		Name:        req.Inputs.Name,
		Description: req.Inputs.Description,
	})
	if err != nil {
		return infer.UpdateResponse[NetworkState]{}, fmt.Errorf("updating network failed: %w", err)
	}

	return infer.UpdateResponse[NetworkState]{
		Output: NetworkState{
			Name:        req.Inputs.Name,
			Description: req.Inputs.Description,
		},
	}, nil
}

// Delete removes a network from NetBird.
func (*Network) Delete(ctx context.Context, req infer.DeleteRequest[NetworkState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:Network[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	err = client.Networks.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting network failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*Network) Diff(ctx context.Context, req infer.DiffRequest[NetworkArgs, NetworkState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Network[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{Kind: p.Update}
	}

	if !equalPtr(req.Inputs.Description, req.State.Description) {
		diff["description"] = p.PropertyDiff{Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:Network[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*Network) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[NetworkArgs], error) {
	p.GetLogger(ctx).Debugf("Check:Network old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())
	args, failures, err := infer.DefaultCheck[NetworkArgs](ctx, req.NewInputs)

	return infer.CheckResponse[NetworkArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*Network) WireDependencies(f infer.FieldSelector, args *NetworkArgs, state *NetworkState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.Description).DependsOn(f.InputField(&args.Description))
}
