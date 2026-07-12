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

// OktaScimIDP represents a NetBird Okta SCIM IdP-sync integration.
type OktaScimIDP struct{}

// Annotate adds a description to the OktaScimIDP resource type.
func (o *OktaScimIDP) Annotate(a infer.Annotator) {
	a.Describe(&o, "A NetBird Okta SCIM identity-provider sync integration.")
}

// OktaScimIDPArgs defines input fields for an Okta SCIM integration.
type OktaScimIDPArgs struct {
	ConnectionName    string    `pulumi:"connectionName"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
}

// Annotate provides documentation for OktaScimIDPArgs fields.
func (o *OktaScimIDPArgs) Annotate(a infer.Annotator) {
	a.Describe(&o.ConnectionName, "The Okta enterprise connection name on Auth0. Changing this forces a replacement.")
	a.Describe(&o.Enabled, "Whether the integration is enabled.")
	a.Describe(&o.ConnectorID, "DEX connector ID for embedded IdP setups.")
	a.Describe(&o.GroupPrefixes, "start_with patterns for groups to sync.")
	a.Describe(&o.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
}

// OktaScimIDPState represents the output state of an Okta SCIM integration.
type OktaScimIDPState struct {
	ConnectionName    string    `pulumi:"connectionName"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
	AuthToken         *string   `provider:"secret"                   pulumi:"authToken,optional"`
	LastSyncedAt      *string   `pulumi:"lastSyncedAt,optional"`
}

// Annotate provides documentation for OktaScimIDPState fields.
func (o *OktaScimIDPState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&o.ConnectionName, "The Okta enterprise connection name on Auth0.")
	annotator.Describe(&o.Enabled, "Whether the integration is enabled.")
	annotator.Describe(&o.ConnectorID, "DEX connector ID for embedded IdP setups.")
	annotator.Describe(&o.GroupPrefixes, "start_with patterns for groups to sync.")
	annotator.Describe(&o.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
	annotator.Describe(&o.AuthToken, "SCIM API token. Returned in full only on creation; masked afterwards, so the created value is preserved.")
	annotator.Describe(&o.LastSyncedAt, "Timestamp of the last synchronization.")
}

