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

// AzureIDP represents a NetBird Azure AD / Entra ID IdP-sync integration.
type AzureIDP struct{}

// Annotate adds a description to the AzureIDP resource type.
func (a *AzureIDP) Annotate(ann infer.Annotator) {
	ann.Describe(&a, "A NetBird Azure AD (Entra ID) identity-provider sync integration.")
}

// AzureHost selects the Microsoft Graph host domain for an Azure integration.
type AzureHost string

const (
	// AzureHostMicrosoftCom is the commercial Microsoft Graph host.
	AzureHostMicrosoftCom AzureHost = AzureHost(nbapi.CreateAzureIntegrationRequestHostMicrosoftCom)
	// AzureHostMicrosoftUS is the US Government Microsoft Graph host.
	AzureHostMicrosoftUS AzureHost = AzureHost(nbapi.CreateAzureIntegrationRequestHostMicrosoftUs)
)

// Values returns the valid enum values for AzureHost.
func (AzureHost) Values() []infer.EnumValue[AzureHost] {
	return []infer.EnumValue[AzureHost]{
		{Name: "microsoft.com", Value: AzureHostMicrosoftCom, Description: "Commercial Microsoft Graph host."},
		{Name: "microsoft.us", Value: AzureHostMicrosoftUS, Description: "US Government Microsoft Graph host."},
	}
}

