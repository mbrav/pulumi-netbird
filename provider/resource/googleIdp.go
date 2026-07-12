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

// idpTimeFormat is the timestamp layout used for IdP-integration output fields.
const idpTimeFormat = "2006-01-02T15:04:05Z07:00"

// GoogleIDP represents a NetBird Google Workspace IdP-sync integration.
type GoogleIDP struct{}

// Annotate adds a description to the GoogleIDP resource type.
func (g *GoogleIDP) Annotate(a infer.Annotator) {
	a.Describe(&g, "A NetBird Google Workspace identity-provider sync integration.")
}

// GoogleIDPArgs defines input fields for a Google IdP integration.
type GoogleIDPArgs struct {
	CustomerID        string    `pulumi:"customerId"`
	ServiceAccountKey string    `provider:"secret"                   pulumi:"serviceAccountKey"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
	SyncInterval      *int      `pulumi:"syncInterval,optional"`
}

// Annotate provides documentation for GoogleIDPArgs fields.
func (g *GoogleIDPArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&g.CustomerID, "Customer ID from Google Workspace account settings.")
	annotator.Describe(&g.ServiceAccountKey, "Base64-encoded Google service account key.")
	annotator.Describe(&g.Enabled, "Whether the integration is enabled.")
	annotator.Describe(&g.ConnectorID, "DEX connector ID for embedded IdP setups.")
	annotator.Describe(&g.GroupPrefixes, "start_with patterns for groups to sync.")
	annotator.Describe(&g.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
	annotator.Describe(&g.SyncInterval, "Sync interval in seconds (minimum 300).")
}

// GoogleIDPState represents the output state of a Google IdP integration.
type GoogleIDPState struct {
	CustomerID        string    `pulumi:"customerId"`
	ServiceAccountKey string    `provider:"secret"                   pulumi:"serviceAccountKey"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
	SyncInterval      *int      `pulumi:"syncInterval,optional"`
	LastSyncedAt      *string   `pulumi:"lastSyncedAt,optional"`
}

// Annotate provides documentation for GoogleIDPState fields.
func (g *GoogleIDPState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&g.CustomerID, "Customer ID from Google Workspace account settings.")
	annotator.Describe(&g.ServiceAccountKey, "Base64-encoded Google service account key. Not returned by the API; preserved from configuration.")
	annotator.Describe(&g.Enabled, "Whether the integration is enabled.")
	annotator.Describe(&g.ConnectorID, "DEX connector ID for embedded IdP setups.")
	annotator.Describe(&g.GroupPrefixes, "start_with patterns for groups to sync.")
	annotator.Describe(&g.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
	annotator.Describe(&g.SyncInterval, "Sync interval in seconds.")
	annotator.Describe(&g.LastSyncedAt, "Timestamp of the last synchronization.")
}

