package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
)

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
		Description: input.Description,
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

	return NetworkArgs{
			Name:        net.Name,
			Description: net.Description,
		}, NetworkState{
			Name:        net.Name,
			Description: net.Description,
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
		Description: new.Description,
	})
	if err != nil {
		return state, fmt.Errorf("updating network failed: %w", err)
	}

	return NetworkState{
		NbID:        updated.Id,
		Name:        updated.Name,
		Description: updated.Description,
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

// Import method to import a network resource into Pulumi by its ID.
//
// To import a NetBird network into your Pulumi stack, use the following Pulumi CLI command:
//
//	pulumi import netbird:index:Network <pulumi-resource-name> <netbird-network-id>
//
// Example:
//
//	pulumi import netbird:index:Network corp-network 4c2d2e8d-bf71-4a0c-9b8f-2f4d2cb723d7
func (Network) Import(ctx context.Context, id string, input NetworkArgs, preview bool) (string, NetworkArgs, NetworkState, error) {
	state := NetworkState{}
	args := NetworkArgs{}

	if preview {
		state.NbID = id
		return id, args, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", args, state, err
	}

	network, err := client.Networks.Get(ctx, id)
	if err != nil {
		return "", args, state, fmt.Errorf("importing network failed: %w", err)
	}

	state = NetworkState{
		NbID:        network.Id,
		Name:        network.Name,
		Description: network.Description,
	}

	args = NetworkArgs{
		Name:        network.Name,
		Description: network.Description,
	}

	return id, args, state, nil
}
