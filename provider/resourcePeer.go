package provider

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// Peer represents a resource for managing NetBird peers.
type Peer struct{}

// PeerArgs represents the input arguments for a peer resource.
type PeerArgs struct {
	Name string `pulumi:"name"`
	NbID string `pulumi:"nbId"`
}

// PeerState represents the state of the peer resource.
type PeerState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	PeerArgs
	Name       string `pulumi:"name"`
	SshEnabled bool   `pulumi:"sshEnabled"`
	NbID       string `pulumi:"nbId"`
}

// Peer annotation
func (Peer) Annotate(a infer.Annotator) {
	a.Describe(&Peer{}, "A NetBird peer representing a connected device.")
}

func (p *PeerArgs) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the peer.")
	a.Describe(&p.NbID, "The ID of the peer.")
}

func (p *PeerState) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the peer.")
	a.Describe(&p.SshEnabled, "Whether SSH is enabled for the peer.")
	// a.Describe(&p.NbID, "The ID of the peer.")
}

// Create method is a no-op since peers cannot be created through the API, only imported.
func (Peer) Create(ctx context.Context, name string, input PeerArgs, preview bool) (string, PeerState, error) {
	state := PeerState{
		NbID: input.NbID,
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

	peer, err := client.Peers.Get(ctx, state.NbID)
	if err != nil {
		return inputs, state, fmt.Errorf("reading peer failed: %w", err)
	}

	return PeerArgs{
			NbID: peer.Id,
		}, PeerState{
			NbID:       peer.Id,
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
		NbID: input.NbID,
	}

	if preview {
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	peer, err := client.Peers.Get(ctx, input.NbID)
	if err != nil {
		return "", state, fmt.Errorf("importing peer failed: %w", err)
	}

	state.Name = peer.Name
	state.SshEnabled = peer.SshEnabled

	return name, state, nil
}
