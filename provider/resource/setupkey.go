package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

const setupKeyStateValid = "valid"

// SetupKey represents a resource for managing NetBird setup keys.
type SetupKey struct{}

// Annotate adds a description to the SetupKey resource type.
func (s *SetupKey) Annotate(a infer.Annotator) {
	a.Describe(&s, "Manages a NetBird setup key.")
}

// SetupKeyArgs represents the input arguments for creating a setup key.
type SetupKeyArgs struct {
	Name                string       `pulumi:"name"`
	Type                SetupKeyType `pulumi:"type"`      // "one-off" | "reusable"
	ExpiresIn           int          `pulumi:"expiresIn"` // seconds
	AutoGroups          []string     `pulumi:"autoGroups"`
	UsageLimit          int          `pulumi:"usageLimit"` // 0 = unlimited
	Ephemeral           *bool        `pulumi:"ephemeral,optional"`
	AllowExtraDNSLabels *bool        `pulumi:"allowExtraDnsLabels,optional"`
}

// Annotate provides documentation for SetupKeyArgs fields.
func (a *SetupKeyArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&a.Name, "Setup key display name.")
	annotator.Describe(&a.Type, "Setup key type: 'one-off' (single use) or 'reusable'.")
	annotator.Describe(&a.ExpiresIn, "Time-to-live in seconds from creation; use 0 for no expiration if supported by the API.")
	annotator.Describe(&a.AutoGroups, "Group IDs to auto-assign to peers created with this key.")
	annotator.Describe(&a.UsageLimit, "Maximum uses for reusable keys; 0 = unlimited.")
	annotator.Describe(&a.Ephemeral, "Whether peers registered with this key are ephemeral (auto-expire).")
	annotator.Describe(&a.AllowExtraDNSLabels, "Allow peers to add extra DNS labels beyond the base peer name.")
}

// SetupKeyState represents the state/output of a setup key resource.
type SetupKeyState struct {
	SetupKeyArgs

	Key       *string `pulumi:"key,optional"`
	Valid     *bool   `pulumi:"valid,optional"`
	Revoked   *bool   `pulumi:"revoked,optional"`
	UsedTimes *int    `pulumi:"usedTimes,optional"`
	LastUsed  *string `pulumi:"lastUsed,optional"`
	Expires   *string `pulumi:"expires,optional"`
	State     *string `pulumi:"state,optional"`
	UpdatedAt *string `pulumi:"updatedAt,optional"`
}

// SetupKeyType defines the kind of setup key accepted by NetBird.
type SetupKeyType string

const (
	// SetupKeyTypeReusable creates a key that can be used multiple times.
	SetupKeyTypeReusable SetupKeyType = SetupKeyType("reusable")
	// SetupKeyTypeOneOff creates a key that can only be used once.
	SetupKeyTypeOneOff SetupKeyType = SetupKeyType("one-off")
)

// Values describes the setup key type enum for schema generation.
func (SetupKeyType) Values() []infer.EnumValue[Type] {
	return []infer.EnumValue[Type]{
		{Name: "reusable", Value: Type(SetupKeyTypeReusable), Description: "Reusable setup key that supports multiple peers."},
		{Name: "one-off", Value: Type(SetupKeyTypeOneOff), Description: "One-off setup key that can be used only once."},
	}
}

