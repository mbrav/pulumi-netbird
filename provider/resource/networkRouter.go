package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// TEST: InputDiff: false

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
	NetworkID  string    `pulumi:"network_id"`
	Enabled    bool      `pulumi:"enabled"`
	Masquerade bool      `pulumi:"masquerade"`
	Metric     int       `pulumi:"metric"`
	Peer       *string   `pulumi:"peer,optional"`
	PeerGroups *[]string `pulumi:"peer_groups,optional"`
}

// Annotate provides documentation for NetworkRouter.
func (NetworkRouter) Annotate(a infer.Annotator) {
	a.Describe(&NetworkRouter{}, "A NetBird network router resource.")
}

// Annotate provides documentation for NetworkRouterArgs fields.
func (a *NetworkRouterArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&a.NetworkID, "ID of the network this router belongs to.")
	annotator.Describe(&a.Enabled, "Whether the router is enabled.")
	annotator.Describe(&a.Masquerade, "Whether masquerading is enabled.")
	annotator.Describe(&a.Metric, "Routing metric value.")
	annotator.Describe(&a.Peer, "Optional peer ID associated with this router.")
	annotator.Describe(&a.PeerGroups, "Optional list of peer group IDs associated with this router.")
}

// Annotate provides documentation for NetworkRouterState fields.
func (s *NetworkRouterState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&s.NetworkID, "ID of the network this router belongs to.")
	annotator.Describe(&s.Enabled, "Whether the router is enabled.")
	annotator.Describe(&s.Masquerade, "Whether masquerading is enabled.")
	annotator.Describe(&s.Metric, "Routing metric value.")
	annotator.Describe(&s.Peer, "Optional peer ID associated with this router.")
	annotator.Describe(&s.PeerGroups, "Optional list of peer group IDs associated with this router.")
}

// Create creates a new NetBird network router.
func (*NetworkRouter) Create(ctx context.Context, req infer.CreateRequest[NetworkRouterArgs]) (infer.CreateResponse[NetworkRouterState], error) {
	p.GetLogger(ctx).Debugf("Create:NetworkRouter network_id=%s", req.Inputs.NetworkID)

	if req.DryRun {
		return infer.CreateResponse[NetworkRouterState]{
			ID: "preview",
			Output: NetworkRouterState{
				NetworkID:  req.Inputs.NetworkID,
				Enabled:    req.Inputs.Enabled,
				Masquerade: req.Inputs.Masquerade,
				Metric:     req.Inputs.Metric,
				Peer:       req.Inputs.Peer,
				PeerGroups: req.Inputs.PeerGroups,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[NetworkRouterState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	router, err := client.Networks.Routers(req.Inputs.NetworkID).Create(ctx, nbapi.NetworkRouterRequest{
		Enabled:    req.Inputs.Enabled,
		Masquerade: req.Inputs.Masquerade,
		Metric:     req.Inputs.Metric,
		Peer:       req.Inputs.Peer,
		PeerGroups: req.Inputs.PeerGroups,
	})
	if err != nil {
		return infer.CreateResponse[NetworkRouterState]{}, fmt.Errorf("creating network router failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:NetworkRouterAPI id=%s", router.Id)

	return infer.CreateResponse[NetworkRouterState]{
		ID: router.Id,
		Output: NetworkRouterState{
			NetworkID:  req.Inputs.NetworkID,
			Enabled:    router.Enabled,
			Masquerade: router.Masquerade,
			Metric:     router.Metric,
			Peer:       router.Peer,
			PeerGroups: router.PeerGroups,
		},
	}, nil
}

// Read fetches the current state of a network router from NetBird.
func (*NetworkRouter) Read(ctx context.Context, req infer.ReadRequest[NetworkRouterArgs, NetworkRouterState]) (infer.ReadResponse[NetworkRouterArgs, NetworkRouterState], error) {
	p.GetLogger(ctx).Debugf("Read:NetworkRouter[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[NetworkRouterArgs, NetworkRouterState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	router, err := client.Networks.Routers(req.State.NetworkID).Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[NetworkRouterArgs, NetworkRouterState]{}, fmt.Errorf("reading network router failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Read:NetworkRouterAPI[%s]", router.Id)

	return infer.ReadResponse[NetworkRouterArgs, NetworkRouterState]{
		ID: req.ID,
		Inputs: NetworkRouterArgs{
			NetworkID:  req.State.NetworkID,
			Enabled:    router.Enabled,
			Masquerade: router.Masquerade,
			Metric:     router.Metric,
			Peer:       router.Peer,
			PeerGroups: router.PeerGroups,
		},
		State: NetworkRouterState{
			NetworkID:  req.State.NetworkID,
			Enabled:    router.Enabled,
			Masquerade: router.Masquerade,
			Metric:     router.Metric,
			Peer:       router.Peer,
			PeerGroups: router.PeerGroups,
		},
	}, nil
}

// Update updates the state of the NetBird network router if needed.
func (*NetworkRouter) Update(ctx context.Context, req infer.UpdateRequest[NetworkRouterArgs, NetworkRouterState]) (infer.UpdateResponse[NetworkRouterState], error) {
	p.GetLogger(ctx).Debugf("Update:NetworkRouter[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[NetworkRouterState]{
			Output: NetworkRouterState{
				NetworkID:  req.Inputs.NetworkID,
				Enabled:    req.Inputs.Enabled,
				Masquerade: req.Inputs.Masquerade,
				Metric:     req.Inputs.Metric,
				Peer:       req.Inputs.Peer,
				PeerGroups: req.Inputs.PeerGroups,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[NetworkRouterState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	router, err := client.Networks.Routers(req.Inputs.NetworkID).Update(ctx, req.ID, nbapi.NetworkRouterRequest{
		Enabled:    req.Inputs.Enabled,
		Masquerade: req.Inputs.Masquerade,
		Metric:     req.Inputs.Metric,
		Peer:       req.Inputs.Peer,
		PeerGroups: req.Inputs.PeerGroups,
	})
	if err != nil {
		return infer.UpdateResponse[NetworkRouterState]{}, fmt.Errorf("updating network router failed: %w", err)
	}

	return infer.UpdateResponse[NetworkRouterState]{
		Output: NetworkRouterState{
			NetworkID:  req.Inputs.NetworkID,
			Enabled:    router.Enabled,
			Masquerade: router.Masquerade,
			Metric:     router.Metric,
			Peer:       router.Peer,
			PeerGroups: router.PeerGroups,
		},
	}, nil
}

// Delete removes a network router from NetBird.
func (*NetworkRouter) Delete(ctx context.Context, req infer.DeleteRequest[NetworkRouterState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:NetworkRouter[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Networks.Routers(req.State.NetworkID).Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting network router failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*NetworkRouter) Diff(ctx context.Context, req infer.DiffRequest[NetworkRouterArgs, NetworkRouterState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:NetworkRouter[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Masquerade != req.State.Masquerade {
		diff["masquerade"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Metric != req.State.Metric {
		diff["metric"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if !equalPtr(req.Inputs.Peer, req.State.Peer) {
		diff["peer"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if !equalSlicePtr(req.Inputs.PeerGroups, req.State.PeerGroups) {
		diff["peer_groups"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	p.GetLogger(ctx).Debugf("Diff:NetworkRouter[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}
