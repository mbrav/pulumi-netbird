package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// TEST: InputDiff: false

// Peer represents a resource for managing NetBird peers.
type Peer struct{}

// Annotate describes the resource and its fields.
func (peer *Peer) Annotate(a infer.Annotator) {
	a.Describe(peer, "A NetBird peer representing a connected device.")
}

// PeerArgs represents the input arguments for a peer resource.
type PeerArgs struct {
	Name                        string `pulumi:"name"`
	InactivityExpirationEnabled bool   `pulumi:"inactivityExpirationEnabled"`
	LoginExpirationEnabled      bool   `pulumi:"loginExpirationEnabled"`
	SSHEnabled                  bool   `pulumi:"sshEnabled"`
	// Cloud Only
	ApprovalRequired *bool `pulumi:"approvalRequired"`
}

// Annotate adds descriptive annotations to the PeerArgs fields for use in generated SDKs.
func (p *PeerArgs) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the peer.")
	a.Describe(&p.InactivityExpirationEnabled, "Whether Inactivity Expiration is enabled.")
	a.Describe(&p.LoginExpirationEnabled, "Whether Login Expiration is enabled.")
	a.Describe(&p.SSHEnabled, "Whether SSH is enabled.")
	a.Deprecate(&p.ApprovalRequired, "Cloud only, not maintained in this provider")
}

// PeerState represents the state of the peer resource.
type PeerState struct {
	Name                        string `pulumi:"name"`
	InactivityExpirationEnabled bool   `pulumi:"inactivityExpirationEnabled"`
	LoginExpirationEnabled      bool   `pulumi:"loginExpirationEnabled"`
	SSHEnabled                  bool   `pulumi:"sshEnabled"`
	// Cloud Only
	ApprovalRequired *bool `pulumi:"approvalRequired"`
}

// Annotate adds descriptive annotations to the PeerState fields for use in generated SDKs.
func (p *PeerState) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the peer.")
	a.Describe(&p.InactivityExpirationEnabled, "Whether Inactivity Expiration is enabled.")
	a.Describe(&p.LoginExpirationEnabled, "Whether Login Expiration is enabled.")
	a.Describe(&p.SSHEnabled, "Whether SSH is enabled.")
	a.Deprecate(&p.ApprovalRequired, "Cloud only, not maintained in this provider")
}

// Create is a no-op; peers must be imported.
func (*Peer) Create(_ context.Context, req infer.CreateRequest[PeerArgs]) (infer.CreateResponse[PeerState], error) {
	state := PeerState{
		Name:                        req.Inputs.Name,
		InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
		LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
		SSHEnabled:                  req.Inputs.SSHEnabled,
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
		return infer.ReadResponse[PeerArgs, PeerState]{}, fmt.Errorf("error getting NetBird client: %w", err)
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
			SSHEnabled:                  peer.SshEnabled,
			ApprovalRequired:            nil,
		},
		State: PeerState{
			Name:                        peer.Name,
			InactivityExpirationEnabled: peer.InactivityExpirationEnabled,
			LoginExpirationEnabled:      peer.LoginExpirationEnabled,
			SSHEnabled:                  peer.SshEnabled,
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
				SSHEnabled:                  req.Inputs.SSHEnabled,
				ApprovalRequired:            nil,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[PeerState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	_, err = client.Peers.Update(ctx, req.ID, nbapi.PeerRequest{
		Name:                        req.Inputs.Name,
		InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
		LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
		SshEnabled:                  req.Inputs.SSHEnabled,
		ApprovalRequired:            nil, // ApprovalRequired is not supported in for Cloud version only
		Ip:                          nil,
	})
	if err != nil {
		return infer.UpdateResponse[PeerState]{}, fmt.Errorf("updating peer failed: %w", err)
	}

	return infer.UpdateResponse[PeerState]{
		Output: PeerState{
			Name:                        req.Inputs.Name,
			InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
			LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
			SSHEnabled:                  req.Inputs.SSHEnabled,
			ApprovalRequired:            nil,
		},
	}, nil
}

// Delete removes a peer from NetBird.
func (*Peer) Delete(ctx context.Context, req infer.DeleteRequest[PeerState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:Peer[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Peers.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting peer failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
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

	if req.Inputs.SSHEnabled != req.State.SSHEnabled {
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

// Check provides input validation and default setting.
func (*Peer) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[PeerArgs], error) {
	p.GetLogger(ctx).Debugf("Check:Peer old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())
	args, failures, err := infer.DefaultCheck[PeerArgs](ctx, req.NewInputs)

	return infer.CheckResponse[PeerArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*Peer) WireDependencies(f infer.FieldSelector, args *PeerArgs, state *PeerState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.InactivityExpirationEnabled).DependsOn(f.InputField(&args.InactivityExpirationEnabled))
	f.OutputField(&state.LoginExpirationEnabled).DependsOn(f.InputField(&args.LoginExpirationEnabled))
	f.OutputField(&state.SSHEnabled).DependsOn(f.InputField(&args.SSHEnabled))
	f.OutputField(&state.ApprovalRequired).DependsOn(f.InputField(&args.ApprovalRequired))
}
