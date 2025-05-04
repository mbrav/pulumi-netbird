package provider

import (
	"context"
	"fmt"
	"strings"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
)

// FIX: Not updateing Enabled state on UPDATE

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

	created, err := client.Networks.Resources(input.NetworkID).Create(ctx, nbapi.NetworkResourceRequest{
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

	updated, err := client.Networks.Resources(state.NetworkID).Update(ctx, state.NbID, nbapi.NetworkResourceRequest{
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

// Import allows importing an existing NetBird network resource by its ID.
//
// To import a NetBird network resource into your Pulumi stack, use:
//
//	pulumi import netbird:index:NetworkResource <pulumi-resource-name> <network-id>/<resource-id>
//
// Example:
//
//	pulumi import netbird:index:NetworkResource core-dns 70abf594-fb68-4ec9-84b9-672758bfc1a3/2ab785be-1439-4ef3-a2b4-cfb622a6f60a
func (NetworkResource) Import(ctx context.Context, name string, input NetworkResourceArgs, preview bool) (string, NetworkResourceState, error) {
	state := NetworkResourceState{}

	if preview {
		// Expect NetworkID and NbID to be passed as "network-id/resource-id" in input.Name
		ids := strings.SplitN(input.Name, "/", 2)
		if len(ids) != 2 {
			return "", state, fmt.Errorf("invalid import ID format, expected <network-id>/<resource-id>")
		}
		state.NetworkID = ids[0]
		state.NbID = ids[1]
		return name, state, nil
	}

	ids := strings.SplitN(input.Name, "/", 2)
	if len(ids) != 2 {
		return "", state, fmt.Errorf("invalid import ID format, expected <network-id>/<resource-id>")
	}
	networkID := ids[0]
	resourceID := ids[1]

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	res, err := client.Networks.Resources(networkID).Get(ctx, resourceID)
	if err != nil {
		return "", state, fmt.Errorf("importing network resource failed: %w", err)
	}

	groupIDs := make([]string, len(res.Groups))
	for i, grp := range res.Groups {
		groupIDs[i] = grp.Id
	}

	state = NetworkResourceState{
		Name:        res.Name,
		Description: res.Description,
		NbID:        res.Id,
		NetworkID:   networkID,
		Address:     res.Address,
		Enabled:     res.Enabled,
		GroupIDs:    &groupIDs,
	}

	return name, state, nil
}
