package resource

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ScimIntegration represents a NetBird generic SCIM identity-provider integration.
type ScimIntegration struct{}

// Annotate adds a description to the ScimIntegration resource type.
func (s *ScimIntegration) Annotate(a infer.Annotator) {
	a.Describe(&s, "A NetBird generic SCIM identity-provider integration.")
}

// ScimIntegrationArgs defines input fields for a SCIM integration.
type ScimIntegrationArgs struct {
	Prefix            string    `pulumi:"prefix"`
	Provider          string    `pulumi:"provider"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
}

// Annotate provides documentation for ScimIntegrationArgs fields.
func (s *ScimIntegrationArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&s.Prefix, "The connection prefix used for the SCIM provider.")
	annotator.Describe(&s.Provider, "Name of the SCIM identity provider. Changing this forces a replacement.")
	annotator.Describe(&s.Enabled, "Whether the integration is enabled.")
	annotator.Describe(&s.ConnectorID, "DEX connector ID for embedded IdP setups.")
	annotator.Describe(&s.GroupPrefixes, "start_with patterns for groups to sync.")
	annotator.Describe(&s.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
}

// ScimIntegrationState represents the output state of a SCIM integration.
type ScimIntegrationState struct {
	Prefix            string    `pulumi:"prefix"`
	Provider          string    `pulumi:"provider"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
	AuthToken         *string   `provider:"secret"                   pulumi:"authToken,optional"`
	LastSyncedAt      *string   `pulumi:"lastSyncedAt,optional"`
}

// Annotate provides documentation for ScimIntegrationState fields.
func (s *ScimIntegrationState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&s.Prefix, "The connection prefix used for the SCIM provider.")
	annotator.Describe(&s.Provider, "Name of the SCIM identity provider.")
	annotator.Describe(&s.Enabled, "Whether the integration is enabled.")
	annotator.Describe(&s.ConnectorID, "DEX connector ID for embedded IdP setups.")
	annotator.Describe(&s.GroupPrefixes, "start_with patterns for groups to sync.")
	annotator.Describe(&s.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
	annotator.Describe(&s.AuthToken, "SCIM API token. Returned in full only on creation; masked afterwards, so the created value is preserved.")
	annotator.Describe(&s.LastSyncedAt, "Timestamp of the last synchronization.")
}

