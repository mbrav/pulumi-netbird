package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Peer represents a resource for managing NetBird peers.
type Peer struct{}

// PeerArgs represents the input arguments for a peer resource.
type PeerArgs struct {
	Name                        string `pulumi:"name"`
	InactivityExpirationEnabled bool   `pulumi:"inactivity_expiration_enabled"`
	LoginExpirationEnabled      bool   `pulumi:"login_expiration_enabled"`
	SshEnabled                  bool   `pulumi:"sshEnabled"`
}

// PeerState represents the state of the peer resource.
type PeerState struct {
	Name                        string `pulumi:"name"`
	InactivityExpirationEnabled bool   `pulumi:"inactivity_expiration_enabled"`
	LoginExpirationEnabled      bool   `pulumi:"login_expiration_enabled"`
	SshEnabled                  bool   `pulumi:"sshEnabled"`
}

// Annotate describes the resource and its fields.
func (Peer) Annotate(a infer.Annotator) {
	a.Describe(&Peer{}, "A NetBird peer representing a connected device.")
}

func (p *PeerArgs) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the peer.")
	a.Describe(&p.InactivityExpirationEnabled, "Whether Inactivity Expiration is enabled.")
	a.Describe(&p.LoginExpirationEnabled, "Whether Login Expiration is enabled.")
	a.Describe(&p.SshEnabled, "Whether SSH is enabled.")
}

func (p *PeerState) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the peer.")
	a.Describe(&p.InactivityExpirationEnabled, "Whether Inactivity Expiration is enabled.")
	a.Describe(&p.LoginExpirationEnabled, "Whether Login Expiration is enabled.")
	a.Describe(&p.SshEnabled, "Whether SSH is enabled.")
}

// Create is a no-op; peers must be imported.
func (*Peer) Create(ctx context.Context, req infer.CreateRequest[PeerArgs]) (infer.CreateResponse[PeerState], error) {
	state := PeerState{}

	if req.DryRun {
		return infer.CreateResponse[PeerState]{
			ID:     "preview",
			Output: state,
		}, nil
	}
	return infer.CreateResponse[PeerState]{
		ID:     req.Inputs.Name,
		Output: state,
	}, nil
}

// Read fetches the current state of a peer from NetBird.
func (*Peer) Read(ctx context.Context, req infer.ReadRequest[PeerArgs, PeerState]) (infer.ReadResponse[PeerArgs, PeerState], error) {
	p.GetLogger(ctx).Debugf("Read:Peer[%s]", req.ID)

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[PeerArgs, PeerState]{}, err
	}

	peer, err := client.Peers.Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[PeerArgs, PeerState]{}, fmt.Errorf("reading peer failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Read:PeerAPI[%s] name=%s", peer.Ip, peer.Name)

	return infer.ReadResponse[PeerArgs, PeerState]{
		ID: req.ID,
		Inputs: PeerArgs{
			Name:                        peer.Name,
			InactivityExpirationEnabled: peer.InactivityExpirationEnabled,
			LoginExpirationEnabled:      peer.LoginExpirationEnabled,
			SshEnabled:                  peer.SshEnabled,
		},
		State: PeerState{
			Name:                        peer.Name,
			InactivityExpirationEnabled: peer.InactivityExpirationEnabled,
			LoginExpirationEnabled:      peer.LoginExpirationEnabled,
			SshEnabled:                  peer.SshEnabled,
		},
	}, nil
}

// Update updates the state of the NetBird Peer if needed.
func (*Peer) Update(ctx context.Context, req infer.UpdateRequest[PeerArgs, PeerState]) (infer.UpdateResponse[PeerState], error) {
	p.GetLogger(ctx).Debugf("Update:Peer[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[PeerState]{
			Output: PeerState{
				Name:                        req.Inputs.Name,
				InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
				LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
				SshEnabled:                  req.Inputs.SshEnabled,
			},
		}, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[PeerState]{}, err
	}

	_, err = client.Peers.Update(ctx, req.ID, nbapi.PeerRequest{
		Name:                        req.Inputs.Name,
		InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
		LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
		SshEnabled:                  req.Inputs.SshEnabled,
	})
	if err != nil {
		return infer.UpdateResponse[PeerState]{}, fmt.Errorf("updating peer failed: %w", err)
	}

	return infer.UpdateResponse[PeerState]{
		Output: PeerState{
			Name:                        req.Inputs.Name,
			InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
			LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
			SshEnabled:                  req.Inputs.SshEnabled,
		},
	}, nil
}

// Diff detects changes between inputs and prior state.
func (*Peer) Diff(ctx context.Context, req infer.DiffRequest[PeerArgs, PeerState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Peer[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{Kind: p.Update}
	}
	if req.Inputs.InactivityExpirationEnabled != req.State.InactivityExpirationEnabled {
		diff["inactivity_expiration_enabled"] = p.PropertyDiff{Kind: p.Update}
	}
	if req.Inputs.LoginExpirationEnabled != req.State.LoginExpirationEnabled {
		diff["login_expiration_enabled"] = p.PropertyDiff{Kind: p.Update}
	}
	if req.Inputs.SshEnabled != req.State.SshEnabled {
		diff["sshEnabled"] = p.PropertyDiff{Kind: p.Update}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}
