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

// SetupKeyType defines the kind of setup key accepted by NetBird.
type SetupKeyType string

const (
	// SetupKeyTypeReusable creates a key that can be used multiple times.
	SetupKeyTypeReusable SetupKeyType = SetupKeyType("reusable")
	// SetupKeyTypeOneOff creates a key that can only be used once.
	SetupKeyTypeOneOff SetupKeyType = SetupKeyType("one-off")
)

// SetupKeyArgs represents the input arguments for creating a setup key.
type SetupKeyArgs struct {
	Name                string       `pulumi:"name"`
	Type                SetupKeyType `pulumi:"type"`
	ExpiresIn           int          `pulumi:"expiresIn"` // seconds
	AutoGroups          []string     `pulumi:"autoGroups"`
	UsageLimit          int          `pulumi:"usageLimit"`
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
	state.Ephemeral = &setupKey.Ephemeral
	state.AllowExtraDNSLabels = &setupKey.AllowExtraDnsLabels

	// Output fields
	state.Key = &setupKey.Key
	state.Revoked = &setupKey.Revoked
	state.UsedTimes = &setupKey.UsedTimes
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

// Diff detects changes between inputs and prior state.
func (*SetupKey) Diff(ctx context.Context, req infer.DiffRequest[SetupKeyArgs, SetupKeyState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:SetupKey[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.UpdateReplace,
		}
	}

	if req.Inputs.Type != req.State.Type {
		diff["type"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.UpdateReplace,
		}
	}

	if req.Inputs.ExpiresIn != req.State.ExpiresIn {
		diff["expiresIn"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.UpdateReplace,
		}
	}

	if req.Inputs.UsageLimit != req.State.UsageLimit {
		diff["usageLimit"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.UpdateReplace,
		}
	}

	if boolVal(req.Inputs.Ephemeral) != boolVal(req.State.Ephemeral) {
		diff["ephemeral"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.UpdateReplace,
		}
	}

	if boolVal(req.Inputs.AllowExtraDNSLabels) != boolVal(req.State.AllowExtraDNSLabels) {
		diff["allowExtraDnsLabels"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.UpdateReplace,
		}
	}

	if !equalSlice(req.Inputs.AutoGroups, req.State.AutoGroups) {
		diff["autoGroups"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	p.GetLogger(ctx).Debugf("Diff:SetupKey[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*SetupKey) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[SetupKeyArgs], error) {
	p.GetLogger(ctx).Debugf("Check:SetupKey old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[SetupKeyArgs](ctx, req.NewInputs)
	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{
			Property: "name",
			Reason:   "name must not be empty",
		})
	}

	if args.Type != SetupKeyTypeReusable && args.Type != SetupKeyTypeOneOff {
		failures = append(failures, p.CheckFailure{
			Property: "type",
			Reason:   "type must be 'reusable' or 'one-off'",
		})
	}

	if args.ExpiresIn < 0 {
		failures = append(failures, p.CheckFailure{
			Property: "expiresIn",
			Reason:   "expiresIn must be greater than or equal to 0",
		})
	}

	if args.UsageLimit < 0 {
		failures = append(failures, p.CheckFailure{
			Property: "usageLimit",
			Reason:   "usageLimit must be greater than or equal to 0",
		})
	}

	for i, groupID := range args.AutoGroups {
		if isBlank(groupID) {
			failures = append(failures, p.CheckFailure{
				Property: fmt.Sprintf("autoGroups[%d]", i),
				Reason:   "group id must not be empty",
			})
		}
	}

	return infer.CheckResponse[SetupKeyArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}