// Create creates a new Google IdP integration.
func (*GoogleIDP) Create(ctx context.Context, req infer.CreateRequest[GoogleIDPArgs]) (infer.CreateResponse[GoogleIDPState], error) {
	p.GetLogger(ctx).Debugf("Create:GoogleIDP customerId=%s", req.Inputs.CustomerID)

	if req.DryRun {
		return infer.CreateResponse[GoogleIDPState]{
			ID:     "preview",
			Output: googleIDPStateFromArgs(req.Inputs, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[GoogleIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	created, err := client.GoogleIDP.Create(ctx, nbapi.CreateGoogleIntegrationRequest{
		CustomerId:        req.Inputs.CustomerID,
		ServiceAccountKey: req.Inputs.ServiceAccountKey,
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
		SyncInterval:      req.Inputs.SyncInterval,
	})
	if err != nil {
		return infer.CreateResponse[GoogleIDPState]{}, fmt.Errorf("creating Google IdP integration failed: %w", err)
	}

	return infer.CreateResponse[GoogleIDPState]{
		ID:     strconv.FormatInt(created.Id, 10),
		Output: googleIDPStateFromAPI(req.Inputs.ServiceAccountKey, *created),
	}, nil
}

// Read fetches the current state of a Google IdP integration from NetBird.
func (*GoogleIDP) Read(ctx context.Context, req infer.ReadRequest[GoogleIDPArgs, GoogleIDPState]) (infer.ReadResponse[GoogleIDPArgs, GoogleIDPState], error) {
	p.GetLogger(ctx).Debugf("Read:GoogleIDP[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[GoogleIDPArgs, GoogleIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	idp, err := client.GoogleIDP.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[GoogleIDPArgs, GoogleIDPState]{
				ID:     "",
				Inputs: GoogleIDPArgs{},  //nolint:exhaustruct
				State:  GoogleIDPState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[GoogleIDPArgs, GoogleIDPState]{}, fmt.Errorf("reading Google IdP integration failed: %w", err)
	}

	state := googleIDPStateFromAPI(req.State.ServiceAccountKey, *idp)

	return infer.ReadResponse[GoogleIDPArgs, GoogleIDPState]{
		ID: req.ID,
		Inputs: GoogleIDPArgs{
			CustomerID:        idp.CustomerId,
			ServiceAccountKey: req.Inputs.ServiceAccountKey,
			Enabled:           &idp.Enabled,
			ConnectorID:       idp.ConnectorId,
			GroupPrefixes:     &idp.GroupPrefixes,
			UserGroupPrefixes: &idp.UserGroupPrefixes,
			SyncInterval:      &idp.SyncInterval,
		},
		State: state,
	}, nil
}

// Update updates a Google IdP integration.
func (*GoogleIDP) Update(ctx context.Context, req infer.UpdateRequest[GoogleIDPArgs, GoogleIDPState]) (infer.UpdateResponse[GoogleIDPState], error) {
	p.GetLogger(ctx).Debugf("Update:GoogleIDP[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[GoogleIDPState]{
			Output: googleIDPStateFromArgs(req.Inputs, req.State.LastSyncedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[GoogleIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	customerID := req.Inputs.CustomerID
	serviceAccountKey := req.Inputs.ServiceAccountKey

	updated, err := client.GoogleIDP.Update(ctx, req.ID, nbapi.UpdateGoogleIntegrationRequest{
		CustomerId:        &customerID,
		ServiceAccountKey: &serviceAccountKey,
		Enabled:           req.Inputs.Enabled,
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
		SyncInterval:      req.Inputs.SyncInterval,
	})
	if err != nil {
		return infer.UpdateResponse[GoogleIDPState]{}, fmt.Errorf("updating Google IdP integration failed: %w", err)
	}

	return infer.UpdateResponse[GoogleIDPState]{
		Output: googleIDPStateFromAPI(req.Inputs.ServiceAccountKey, *updated),
	}, nil
}

// Delete removes a Google IdP integration from NetBird.
func (*GoogleIDP) Delete(ctx context.Context, req infer.DeleteRequest[GoogleIDPState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:GoogleIDP[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.GoogleIDP.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting Google IdP integration failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*GoogleIDP) Diff(ctx context.Context, req infer.DiffRequest[GoogleIDPArgs, GoogleIDPState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:GoogleIDP[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.CustomerID != req.State.CustomerID {
		diff["customerId"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.ServiceAccountKey != req.State.ServiceAccountKey {
		diff["serviceAccountKey"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
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

	if !equalPtr(req.Inputs.SyncInterval, req.State.SyncInterval) {
		diff["syncInterval"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a Google IdP integration.
func (*GoogleIDP) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[GoogleIDPArgs], error) {
	p.GetLogger(ctx).Debugf("Check:GoogleIDP old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[GoogleIDPArgs](ctx, req.NewInputs)

	if args.Enabled == nil {
		enabled := true
		args.Enabled = &enabled
	}

	if isBlank(args.CustomerID) {
		failures = append(failures, p.CheckFailure{Property: "customerId", Reason: "customerId must not be empty"})
	}

	if isBlank(args.ServiceAccountKey) {
		failures = append(failures, p.CheckFailure{Property: "serviceAccountKey", Reason: "serviceAccountKey must not be empty"})
	}

	if args.SyncInterval != nil && *args.SyncInterval < 300 {
		failures = append(failures, p.CheckFailure{Property: "syncInterval", Reason: "syncInterval must be at least 300 seconds"})
	}

	return infer.CheckResponse[GoogleIDPArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*GoogleIDP) WireDependencies(field infer.FieldSelector, args *GoogleIDPArgs, state *GoogleIDPState) {
	field.OutputField(&state.CustomerID).DependsOn(field.InputField(&args.CustomerID))
	field.OutputField(&state.ServiceAccountKey).DependsOn(field.InputField(&args.ServiceAccountKey))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.ConnectorID).DependsOn(field.InputField(&args.ConnectorID))
	field.OutputField(&state.GroupPrefixes).DependsOn(field.InputField(&args.GroupPrefixes))
	field.OutputField(&state.UserGroupPrefixes).DependsOn(field.InputField(&args.UserGroupPrefixes))
	field.OutputField(&state.SyncInterval).DependsOn(field.InputField(&args.SyncInterval))
}

func googleIDPStateFromArgs(args GoogleIDPArgs, lastSyncedAt *string) GoogleIDPState {
	return GoogleIDPState{
		CustomerID:        args.CustomerID,
		ServiceAccountKey: args.ServiceAccountKey,
		Enabled:           args.Enabled,
		ConnectorID:       args.ConnectorID,
		GroupPrefixes:     args.GroupPrefixes,
		UserGroupPrefixes: args.UserGroupPrefixes,
		SyncInterval:      args.SyncInterval,
		LastSyncedAt:      lastSyncedAt,
	}
}

func googleIDPStateFromAPI(serviceAccountKey string, idp nbapi.GoogleIntegration) GoogleIDPState {
	lastSynced := idp.LastSyncedAt.Format(idpTimeFormat)

	return GoogleIDPState{
		CustomerID:        idp.CustomerId,
		ServiceAccountKey: serviceAccountKey,
		Enabled:           &idp.Enabled,
		ConnectorID:       idp.ConnectorId,
		GroupPrefixes:     &idp.GroupPrefixes,
		UserGroupPrefixes: &idp.UserGroupPrefixes,
		SyncInterval:      &idp.SyncInterval,
		LastSyncedAt:      &lastSynced,
	}
}
