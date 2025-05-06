package provider

import (
	"context"
	"fmt"
	"strings"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	"github.com/pulumi/pulumi-go-provider/infer"
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
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	NetworkResourceArgs
	Name        string    `pulumi:"name"`
	Description *string   `pulumi:"description,optional"`
	NetworkID   string    `pulumi:"network_id"`
	Address     string    `pulumi:"address"`
	Enabled     bool      `pulumi:"enabled"`
	GroupIDs    *[]string `pulumi:"group_ids,optional"`
	NbID        string    `pulumi:"nbId"`
}

// NetworkResource annotation
func (NetworkResource) Annotate(a infer.Annotator) {
	a.Describe(&NetworkResource{}, "A NetBird network resource, such as a CIDR range assigned to the network.")
}

func (n *NetworkResourceArgs) Annotate(a infer.Annotator) {
	a.Describe(&n.Name, "The name of the network resource.")
	a.Describe(&n.Description, "An optional description of the network resource.")
	a.Describe(&n.NetworkID, "The ID of the associated network.")
	a.Describe(&n.Address, "The IP address or subnet of the network resource.")
	a.Describe(&n.Enabled, "Indicates if the resource is currently enabled.")
	a.Describe(&n.GroupIDs, "Optional list of group IDs to associate with this network resource.")
}

func (n *NetworkResourceState) Annotate(a infer.Annotator) {
	a.Describe(&n.NbID, "The internal NetBird ID of the network resource.")
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
