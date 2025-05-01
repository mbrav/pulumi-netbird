package provider

import (
	"context"
	"fmt"

	"github.com/netbirdio/netbird/management/server/http/api"
)

// FIX: Not updateing Enabled state on UPDATE

// NetworkResource represents a Pulumi resource for NetBird network resources.
type NetworkResource struct{}

// NetworkResourceArgs represents the input arguments for creating or updating a network resource.
type NetworkResourceArgs struct {
	Name        string    `pulumi:"name"`
	Description *string   `pulumi:"description,optional"`
	NetworkID   string    `pulumi:"network_id"`
	Address     string    `pulumi:"address"`
	Enabled     bool      `pulumi:"enabled"`
	GroupIDs    *[]string `pulumi:"group_ids,optional"`
}

// NetworkResourceState represents the state of a network resource.
type NetworkResourceState struct {
	Name        string    `pulumi:"name"`
	Description *string   `pulumi:"description"`
	NbID        string    `pulumi:"nbId"`
	NetworkID   string    `pulumi:"network_id"`
	Address     string    `pulumi:"address"`
	Enabled     bool      `pulumi:"enabled"`
	GroupIDs    *[]string `pulumi:"group_ids,optional"`
}

func (NetworkResource) Create(ctx context.Context, name string, input NetworkResourceArgs, preview bool) (string, NetworkResourceState, error) {
	state := NetworkResourceState{
		Name:        input.Name,
		Description: input.Description,
		NetworkID:   input.NetworkID,
		Address:     input.Address,
		Enabled:     input.Enabled,
		GroupIDs:    input.GroupIDs,
	}

	if preview {
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	created, err := client.Networks.Resources(input.NetworkID).Create(ctx, api.NetworkResourceRequest{
		Name:        input.Name,
		Address:     input.Address,
		Description: input.Description,
		Enabled:     input.Enabled,
		Groups:      *input.GroupIDs,
	})
	if err != nil {
		return "", state, fmt.Errorf("creating network resource failed: %w", err)
	}

	state.NbID = created.Id
	return name, state, nil
}

func (NetworkResource) Read(ctx context.Context, id string, input NetworkResourceArgs, state NetworkResourceState) (NetworkResourceArgs, NetworkResourceState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return input, state, err
	}

	res, err := client.Networks.Resources(state.NetworkID).Get(ctx, state.NbID)
	if err != nil {
		return input, state, fmt.Errorf("reading network resource failed: %w", err)
	}

	groupIDs := make([]string, len(res.Groups))
	for i, peer := range res.Groups {
		groupIDs[i] = peer.Id
	}

	return NetworkResourceArgs{
			Name:        res.Name,
			Description: res.Description,
			NetworkID:   state.NetworkID,
			Address:     res.Address,
			Enabled:     res.Enabled,
			GroupIDs:    &groupIDs,
		}, NetworkResourceState{
			Name:        res.Name,
			Description: res.Description,
			NbID:        res.Id,
			NetworkID:   state.NetworkID,
			Address:     res.Address,
			Enabled:     res.Enabled,
			GroupIDs:    &groupIDs,
		}, nil
}

func (NetworkResource) Update(ctx context.Context, id string, old NetworkResourceArgs, new NetworkResourceArgs, state NetworkResourceState) (NetworkResourceState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return state, err
	}

	updated, err := client.Networks.Resources(state.NetworkID).Update(ctx, state.NbID, api.NetworkResourceRequest{
		Name:        new.Name,
		Address:     new.Address,
		Description: new.Description,
		Enabled:     new.Enabled,
		Groups:      *new.GroupIDs,
	})
	if err != nil {
		return state, fmt.Errorf("updating network resource failed: %w", err)
	}

	groupIDs := make([]string, len(updated.Groups))
	for i, peer := range updated.Groups {
		groupIDs[i] = peer.Id
	}

	return NetworkResourceState{
		Name:        updated.Name,
		Description: updated.Description,
		NbID:        updated.Id,
		NetworkID:   state.NetworkID,
		Address:     updated.Address,
		Enabled:     updated.Enabled,
		GroupIDs:    &groupIDs,
	}, nil
}

func (NetworkResource) Delete(ctx context.Context, id string, state NetworkResourceState) error {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return err
	}

	if err := client.Networks.Resources(state.NetworkID).Delete(ctx, state.NbID); err != nil {
		return fmt.Errorf("deleting network resource failed: %w", err)
	}
	return nil
}