// Create creates a new SCIM integration.
func (*ScimIntegration) Create(ctx context.Context, req infer.CreateRequest[ScimIntegrationArgs]) (infer.CreateResponse[ScimIntegrationState], error) {
	p.GetLogger(ctx).Debugf("Create:ScimIntegration provider=%s prefix=%s", req.Inputs.Provider, req.Inputs.Prefix)

	if req.DryRun {
		return infer.CreateResponse[ScimIntegrationState]{
			ID:     "preview",
			Output: scimIntegrationStateFromArgs(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[ScimIntegrationState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	created, err := client.SCIM.Create(ctx, nbapi.CreateScimIntegrationRequest{
		Prefix:            req.Inputs.Prefix,
		Provider:          req.Inputs.Provider,
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
	})
	if err != nil {
		return infer.CreateResponse[ScimIntegrationState]{}, fmt.Errorf("creating SCIM integration failed: %w", err)
	}

	authToken := created.AuthToken

	return infer.CreateResponse[ScimIntegrationState]{
		ID:     strconv.FormatInt(created.Id, 10),
		Output: scimIntegrationStateFromAPI(&authToken, *created),
	}, nil
}

// Read fetches the current state of a SCIM integration from NetBird.
func (*ScimIntegration) Read(ctx context.Context, req infer.ReadRequest[ScimIntegrationArgs, ScimIntegrationState]) (infer.ReadResponse[ScimIntegrationArgs, ScimIntegrationState], error) {
	p.GetLogger(ctx).Debugf("Read:ScimIntegration[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[ScimIntegrationArgs, ScimIntegrationState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	idp, err := client.SCIM.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[ScimIntegrationArgs, ScimIntegrationState]{
				ID:     "",
				Inputs: ScimIntegrationArgs{},  //nolint:exhaustruct
				State:  ScimIntegrationState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[ScimIntegrationArgs, ScimIntegrationState]{}, fmt.Errorf("reading SCIM integration failed: %w", err)
	}

	// The auth token is masked on read; preserve the created value.
	return infer.ReadResponse[ScimIntegrationArgs, ScimIntegrationState]{
		ID: req.ID,
		Inputs: ScimIntegrationArgs{
			Prefix:            idp.Prefix,
			Provider:          idp.Provider,
			Enabled:           &idp.Enabled,
			ConnectorID:       idp.ConnectorId,
			GroupPrefixes:     &idp.GroupPrefixes,
			UserGroupPrefixes: &idp.UserGroupPrefixes,
		},
		State: scimIntegrationStateFromAPI(req.State.AuthToken, *idp),
	}, nil
}

// Update updates a SCIM integration.
func (*ScimIntegration) Update(ctx context.Context, req infer.UpdateRequest[ScimIntegrationArgs, ScimIntegrationState]) (infer.UpdateResponse[ScimIntegrationState], error) {
	p.GetLogger(ctx).Debugf("Update:ScimIntegration[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[ScimIntegrationState]{
			Output: scimIntegrationStateFromArgs(req.Inputs, req.State.AuthToken, req.State.LastSyncedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[ScimIntegrationState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	prefix := req.Inputs.Prefix

	updated, err := client.SCIM.Update(ctx, req.ID, nbapi.UpdateScimIntegrationRequest{
		Prefix:            &prefix,
		Enabled:           req.Inputs.Enabled,
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
	})
	if err != nil {
		return infer.UpdateResponse[ScimIntegrationState]{}, fmt.Errorf("updating SCIM integration failed: %w", err)
	}

	return infer.UpdateResponse[ScimIntegrationState]{
		Output: scimIntegrationStateFromAPI(req.State.AuthToken, *updated),
	}, nil
}

// Delete removes a SCIM integration from NetBird.
func (*ScimIntegration) Delete(ctx context.Context, req infer.DeleteRequest[ScimIntegrationState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:ScimIntegration[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.SCIM.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting SCIM integration failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*ScimIntegration) Diff(ctx context.Context, req infer.DiffRequest[ScimIntegrationArgs, ScimIntegrationState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:ScimIntegration[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Prefix != req.State.Prefix {
		diff["prefix"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Provider != req.State.Provider {
		diff["provider"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if boolVal(req.Inputs.Enabled) != boolVal(req.State.Enabled) {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.ConnectorID, req.State.ConnectorID) {
		diff["connectorId"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.GroupPrefixes, req.State.GroupPrefixes) {
		diff["groupPrefixes"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.UserGroupPrefixes, req.State.UserGroupPrefixes) {
		diff["userGroupPrefixes"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a SCIM integration.
func (*ScimIntegration) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[ScimIntegrationArgs], error) {
	p.GetLogger(ctx).Debugf("Check:ScimIntegration old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[ScimIntegrationArgs](ctx, req.NewInputs)

	if args.Enabled == nil {
		enabled := true
		args.Enabled = &enabled
	}

	if isBlank(args.Prefix) {
		failures = append(failures, p.CheckFailure{Property: "prefix", Reason: "prefix must not be empty"})
	}

	if isBlank(args.Provider) {
		failures = append(failures, p.CheckFailure{Property: "provider", Reason: "provider must not be empty"})
	}

	return infer.CheckResponse[ScimIntegrationArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*ScimIntegration) WireDependencies(field infer.FieldSelector, args *ScimIntegrationArgs, state *ScimIntegrationState) {
	field.OutputField(&state.Prefix).DependsOn(field.InputField(&args.Prefix))
	field.OutputField(&state.Provider).DependsOn(field.InputField(&args.Provider))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.ConnectorID).DependsOn(field.InputField(&args.ConnectorID))
	field.OutputField(&state.GroupPrefixes).DependsOn(field.InputField(&args.GroupPrefixes))
	field.OutputField(&state.UserGroupPrefixes).DependsOn(field.InputField(&args.UserGroupPrefixes))
}

func scimIntegrationStateFromArgs(args ScimIntegrationArgs, authToken *string, lastSyncedAt *string) ScimIntegrationState {
	return ScimIntegrationState{
		Prefix:            args.Prefix,
		Provider:          args.Provider,
		Enabled:           args.Enabled,
		ConnectorID:       args.ConnectorID,
		GroupPrefixes:     args.GroupPrefixes,
		UserGroupPrefixes: args.UserGroupPrefixes,
		AuthToken:         authToken,
		LastSyncedAt:      lastSyncedAt,
	}
}

func scimIntegrationStateFromAPI(authToken *string, idp nbapi.ScimIntegration) ScimIntegrationState {
	lastSynced := idp.LastSyncedAt.Format(idpTimeFormat)

	return ScimIntegrationState{
		Prefix:            idp.Prefix,
		Provider:          idp.Provider,
		Enabled:           &idp.Enabled,
		ConnectorID:       idp.ConnectorId,
		GroupPrefixes:     &idp.GroupPrefixes,
		UserGroupPrefixes: &idp.UserGroupPrefixes,
		AuthToken:         authToken,
		LastSyncedAt:      &lastSynced,
	}
}