// Create creates a new NetBird setup key.
func (*SetupKey) Create(ctx context.Context, req infer.CreateRequest[SetupKeyArgs]) (infer.CreateResponse[SetupKeyState], error) {
	p.GetLogger(ctx).Debugf("Create:SetupKey name=%s, type=%s", req.Inputs.Name, req.Inputs.Type)

	if req.DryRun {
		return infer.CreateResponse[SetupKeyState]{
			ID: "preview",
			Output: SetupKeyState{
				SetupKeyArgs: req.Inputs,
				Key:          nil,
				Valid:        nil,
				Revoked:      nil,
				UsedTimes:    nil,
				LastUsed:     nil,
				Expires:      nil,
				State:        nil,
				UpdatedAt:    nil,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[SetupKeyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// Use CreateSetupKeyRequest for creation
	createReq := nbapi.CreateSetupKeyRequest{
		Name:                req.Inputs.Name,
		Type:                string(req.Inputs.Type),
		ExpiresIn:           req.Inputs.ExpiresIn,
		AutoGroups:          req.Inputs.AutoGroups,
		UsageLimit:          req.Inputs.UsageLimit,
		Ephemeral:           req.Inputs.Ephemeral,
		AllowExtraDnsLabels: req.Inputs.AllowExtraDNSLabels,
	}

	setupKey, err := client.SetupKeys.Create(ctx, createReq)
	if err != nil {
		return infer.CreateResponse[SetupKeyState]{}, fmt.Errorf("creating setup key failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:SetupKeyAPI name=%s, id=%s", setupKey.Name, setupKey.Id)

	// Convert time.Time to string
	key := setupKey.Key
	expires := setupKey.Expires.Format("2006-01-02T15:04:05Z07:00")
	lastUsed := setupKey.LastUsed.Format("2006-01-02T15:04:05Z07:00")
	updatedAt := setupKey.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	state := setupKey.State
	revoked := setupKey.Revoked
	usedTimes := setupKey.UsedTimes

	// Note: SetupKey doesn't have a Valid field in the API, using State instead
	valid := state == setupKeyStateValid

	stateObj := SetupKeyState{
		SetupKeyArgs: req.Inputs,
		Key:          &key,
		Valid:        &valid,
		Revoked:      &revoked,
		UsedTimes:    &usedTimes,
		LastUsed:     &lastUsed,
		Expires:      &expires,
		State:        &state,
		UpdatedAt:    &updatedAt,
	}

	return infer.CreateResponse[SetupKeyState]{
		ID:     setupKey.Id,
		Output: stateObj,
	}, nil
}

// Read fetches the current state of a setup key resource from NetBird.
func (*SetupKey) Read(ctx context.Context, setupKeyID string, state SetupKeyState) (SetupKeyState, error) {
	p.GetLogger(ctx).Debugf("Read:SetupKey id=%s", setupKeyID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return state, fmt.Errorf("error getting NetBird client: %w", err)
	}

	setupKey, err := client.SetupKeys.Get(ctx, setupKeyID)
	if err != nil {
		return state, fmt.Errorf("reading setup key failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Read:SetupKeyAPI name=%s, id=%s", setupKey.Name, setupKey.Id)

	state.Name = setupKey.Name
	state.Type = SetupKeyType(setupKey.Type)
	state.AutoGroups = setupKey.AutoGroups
	state.UsageLimit = setupKey.UsageLimit
	ephemeral := setupKey.Ephemeral
	state.Ephemeral = &ephemeral
	allowExtraDNS := setupKey.AllowExtraDnsLabels
	state.AllowExtraDNSLabels = &allowExtraDNS

	// Output fields
	key := setupKey.Key
	state.Key = &key
	revoked := setupKey.Revoked
	state.Revoked = &revoked
	usedTimes := setupKey.UsedTimes
	state.UsedTimes = &usedTimes
	expires := setupKey.Expires.Format("2006-01-02T15:04:05Z07:00")
	state.Expires = &expires
	lastUsed := setupKey.LastUsed.Format("2006-01-02T15:04:05Z07:00")
	state.LastUsed = &lastUsed
	stateStr := setupKey.State
	state.State = &stateStr
	updatedAt := setupKey.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	state.UpdatedAt = &updatedAt
	valid := stateStr == setupKeyStateValid
	state.Valid = &valid

	return state, nil
}

// Update updates the state of the setup key if needed.
func (*SetupKey) Update(ctx context.Context, req infer.UpdateRequest[SetupKeyArgs, SetupKeyState]) (infer.UpdateResponse[SetupKeyState], error) {
	p.GetLogger(ctx).Debugf("Update:SetupKey[%s] name=%s", req.ID, req.Inputs.Name)

	// Check for non-updatable field changes (would require replace)
	if req.Inputs.Name != req.State.Name ||
		req.Inputs.Type != req.State.Type ||
		req.Inputs.ExpiresIn != req.State.ExpiresIn ||
		req.Inputs.UsageLimit != req.State.UsageLimit ||
		boolVal(req.Inputs.Ephemeral) != boolVal(req.State.Ephemeral) ||
		boolVal(req.Inputs.AllowExtraDNSLabels) != boolVal(req.State.AllowExtraDNSLabels) {
		p.GetLogger(ctx).Warningf("Update:SetupKey[%s] non-updatable fields changed, resource needs replacement", req.ID)

		return infer.UpdateResponse[SetupKeyState]{}, errors.New("non-updatable fields changed, resource needs replacement")
	}

	if req.DryRun {
		return infer.UpdateResponse[SetupKeyState]{
			Output: SetupKeyState{
				SetupKeyArgs: req.Inputs,
				Key:          nil,
				Valid:        nil,
				Revoked:      nil,
				UsedTimes:    nil,
				LastUsed:     nil,
				Expires:      nil,
				State:        nil,
				UpdatedAt:    nil,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[SetupKeyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// Only AutoGroups and Revoked can be updated
	updateReq := nbapi.SetupKeyRequest{
		AutoGroups: req.Inputs.AutoGroups,
		Revoked:    req.State.Revoked != nil && *req.State.Revoked,
	}

	updated, err := client.SetupKeys.Update(ctx, req.ID, updateReq)
	if err != nil {
		return infer.UpdateResponse[SetupKeyState]{}, fmt.Errorf("updating setup key failed: %w", err)
	}

	out := req.State
	out.AutoGroups = req.Inputs.AutoGroups
	revoked := updated.Revoked
	out.Revoked = &revoked
	stateStr := updated.State
	out.State = &stateStr
	valid := stateStr == setupKeyStateValid
	out.Valid = &valid
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	out.UpdatedAt = &updatedAt

	return infer.UpdateResponse[SetupKeyState]{Output: out}, nil
}

// Delete removes a setup key from NetBird.
func (*SetupKey) Delete(ctx context.Context, req infer.DeleteRequest[SetupKeyState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:SetupKey[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.SetupKeys.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting setup key failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}