// Create creates a new Okta SCIM integration.
func (*OktaScimIDP) Create(ctx context.Context, req infer.CreateRequest[OktaScimIDPArgs]) (infer.CreateResponse[OktaScimIDPState], error) {
	p.GetLogger(ctx).Debugf("Create:OktaScimIDP connectionName=%s", req.Inputs.ConnectionName)

	if req.DryRun {
		return infer.CreateResponse[OktaScimIDPState]{
			ID:     "preview",
			Output: oktaScimIDPStateFromArgs(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[OktaScimIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	created, err := client.OktaScimIDP.Create(ctx, nbapi.CreateOktaScimIntegrationRequest{
		ConnectionName:    req.Inputs.ConnectionName,
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
	})
	if err != nil {
		return infer.CreateResponse[OktaScimIDPState]{}, fmt.Errorf("creating Okta SCIM integration failed: %w", err)
	}

	authToken := created.AuthToken

	return infer.CreateResponse[OktaScimIDPState]{
		ID:     strconv.FormatInt(created.Id, 10),
		Output: oktaScimIDPStateFromAPI(req.Inputs.ConnectionName, &authToken, *created),
	}, nil
}

// Read fetches the current state of an Okta SCIM integration from NetBird.
func (*OktaScimIDP) Read(ctx context.Context, req infer.ReadRequest[OktaScimIDPArgs, OktaScimIDPState]) (infer.ReadResponse[OktaScimIDPArgs, OktaScimIDPState], error) {
	p.GetLogger(ctx).Debugf("Read:OktaScimIDP[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[OktaScimIDPArgs, OktaScimIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	idp, err := client.OktaScimIDP.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[OktaScimIDPArgs, OktaScimIDPState]{
				ID:     "",
				Inputs: OktaScimIDPArgs{},  //nolint:exhaustruct
				State:  OktaScimIDPState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[OktaScimIDPArgs, OktaScimIDPState]{}, fmt.Errorf("reading Okta SCIM integration failed: %w", err)
	}

	// connectionName is not returned by the API; auth token is masked on read.
	return infer.ReadResponse[OktaScimIDPArgs, OktaScimIDPState]{
		ID: req.ID,
		Inputs: OktaScimIDPArgs{
			ConnectionName:    req.Inputs.ConnectionName,
			Enabled:           &idp.Enabled,
			ConnectorID:       idp.ConnectorId,
			GroupPrefixes:     &idp.GroupPrefixes,
			UserGroupPrefixes: &idp.UserGroupPrefixes,
		},
		State: oktaScimIDPStateFromAPI(req.State.ConnectionName, req.State.AuthToken, *idp),
	}, nil
}

// Update updates an Okta SCIM integration.
func (*OktaScimIDP) Update(ctx context.Context, req infer.UpdateRequest[OktaScimIDPArgs, OktaScimIDPState]) (infer.UpdateResponse[OktaScimIDPState], error) {
	p.GetLogger(ctx).Debugf("Update:OktaScimIDP[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[OktaScimIDPState]{
			Output: oktaScimIDPStateFromArgs(req.Inputs, req.State.AuthToken, req.State.LastSyncedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[OktaScimIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.OktaScimIDP.Update(ctx, req.ID, nbapi.UpdateOktaScimIntegrationRequest{
		Enabled:           req.Inputs.Enabled,
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
	})
	if err != nil {
		return infer.UpdateResponse[OktaScimIDPState]{}, fmt.Errorf("updating Okta SCIM integration failed: %w", err)
	}

	return infer.UpdateResponse[OktaScimIDPState]{
		Output: oktaScimIDPStateFromAPI(req.Inputs.ConnectionName, req.State.AuthToken, *updated),
	}, nil
}

// Delete removes an Okta SCIM integration from NetBird.
func (*OktaScimIDP) Delete(ctx context.Context, req infer.DeleteRequest[OktaScimIDPState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:OktaScimIDP[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.OktaScimIDP.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting Okta SCIM integration failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*OktaScimIDP) Diff(ctx context.Context, req infer.DiffRequest[OktaScimIDPArgs, OktaScimIDPState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:OktaScimIDP[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.ConnectionName != req.State.ConnectionName {
		diff["connectionName"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
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

// Check validates input fields for an Okta SCIM integration.
func (*OktaScimIDP) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[OktaScimIDPArgs], error) {
	p.GetLogger(ctx).Debugf("Check:OktaScimIDP old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[OktaScimIDPArgs](ctx, req.NewInputs)

	if args.Enabled == nil {
		enabled := true
		args.Enabled = &enabled
	}

	if isBlank(args.ConnectionName) {
		failures = append(failures, p.CheckFailure{Property: "connectionName", Reason: "connectionName must not be empty"})
	}

	return infer.CheckResponse[OktaScimIDPArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*OktaScimIDP) WireDependencies(field infer.FieldSelector, args *OktaScimIDPArgs, state *OktaScimIDPState) {
	field.OutputField(&state.ConnectionName).DependsOn(field.InputField(&args.ConnectionName))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.ConnectorID).DependsOn(field.InputField(&args.ConnectorID))
	field.OutputField(&state.GroupPrefixes).DependsOn(field.InputField(&args.GroupPrefixes))
	field.OutputField(&state.UserGroupPrefixes).DependsOn(field.InputField(&args.UserGroupPrefixes))
}

func oktaScimIDPStateFromArgs(args OktaScimIDPArgs, authToken *string, lastSyncedAt *string) OktaScimIDPState {
	return OktaScimIDPState{
		ConnectionName:    args.ConnectionName,
		Enabled:           args.Enabled,
		ConnectorID:       args.ConnectorID,
		GroupPrefixes:     args.GroupPrefixes,
		UserGroupPrefixes: args.UserGroupPrefixes,
		AuthToken:         authToken,
		LastSyncedAt:      lastSyncedAt,
	}
}

func oktaScimIDPStateFromAPI(connectionName string, authToken *string, idp nbapi.OktaScimIntegration) OktaScimIDPState {
	lastSynced := idp.LastSyncedAt.Format(idpTimeFormat)

	return OktaScimIDPState{
		ConnectionName:    connectionName,
		Enabled:           &idp.Enabled,
		ConnectorID:       idp.ConnectorId,
		GroupPrefixes:     &idp.GroupPrefixes,
		UserGroupPrefixes: &idp.UserGroupPrefixes,
		AuthToken:         authToken,
		LastSyncedAt:      &lastSynced,
	}
}
