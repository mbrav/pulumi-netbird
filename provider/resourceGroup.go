package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
)

// Group represents a resource for managing NetBird groups.
type Group struct{}

// GroupArgs represents the input arguments for creating or updating a group.
type GroupArgs struct {
	Name  string    `pulumi:"name"`
	Peers *[]string `pulumi:"peers,optional"`
}

// GroupState represents the state of the group resource.
type GroupState struct {
	Name  string    `pulumi:"name"`
	Peers *[]string `pulumi:"peers,optional"`
	NbID  string    `pulumi:"nbId"`
}

// Create a new group resource.
func (Group) Create(ctx context.Context, name string, input GroupArgs, preview bool) (string, GroupState, error) {
	state := GroupState{
		Name:  input.Name,
		Peers: input.Peers,
	}

	if preview {
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	created, err := client.Groups.Create(ctx, nbapi.GroupRequest{
		Name:  input.Name,
		Peers: input.Peers,
	})
	if err != nil {
		return "", state, fmt.Errorf("creating group failed: %w", err)
	}

	state.NbID = created.Id
	return name, state, nil
}

// Read retrieves an existing group resource by ID.
func (Group) Read(ctx context.Context, id string, inputs GroupArgs, state GroupState) (GroupArgs, GroupState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return inputs, state, err
	}

	group, err := client.Groups.Get(ctx, state.NbID)
	if err != nil {
		return inputs, state, fmt.Errorf("reading group failed: %w", err)
	}

	// Convert []api.PeerMinimum to []string
	peerIDs := make([]string, len(group.Peers))
	for i, peer := range group.Peers {
		peerIDs[i] = peer.Id
	}

	return GroupArgs{
			Name:  group.Name,
			Peers: &peerIDs,
		}, GroupState{
			Name:  group.Name,
			Peers: &peerIDs,
			NbID:  group.Id,
		}, nil
}

// Update modifies an existing group resource.
func (Group) Update(ctx context.Context, id string, old GroupArgs, new GroupArgs, state GroupState) (GroupState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return state, err
	}

	updated, err := client.Groups.Update(ctx, state.NbID, nbapi.GroupRequest{
		Name:  new.Name,
		Peers: new.Peers,
	})
	if err != nil {
		return state, fmt.Errorf("updating group failed: %w", err)
	}

	peerIDs := make([]string, len(updated.Peers))
	for i, peer := range updated.Peers {
		peerIDs[i] = peer.Id
	}

	return GroupState{
		NbID:  updated.Id,
		Name:  updated.Name,
		Peers: &peerIDs,
	}, nil
}

// Delete removes an existing group resource.
func (Group) Delete(ctx context.Context, id string, props GroupState) error {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return err
	}

	if err := client.Groups.Delete(ctx, props.NbID); err != nil {
		return fmt.Errorf("deleting group failed: %w", err)
	}
	return nil
}
