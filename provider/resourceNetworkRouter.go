package provider

import (
	"context"
	"fmt"
	"strings"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// FIX: Recreate resource on UPDATE

// NetworkRouter represents a Pulumi resource for NetBird network routers.
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

// NetworkRouterState represents the state of a network router.
type NetworkRouterState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	NetworkRouterArgs
	NetworkID  string    `pulumi:"network_id"`
	Enabled    bool      `pulumi:"enabled"`
	Masquerade bool      `pulumi:"masquerade"`
	Metric     int       `pulumi:"metric"`
	Peer       *string   `pulumi:"peer,optional"`
	PeerGroups *[]string `pulumi:"peer_groups,optional"`
	NbID       string    `pulumi:"nbId"`
}

// NetworkRouter annotation
func (NetworkRouter) Annotate(a infer.Annotator) {
	a.Describe(&NetworkRouter{}, "A NetBird router used to route traffic between peers or networks.")
}

func (r *NetworkRouterArgs) Annotate(a infer.Annotator) {
	a.Describe(&r.NetworkID, "The ID of the network that this router is associated with.")
	a.Describe(&r.Enabled, "Whether the router is enabled.")
	a.Describe(&r.Masquerade, "Whether NAT masquerading is enabled for this router.")
	a.Describe(&r.Metric, "The routing metric (priority) for this router.")
	a.Describe(&r.Peer, "Optional peer ID to route through.")
	a.Describe(&r.PeerGroups, "Optional list of peer group IDs to use as routing targets.")
}

func (r *NetworkRouterState) Annotate(a infer.Annotator) {
	a.Describe(&r.NbID, "The internal NetBird ID of the router.")
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

// Import allows importing an existing NetBird network router resource by its ID.
//
// Expected import ID format: <network-id>/<router-id>
//
// Example:
//
//	pulumi import netbird:index:NetworkRouter core-router 12345678-abcd-ef01-2345-6789abcdef01/abcdef12-3456-7890-abcd-ef1234567890
func (NetworkRouter) Import(ctx context.Context, name string, input NetworkRouterArgs, preview bool) (string, NetworkRouterState, error) {
	state := NetworkRouterState{}

	ids := strings.SplitN(name, "/", 2)
	if len(ids) != 2 {
		return "", state, fmt.Errorf("invalid import ID format, expected <network-id>/<router-id>")
	}
	networkID := ids[0]
	routerID := ids[1]

	if preview {
		state.NetworkID = networkID
		state.NbID = routerID
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	router, err := client.Networks.Routers(networkID).Get(ctx, routerID)
	if err != nil {
		return "", state, fmt.Errorf("importing network router failed: %w", err)
	}

	state = NetworkRouterState{
		NbID:       router.Id,
		NetworkID:  networkID,
		Enabled:    router.Enabled,
		Masquerade: router.Masquerade,
		Metric:     router.Metric,
		Peer:       router.Peer,
		PeerGroups: router.PeerGroups,
	}

	return name, state, nil
}
