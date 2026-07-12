package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// IngressPeer represents a NetBird ingress peer resource.
type IngressPeer struct{}

// Annotate adds a description to the IngressPeer resource type.
func (i *IngressPeer) Annotate(a infer.Annotator) {
	a.Describe(&i, "A NetBird ingress peer: an existing peer designated to receive forwarded ingress traffic.")
}

// IngressPeerArgs defines input fields for creating an ingress peer.
type IngressPeerArgs struct {
	PeerID   string `pulumi:"peerId"`
	Enabled  bool   `pulumi:"enabled"`
	Fallback bool   `pulumi:"fallback"`
}

// Annotate provides documentation for IngressPeerArgs fields.
func (i *IngressPeerArgs) Annotate(a infer.Annotator) {
	a.Describe(&i.PeerID, "ID of the peer used as an ingress peer. Changing this forces a replacement.")
	a.Describe(&i.Enabled, "Whether the ingress peer is enabled.")
	a.Describe(&i.Fallback, "Whether this ingress peer may be used as a fallback when no ingress peer exists in the forwarded peer's region.")
}

// IngressAvailablePorts reports how many forwarding ports remain on an ingress peer.
type IngressAvailablePorts struct {
	TCP int `pulumi:"tcp"`
	UDP int `pulumi:"udp"`
}

// Annotate provides documentation for IngressAvailablePorts fields.
func (i *IngressAvailablePorts) Annotate(a infer.Annotator) {
	a.Describe(&i.TCP, "Number of available TCP ports left on the ingress peer.")
	a.Describe(&i.UDP, "Number of available UDP ports left on the ingress peer.")
}

// IngressPeerState represents the output state of an ingress peer resource.
type IngressPeerState struct {
	PeerID         string                 `pulumi:"peerId"`
	Enabled        bool                   `pulumi:"enabled"`
	Fallback       bool                   `pulumi:"fallback"`
	IngressIP      *string                `pulumi:"ingressIp,optional"`
	Region         *string                `pulumi:"region,optional"`
	Connected      *bool                  `pulumi:"connected,optional"`
	AvailablePorts *IngressAvailablePorts `pulumi:"availablePorts,optional"`
}

// Annotate provides documentation for IngressPeerState fields.
func (i *IngressPeerState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&i.PeerID, "ID of the peer used as an ingress peer.")
	annotator.Describe(&i.Enabled, "Whether the ingress peer is enabled.")
	annotator.Describe(&i.Fallback, "Whether this ingress peer may be used as a fallback.")
	annotator.Describe(&i.IngressIP, "Ingress IP address where forwarded traffic arrives.")
	annotator.Describe(&i.Region, "Region of the ingress peer.")
	annotator.Describe(&i.Connected, "Whether the ingress peer is connected to the management server.")
	annotator.Describe(&i.AvailablePorts, "Forwarding ports remaining on the ingress peer.")
}