// AzureIDPArgs defines input fields for an Azure IdP integration.
type AzureIDPArgs struct {
	ClientID          string    `pulumi:"clientId"`
	ClientSecret      string    `provider:"secret"                   pulumi:"clientSecret"`
	TenantID          string    `pulumi:"tenantId"`
	Host              AzureHost `pulumi:"host"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
	SyncInterval      *int      `pulumi:"syncInterval,optional"`
}

// Annotate provides documentation for AzureIDPArgs fields.
func (a *AzureIDPArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ClientID, "Azure AD application (client) ID.")
	ann.Describe(&a.ClientSecret, "Base64-encoded Azure AD client secret.")
	ann.Describe(&a.TenantID, "Azure AD tenant ID.")
	ann.Describe(&a.Host, "Azure host domain for the Graph API. Changing this forces a replacement.")
	ann.Describe(&a.Enabled, "Whether the integration is enabled.")
	ann.Describe(&a.ConnectorID, "DEX connector ID for embedded IdP setups.")
	ann.Describe(&a.GroupPrefixes, "start_with patterns for groups to sync.")
	ann.Describe(&a.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
	ann.Describe(&a.SyncInterval, "Sync interval in seconds (minimum 300).")
}

// AzureIDPState represents the output state of an Azure IdP integration.
type AzureIDPState struct {
	ClientID          string    `pulumi:"clientId"`
	ClientSecret      string    `provider:"secret"                   pulumi:"clientSecret"`
	TenantID          string    `pulumi:"tenantId"`
	Host              AzureHost `pulumi:"host"`
	Enabled           *bool     `pulumi:"enabled,optional"`
	ConnectorID       *string   `pulumi:"connectorId,optional"`
	GroupPrefixes     *[]string `pulumi:"groupPrefixes,optional"`
	UserGroupPrefixes *[]string `pulumi:"userGroupPrefixes,optional"`
	SyncInterval      *int      `pulumi:"syncInterval,optional"`
	LastSyncedAt      *string   `pulumi:"lastSyncedAt,optional"`
}

// Annotate provides documentation for AzureIDPState fields.
func (a *AzureIDPState) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ClientID, "Azure AD application (client) ID.")
	ann.Describe(&a.ClientSecret, "Base64-encoded Azure AD client secret. Not returned by the API; preserved from configuration.")
	ann.Describe(&a.TenantID, "Azure AD tenant ID.")
	ann.Describe(&a.Host, "Azure host domain for the Graph API.")
	ann.Describe(&a.Enabled, "Whether the integration is enabled.")
	ann.Describe(&a.ConnectorID, "DEX connector ID for embedded IdP setups.")
	ann.Describe(&a.GroupPrefixes, "start_with patterns for groups to sync.")
	ann.Describe(&a.UserGroupPrefixes, "start_with patterns for groups whose users to sync.")
	ann.Describe(&a.SyncInterval, "Sync interval in seconds.")
	ann.Describe(&a.LastSyncedAt, "Timestamp of the last synchronization.")
}

// Create creates a new Azure IdP integration.
func (*AzureIDP) Create(ctx context.Context, req infer.CreateRequest[AzureIDPArgs]) (infer.CreateResponse[AzureIDPState], error) {
	p.GetLogger(ctx).Debugf("Create:AzureIDP clientId=%s tenantId=%s", req.Inputs.ClientID, req.Inputs.TenantID)

	if req.DryRun {
		return infer.CreateResponse[AzureIDPState]{
			ID:     "preview",
			Output: azureIDPStateFromArgs(req.Inputs, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[AzureIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	created, err := client.AzureIDP.Create(ctx, nbapi.CreateAzureIntegrationRequest{
		ClientId:          req.Inputs.ClientID,
		ClientSecret:      req.Inputs.ClientSecret,
		TenantId:          req.Inputs.TenantID,
		Host:              nbapi.CreateAzureIntegrationRequestHost(req.Inputs.Host),
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
		SyncInterval:      req.Inputs.SyncInterval,
	})
	if err != nil {
		return infer.CreateResponse[AzureIDPState]{}, fmt.Errorf("creating Azure IdP integration failed: %w", err)
	}

	return infer.CreateResponse[AzureIDPState]{
		ID:     strconv.FormatInt(created.Id, 10),
		Output: azureIDPStateFromAPI(req.Inputs.ClientSecret, *created),
	}, nil
}

// Read fetches the current state of an Azure IdP integration from NetBird.
func (*AzureIDP) Read(ctx context.Context, req infer.ReadRequest[AzureIDPArgs, AzureIDPState]) (infer.ReadResponse[AzureIDPArgs, AzureIDPState], error) {
	p.GetLogger(ctx).Debugf("Read:AzureIDP[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[AzureIDPArgs, AzureIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	idp, err := client.AzureIDP.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[AzureIDPArgs, AzureIDPState]{
				ID:     "",
				Inputs: AzureIDPArgs{},  //nolint:exhaustruct
				State:  AzureIDPState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[AzureIDPArgs, AzureIDPState]{}, fmt.Errorf("reading Azure IdP integration failed: %w", err)
	}

	return infer.ReadResponse[AzureIDPArgs, AzureIDPState]{
		ID: req.ID,
		Inputs: AzureIDPArgs{
			ClientID:          idp.ClientId,
			ClientSecret:      req.Inputs.ClientSecret,
			TenantID:          idp.TenantId,
			Host:              AzureHost(idp.Host),
			Enabled:           &idp.Enabled,
			ConnectorID:       idp.ConnectorId,
			GroupPrefixes:     &idp.GroupPrefixes,
			UserGroupPrefixes: &idp.UserGroupPrefixes,
			SyncInterval:      &idp.SyncInterval,
		},
		State: azureIDPStateFromAPI(req.State.ClientSecret, *idp),
	}, nil
}

// Update updates an Azure IdP integration.
func (*AzureIDP) Update(ctx context.Context, req infer.UpdateRequest[AzureIDPArgs, AzureIDPState]) (infer.UpdateResponse[AzureIDPState], error) {
	p.GetLogger(ctx).Debugf("Update:AzureIDP[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[AzureIDPState]{
			Output: azureIDPStateFromArgs(req.Inputs, req.State.LastSyncedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AzureIDPState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	clientID := req.Inputs.ClientID
	clientSecret := req.Inputs.ClientSecret
	tenantID := req.Inputs.TenantID

	updated, err := client.AzureIDP.Update(ctx, req.ID, nbapi.UpdateAzureIntegrationRequest{
		ClientId:          &clientID,
		ClientSecret:      &clientSecret,
		TenantId:          &tenantID,
		Enabled:           req.Inputs.Enabled,
		ConnectorId:       req.Inputs.ConnectorID,
		GroupPrefixes:     req.Inputs.GroupPrefixes,
		UserGroupPrefixes: req.Inputs.UserGroupPrefixes,
		SyncInterval:      req.Inputs.SyncInterval,
	})
	if err != nil {
		return infer.UpdateResponse[AzureIDPState]{}, fmt.Errorf("updating Azure IdP integration failed: %w", err)
	}

	return infer.UpdateResponse[AzureIDPState]{
		Output: azureIDPStateFromAPI(req.Inputs.ClientSecret, *updated),
	}, nil
}

// Delete removes an Azure IdP integration from NetBird.
func (*AzureIDP) Delete(ctx context.Context, req infer.DeleteRequest[AzureIDPState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:AzureIDP[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.AzureIDP.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting Azure IdP integration failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*AzureIDP) Diff(ctx context.Context, req infer.DiffRequest[AzureIDPArgs, AzureIDPState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:AzureIDP[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.ClientID != req.State.ClientID {
		diff["clientId"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.ClientSecret != req.State.ClientSecret {
		diff["clientSecret"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.TenantID != req.State.TenantID {
		diff["tenantId"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Host != req.State.Host {
		diff["host"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
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

// Check validates input fields for an Azure IdP integration.
func (*AzureIDP) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[AzureIDPArgs], error) {
	p.GetLogger(ctx).Debugf("Check:AzureIDP old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[AzureIDPArgs](ctx, req.NewInputs)

	if args.Enabled == nil {
		enabled := true
		args.Enabled = &enabled
	}

	if isBlank(args.ClientID) {
		failures = append(failures, p.CheckFailure{Property: "clientId", Reason: "clientId must not be empty"})
	}

	if isBlank(args.ClientSecret) {
		failures = append(failures, p.CheckFailure{Property: "clientSecret", Reason: "clientSecret must not be empty"})
	}

	if isBlank(args.TenantID) {
		failures = append(failures, p.CheckFailure{Property: "tenantId", Reason: "tenantId must not be empty"})
	}

	if args.Host != AzureHostMicrosoftCom && args.Host != AzureHostMicrosoftUS {
		failures = append(failures, p.CheckFailure{Property: "host", Reason: "host must be 'microsoft.com' or 'microsoft.us'"})
	}

	if args.SyncInterval != nil && *args.SyncInterval < 300 {
		failures = append(failures, p.CheckFailure{Property: "syncInterval", Reason: "syncInterval must be at least 300 seconds"})
	}

	return infer.CheckResponse[AzureIDPArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*AzureIDP) WireDependencies(field infer.FieldSelector, args *AzureIDPArgs, state *AzureIDPState) {
	field.OutputField(&state.ClientID).DependsOn(field.InputField(&args.ClientID))
	field.OutputField(&state.ClientSecret).DependsOn(field.InputField(&args.ClientSecret))
	field.OutputField(&state.TenantID).DependsOn(field.InputField(&args.TenantID))
	field.OutputField(&state.Host).DependsOn(field.InputField(&args.Host))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.ConnectorID).DependsOn(field.InputField(&args.ConnectorID))
	field.OutputField(&state.GroupPrefixes).DependsOn(field.InputField(&args.GroupPrefixes))
	field.OutputField(&state.UserGroupPrefixes).DependsOn(field.InputField(&args.UserGroupPrefixes))
	field.OutputField(&state.SyncInterval).DependsOn(field.InputField(&args.SyncInterval))
}

func azureIDPStateFromArgs(args AzureIDPArgs, lastSyncedAt *string) AzureIDPState {
	return AzureIDPState{
		ClientID:          args.ClientID,
		ClientSecret:      args.ClientSecret,
		TenantID:          args.TenantID,
		Host:              args.Host,
		Enabled:           args.Enabled,
		ConnectorID:       args.ConnectorID,
		GroupPrefixes:     args.GroupPrefixes,
		UserGroupPrefixes: args.UserGroupPrefixes,
		SyncInterval:      args.SyncInterval,
		LastSyncedAt:      lastSyncedAt,
	}
}

func azureIDPStateFromAPI(clientSecret string, idp nbapi.AzureIntegration) AzureIDPState {
	lastSynced := idp.LastSyncedAt.Format(idpTimeFormat)

	return AzureIDPState{
		ClientID:          idp.ClientId,
		ClientSecret:      clientSecret,
		TenantID:          idp.TenantId,
		Host:              AzureHost(idp.Host),
		Enabled:           &idp.Enabled,
		ConnectorID:       idp.ConnectorId,
		GroupPrefixes:     &idp.GroupPrefixes,
		UserGroupPrefixes: &idp.UserGroupPrefixes,
		SyncInterval:      &idp.SyncInterval,
		LastSyncedAt:      &lastSynced,
	}
}
