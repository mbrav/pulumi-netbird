package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
)

// FIX: Recreate resource on UPDATE

// NetworkRouter represents a Pulumi resource for NetBird network resources.
type NetworkRouter struct{}

// NetworkRouterArgs represents the input arguments for creating or updating a network router.
type NetworkRouterArgs struct {
	NetworkID  string    `pulumi:"network_id"`
	Enabled    bool      `pulumi:"enabled"`
	Masquerade bool      `pulumi:"masquerade"`
	Metric     int       `pulumi:"metric"`
	Peer       *string   `pulumi:"peer,optional"`
	PeerGroups *[]string `pulumi:"peer_groups,optional"`
}

// NetworkResourceArgs represents the state of a network router.
type NetworkRouterState struct {
	NbID       string    `pulumi:"nbId"`
	NetworkID  string    `pulumi:"network_id"`
	Enabled    bool      `pulumi:"enabled"`
	Masquerade bool      `pulumi:"masquerade"`
	Metric     int       `pulumi:"metric"`
	Peer       *string   `pulumi:"peer"`
	PeerGroups *[]string `pulumi:"peer_groups"`
}

func (NetworkRouter) Create(ctx context.Context, name string, input NetworkRouterArgs, preview bool) (string, NetworkRouterState, error) {
	state := NetworkRouterState{
		NetworkID:  input.NetworkID,
		Enabled:    input.Enabled,
		Masquerade: input.Masquerade,
		Metric:     input.Metric,
		Peer:       input.Peer,
		PeerGroups: input.PeerGroups,
	}

	if preview {
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	created, err := client.Networks.Routers(input.NetworkID).Create(ctx, nbapi.NetworkRouterRequest{
		Enabled:    input.Enabled,
		Masquerade: input.Masquerade,
		Metric:     input.Metric,
		Peer:       input.Peer,
		PeerGroups: input.PeerGroups,
	})
	if err != nil {
		return "", state, fmt.Errorf("creating network router failed: %w", err)
	}

	state.NbID = created.Id
	return name, state, nil
}

func (NetworkRouter) Read(ctx context.Context, id string, input NetworkRouterArgs, state NetworkRouterState) (NetworkRouterArgs, NetworkRouterState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return input, state, err
	}
	res, err := client.Networks.Routers(state.NetworkID).Get(ctx, state.NbID)
	if err != nil {
		return input, state, fmt.Errorf("reading network router failed: %w", err)
	}
	return NetworkRouterArgs{
			NetworkID:  state.NetworkID,
			Enabled:    res.Enabled,
			Masquerade: res.Masquerade,
			Metric:     res.Metric,
			Peer:       res.Peer,
			PeerGroups: res.PeerGroups,
		}, NetworkRouterState{
			NbID:       res.Id,
			NetworkID:  state.NetworkID,
			Enabled:    res.Enabled,
			Masquerade: res.Masquerade,
			Metric:     res.Metric,
			Peer:       res.Peer,
			PeerGroups: res.PeerGroups,
		}, nil
}

func (NetworkRouter) Update(ctx context.Context, id string, old NetworkRouterArgs, new NetworkRouterArgs, state NetworkRouterState) (NetworkRouterState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return state, err
	}
	updated, err := client.Networks.Routers(state.NetworkID).Update(ctx, state.NbID, nbapi.NetworkRouterRequest{
		Enabled:    new.Enabled,
		Masquerade: new.Masquerade,
		Metric:     new.Metric,
		Peer:       new.Peer,
		PeerGroups: new.PeerGroups,
	})
	if err != nil {
		return state, fmt.Errorf("updating network router failed: %w", err)
	}
	return NetworkRouterState{
		NbID:       updated.Id,
		NetworkID:  state.NetworkID,
		Enabled:    updated.Enabled,
		Masquerade: updated.Masquerade,
		Metric:     updated.Metric,
		Peer:       updated.Peer,
		PeerGroups: updated.PeerGroups,
	}, nil
}

func (NetworkRouter) Delete(ctx context.Context, id string, state NetworkRouterState) error {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return err
	}
	if err := client.Networks.Routers(state.NetworkID).Delete(ctx, state.NbID); err != nil {
		return fmt.Errorf("deleting network router failed: %w", err)
	}
	return nil
}