// Create creates a new ingress peer.
func (*IngressPeer) Create(ctx context.Context, req infer.CreateRequest[IngressPeerArgs]) (infer.CreateResponse[IngressPeerState], error) {
	p.GetLogger(ctx).Debugf("Create:IngressPeer peerId=%s", req.Inputs.PeerID)

	if req.DryRun {
		return infer.CreateResponse[IngressPeerState]{
			ID:     "preview",
			Output: ingressPeerDryRun(req.Inputs),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[IngressPeerState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	created, err := client.Ingress.Create(ctx, nbapi.IngressPeerCreateRequest{
		PeerId:   req.Inputs.PeerID,
		Enabled:  req.Inputs.Enabled,
		Fallback: req.Inputs.Fallback,
	})
	if err != nil {
		return infer.CreateResponse[IngressPeerState]{}, fmt.Errorf("creating ingress peer failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:IngressPeerAPI id=%s", created.Id)

	return infer.CreateResponse[IngressPeerState]{
		ID:     created.Id,
		Output: ingressPeerStateFromAPI(*created),
	}, nil
}

// Read fetches the current state of an ingress peer from NetBird.
func (*IngressPeer) Read(ctx context.Context, req infer.ReadRequest[IngressPeerArgs, IngressPeerState]) (infer.ReadResponse[IngressPeerArgs, IngressPeerState], error) {
	p.GetLogger(ctx).Debugf("Read:IngressPeer[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[IngressPeerArgs, IngressPeerState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	peer, err := client.Ingress.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[IngressPeerArgs, IngressPeerState]{
				ID:     "",
				Inputs: IngressPeerArgs{},  //nolint:exhaustruct
				State:  IngressPeerState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[IngressPeerArgs, IngressPeerState]{}, fmt.Errorf("reading ingress peer failed: %w", err)
	}

	return infer.ReadResponse[IngressPeerArgs, IngressPeerState]{
		ID: req.ID,
		Inputs: IngressPeerArgs{
			PeerID:   peer.PeerId,
			Enabled:  peer.Enabled,
			Fallback: peer.Fallback,
		},
		State: ingressPeerStateFromAPI(*peer),
	}, nil
}

// Update updates the mutable fields of an ingress peer.
func (*IngressPeer) Update(ctx context.Context, req infer.UpdateRequest[IngressPeerArgs, IngressPeerState]) (infer.UpdateResponse[IngressPeerState], error) {
	p.GetLogger(ctx).Debugf("Update:IngressPeer[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[IngressPeerState]{
			Output: ingressPeerDryRun(req.Inputs),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[IngressPeerState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.Ingress.Update(ctx, req.ID, nbapi.IngressPeerUpdateRequest{
		Enabled:  req.Inputs.Enabled,
		Fallback: req.Inputs.Fallback,
	})
	if err != nil {
		return infer.UpdateResponse[IngressPeerState]{}, fmt.Errorf("updating ingress peer failed: %w", err)
	}

	return infer.UpdateResponse[IngressPeerState]{
		Output: ingressPeerStateFromAPI(*updated),
	}, nil
}

// Delete removes an ingress peer from NetBird.
func (*IngressPeer) Delete(ctx context.Context, req infer.DeleteRequest[IngressPeerState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:IngressPeer[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Ingress.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting ingress peer failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*IngressPeer) Diff(ctx context.Context, req infer.DiffRequest[IngressPeerArgs, IngressPeerState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:IngressPeer[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.PeerID != req.State.PeerID {
		diff["peerId"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Fallback != req.State.Fallback {
		diff["fallback"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for an ingress peer.
func (*IngressPeer) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[IngressPeerArgs], error) {
	p.GetLogger(ctx).Debugf("Check:IngressPeer old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[IngressPeerArgs](ctx, req.NewInputs)

	if isBlank(args.PeerID) {
		failures = append(failures, p.CheckFailure{
			Property: "peerId",
			Reason:   "peerId must not be empty",
		})
	}

	return infer.CheckResponse[IngressPeerArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*IngressPeer) WireDependencies(field infer.FieldSelector, args *IngressPeerArgs, state *IngressPeerState) {
	field.OutputField(&state.PeerID).DependsOn(field.InputField(&args.PeerID))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.Fallback).DependsOn(field.InputField(&args.Fallback))
}

// ingressPeerDryRun builds a preview state from inputs alone.
func ingressPeerDryRun(inputs IngressPeerArgs) IngressPeerState {
	return IngressPeerState{
		PeerID:         inputs.PeerID,
		Enabled:        inputs.Enabled,
		Fallback:       inputs.Fallback,
		IngressIP:      nil,
		Region:         nil,
		Connected:      nil,
		AvailablePorts: nil,
	}
}

// ingressPeerStateFromAPI maps an API ingress peer into resource state.
func ingressPeerStateFromAPI(peer nbapi.IngressPeer) IngressPeerState {
	ingressIP := peer.IngressIp
	region := peer.Region
	connected := peer.Connected

	return IngressPeerState{
		PeerID:    peer.PeerId,
		Enabled:   peer.Enabled,
		Fallback:  peer.Fallback,
		IngressIP: &ingressIP,
		Region:    &region,
		Connected: &connected,
		AvailablePorts: &IngressAvailablePorts{
			TCP: peer.AvailablePorts.Tcp,
			UDP: peer.AvailablePorts.Udp,
		},
	}
}
