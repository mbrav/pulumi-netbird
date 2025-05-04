package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
)

// Network represents a resource for managing NetBird networks.
type Network struct{}

// NetworkArgs represents the input arguments for creating or updating a network.
type NetworkArgs struct {
	Name        string `pulumi:"name"`
	Description string `pulumi:"description"`
}

// NetworkState represents the state of the network resource.
type NetworkState struct {
	Name        string `pulumi:"name"`
	Description string `pulumi:"description"`
	NbID        string `pulumi:"nbId"`
}

// Create a new network resource.
func (Network) Create(ctx context.Context, name string, input NetworkArgs, preview bool) (string, NetworkState, error) {
	state := NetworkState{
		Name:        input.Name,
		Description: input.Description,
	}

	if preview {
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	created, err := client.Networks.Create(ctx, nbapi.NetworkRequest{
		Name:        input.Name,
		Description: &input.Description,
	})
	if err != nil {
		return "", state, fmt.Errorf("creating network failed: %w", err)
	}

	state.NbID = created.Id
	return name, state, nil
}

// Read an existing network resource by ID.
func (Network) Read(ctx context.Context, id string, inputs NetworkArgs, state NetworkState) (NetworkArgs, NetworkState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return inputs, state, err
	}

	net, err := client.Networks.Get(ctx, state.NbID)
	if err != nil {
		return inputs, state, fmt.Errorf("reading network failed: %w", err)
	}

	desc := stringOrEmpty(net.Description)

	return NetworkArgs{
			Name:        net.Name,
			Description: desc,
		}, NetworkState{
			Name:        net.Name,
			Description: desc,
			NbID:        net.Id,
		}, nil
}

// Update an existing network resource.
func (Network) Update(ctx context.Context, id string, old NetworkArgs, new NetworkArgs, state NetworkState) (NetworkState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return state, err
	}

	updated, err := client.Networks.Update(ctx, state.NbID, nbapi.NetworkRequest{
		Name:        new.Name,
		Description: &new.Description,
	})
	if err != nil {
		return state, fmt.Errorf("updating network failed: %w", err)
	}

	return NetworkState{
		NbID:        updated.Id,
		Name:        updated.Name,
		Description: stringOrEmpty(updated.Description),
	}, nil
}

// Delete an existing network resource.
func (Network) Delete(ctx context.Context, id string, props NetworkState) error {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return err
	}

	if err := client.Networks.Delete(ctx, props.NbID); err != nil {
		return fmt.Errorf("deleting network failed: %w", err)
	}
	return nil
}
