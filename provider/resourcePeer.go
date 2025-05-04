package provider

import (
	"context"
	"fmt"
)

// Create method is a no-op since peers cannot be created through the API, only imported.
func (Peer) Create(ctx context.Context, name string, input PeerArgs, preview bool) (string, PeerState, error) {
	state := PeerState{
		PeerID: input.PeerID,
	}

	if preview {
		return name, state, nil
	}

	// Peer cannot be created via the API, only imported, so we return the state as-is
	return name, state, nil
}

// Read method retrieves the state of a peer resource from NetBird API.
func (Peer) Read(ctx context.Context, id string, inputs PeerArgs, state PeerState) (PeerArgs, PeerState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return inputs, state, err
	}

	peer, err := client.Peers.Get(ctx, state.PeerID)
	if err != nil {
		return inputs, state, fmt.Errorf("reading peer failed: %w", err)
	}

	return PeerArgs{
			PeerID: peer.Id,
		}, PeerState{
			PeerID:     peer.Id,
			Name:       peer.Name,
			SshEnabled: peer.SshEnabled,
		}, nil
}

// Update method is a no-op for peers as they cannot be updated via the API.
func (Peer) Update(ctx context.Context, id string, old PeerArgs, new PeerArgs, state PeerState) (PeerState, error) {
	// Peers cannot be updated through Pulumi or the API, so we return the current state.
	return state, nil
}

// Delete method is a no-op for peers as they cannot be deleted through the API.
func (Peer) Delete(ctx context.Context, id string, props PeerState) error {
	// Peers cannot be deleted via Pulumi or the API, so no action is required here.
	return nil
}

// Import method to import a peer resource into Pulumi by its ID.
//
// To import a NetBird peer into your Pulumi stack, use the following Pulumi CLI command:
//
//	pulumi import <pulumi-resource-type> <pulumi-resource-name> <netbird-peer-id>
//
// Replace the placeholders as follows:
//   - <pulumi-resource-type>: the fully qualified Pulumi type, e.g. "netbird:index:Peer"
//   - <pulumi-resource-name>: the logical name you want to give to this resource in your Pulumi program
//   - <netbird-peer-id>: the actual Peer ID from NetBird (e.g., "17a3fa1e-cb8b-4c4b-bfdd-0000abcdef01")
//
// Example:
//
//	pulumi import netbird:index:Peer my-peer 17a3fa1e-cb8b-4c4b-bfdd-0000abcdef01
func (Peer) Import(ctx context.Context, name string, input PeerArgs, preview bool) (string, PeerState, error) {
	state := PeerState{
		PeerID: input.PeerID,
	}

	if preview {
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	peer, err := client.Peers.Get(ctx, input.PeerID)
	if err != nil {
		return "", state, fmt.Errorf("importing peer failed: %w", err)
	}

	state.Name = peer.Name
	state.SshEnabled = peer.SshEnabled

	return name, state, nil
}
