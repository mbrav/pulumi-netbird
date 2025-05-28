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

// Peer represents a resource for managing NetBird peers.
type Peer struct{}

// PeerArgs represents the input arguments for a peer resource.
type PeerArgs struct {
	Name                        string `pulumi:"name"`
	InactivityExpirationEnabled bool   `pulumi:"inactivityExpirationEnabled"`
	LoginExpirationEnabled      bool   `pulumi:"loginExpirationEnabled"`
	SshEnabled                  bool   `pulumi:"sshEnabled"`
	// Cloud Only
	ApprovalRequired *bool `pulumi:"approvalRequired"`
}

// PeerState represents the state of the peer resource.
type PeerState struct {
	Name                        string `pulumi:"name"`
	InactivityExpirationEnabled bool   `pulumi:"inactivityExpirationEnabled"`
	LoginExpirationEnabled      bool   `pulumi:"loginExpirationEnabled"`
	SshEnabled                  bool   `pulumi:"sshEnabled"`
	// Cloud Only
	ApprovalRequired *bool `pulumi:"approvalRequired"`
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
	a.Deprecate(&p.ApprovalRequired, "Cloud only, not maintained in this provider")
}

func (p *PeerState) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the peer.")
	a.Describe(&p.InactivityExpirationEnabled, "Whether Inactivity Expiration is enabled.")
	a.Describe(&p.LoginExpirationEnabled, "Whether Login Expiration is enabled.")
	a.Describe(&p.SshEnabled, "Whether SSH is enabled.")
	a.Deprecate(&p.ApprovalRequired, "Cloud only, not maintained in this provider")
}

// Create is a no-op; peers must be imported.
func (*Peer) Create(ctx context.Context, req infer.CreateRequest[PeerArgs]) (infer.CreateResponse[PeerState], error) {
	state := PeerState{
		Name:                        req.Inputs.Name,
		InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
		LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
		SshEnabled:                  req.Inputs.SshEnabled,
		ApprovalRequired:            nil,
	}

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

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[PeerArgs, PeerState]{}, fmt.Errorf("error getting Netbird client: %w", err)
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
			ApprovalRequired:            nil,
		},
		State: PeerState{
			Name:                        peer.Name,
			InactivityExpirationEnabled: peer.InactivityExpirationEnabled,
			LoginExpirationEnabled:      peer.LoginExpirationEnabled,
			SshEnabled:                  peer.SshEnabled,
			ApprovalRequired:            nil,
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
				ApprovalRequired:            nil,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[PeerState]{}, fmt.Errorf("error getting Netbird client: %w", err)
	}

	_, err = client.Peers.Update(ctx, req.ID, nbapi.PeerRequest{
		Name:                        req.Inputs.Name,
		InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
		LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
		SshEnabled:                  req.Inputs.SshEnabled,
		ApprovalRequired:            nil, // ApprovalRequired is not supported in for Cloud version only
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
			ApprovalRequired:            nil,
		},
	}, nil
}

// Diff detects changes between inputs and prior state.
func (*Peer) Diff(ctx context.Context, req infer.DiffRequest[PeerArgs, PeerState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Peer[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	// if *req.Inputs.ApprovalRequired != *req.State.ApprovalRequired {
	// 	diff["approvalRequired"] = p.PropertyDiff{
	// 		InputDiff: false,
	// 		Kind:      p.Update,
	// 	}
	// }

	if req.Inputs.InactivityExpirationEnabled != req.State.InactivityExpirationEnabled {
		diff["inactivityExpirationEnabled"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.LoginExpirationEnabled != req.State.LoginExpirationEnabled {
		diff["loginExpirationEnabled"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.SshEnabled != req.State.SshEnabled {
		diff["sshEnabled"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}
