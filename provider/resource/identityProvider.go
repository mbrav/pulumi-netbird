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

// IdentityProvider represents a NetBird identity provider (OIDC) configuration resource.
type IdentityProvider struct{}

// Annotate adds a description to the IdentityProvider resource type.
func (i *IdentityProvider) Annotate(a infer.Annotator) {
	a.Describe(&i, "A NetBird identity provider (OIDC) configuration for self-hosted authentication.")
}

// IdentityProviderType defines the kind of identity provider.
type IdentityProviderType string

const (
	// IdentityProviderTypeADFS is Microsoft AD FS.
	IdentityProviderTypeADFS IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypeAdfs)
	// IdentityProviderTypeEntra is Microsoft Entra ID.
	IdentityProviderTypeEntra IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypeEntra)
	// IdentityProviderTypeGoogle is Google.
	IdentityProviderTypeGoogle IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypeGoogle)
	// IdentityProviderTypeMicrosoft is Microsoft.
	IdentityProviderTypeMicrosoft IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypeMicrosoft)
	// IdentityProviderTypeOIDC is a generic OIDC provider.
	IdentityProviderTypeOIDC IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypeOidc)
	// IdentityProviderTypeOkta is Okta.
	IdentityProviderTypeOkta IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypeOkta)
	// IdentityProviderTypePocketID is PocketID.
	IdentityProviderTypePocketID IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypePocketid)
	// IdentityProviderTypeZitadel is Zitadel.
	IdentityProviderTypeZitadel IdentityProviderType = IdentityProviderType(nbapi.IdentityProviderTypeZitadel)
)

// Values returns the valid enum values for IdentityProviderType.
func (IdentityProviderType) Values() []infer.EnumValue[IdentityProviderType] {
	return []infer.EnumValue[IdentityProviderType]{
		{Name: "adfs", Value: IdentityProviderTypeADFS, Description: "Microsoft AD FS."},
		{Name: "entra", Value: IdentityProviderTypeEntra, Description: "Microsoft Entra ID."},
		{Name: "google", Value: IdentityProviderTypeGoogle, Description: "Google."},
		{Name: "microsoft", Value: IdentityProviderTypeMicrosoft, Description: "Microsoft."},
		{Name: "oidc", Value: IdentityProviderTypeOIDC, Description: "Generic OIDC provider."},
		{Name: "okta", Value: IdentityProviderTypeOkta, Description: "Okta."},
		{Name: "pocketid", Value: IdentityProviderTypePocketID, Description: "PocketID."},
		{Name: "zitadel", Value: IdentityProviderTypeZitadel, Description: "Zitadel."},
	}
}

// IdentityProviderArgs defines input fields for an identity provider.
type IdentityProviderArgs struct {
	Name         string               `pulumi:"name"`
	Type         IdentityProviderType `pulumi:"type"`
	Issuer       string               `pulumi:"issuer"`
	ClientID     string               `pulumi:"clientId"`
	ClientSecret string               `provider:"secret" pulumi:"clientSecret"`
}

// Annotate provides documentation for IdentityProviderArgs fields.
func (i *IdentityProviderArgs) Annotate(a infer.Annotator) {
	a.Describe(&i.Name, "Human-readable name for the identity provider.")
	a.Describe(&i.Type, "Type of identity provider.")
	a.Describe(&i.Issuer, "OIDC issuer URL.")
	a.Describe(&i.ClientID, "OAuth2 client ID.")
	a.Describe(&i.ClientSecret, "OAuth2 client secret.")
}

// IdentityProviderState represents the output state of an identity provider resource.
type IdentityProviderState struct {
	Name         string               `pulumi:"name"`
	Type         IdentityProviderType `pulumi:"type"`
	Issuer       string               `pulumi:"issuer"`
	ClientID     string               `pulumi:"clientId"`
	ClientSecret string               `provider:"secret" pulumi:"clientSecret"`
}

// Annotate provides documentation for IdentityProviderState fields.
func (i *IdentityProviderState) Annotate(a infer.Annotator) {
	a.Describe(&i.Name, "Human-readable name for the identity provider.")
	a.Describe(&i.Type, "Type of identity provider.")
	a.Describe(&i.Issuer, "OIDC issuer URL.")
	a.Describe(&i.ClientID, "OAuth2 client ID.")
	a.Describe(&i.ClientSecret, "OAuth2 client secret. Not returned by the API; preserved from configuration.")
}

// Create creates a new identity provider.
func (*IdentityProvider) Create(ctx context.Context, req infer.CreateRequest[IdentityProviderArgs]) (infer.CreateResponse[IdentityProviderState], error) {
	p.GetLogger(ctx).Debugf("Create:IdentityProvider name=%s type=%s", req.Inputs.Name, req.Inputs.Type)

	if req.DryRun {
		return infer.CreateResponse[IdentityProviderState]{
			ID:     "preview",
			Output: identityProviderStateFromArgs(req.Inputs),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[IdentityProviderState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	idp, err := client.IdentityProviders.Create(ctx, nbapi.IdentityProviderRequest{
		Name:         req.Inputs.Name,
		Type:         nbapi.IdentityProviderType(req.Inputs.Type),
		Issuer:       req.Inputs.Issuer,
		ClientId:     req.Inputs.ClientID,
		ClientSecret: req.Inputs.ClientSecret,
	})
	if err != nil {
		return infer.CreateResponse[IdentityProviderState]{}, fmt.Errorf("creating identity provider failed: %w", err)
	}

	if idp.Id == nil {
		return infer.CreateResponse[IdentityProviderState]{}, errors.New("creating identity provider failed: API returned no ID")
	}

	return infer.CreateResponse[IdentityProviderState]{
		ID:     *idp.Id,
		Output: identityProviderStateFromArgs(req.Inputs),
	}, nil
}

// Read fetches the current state of an identity provider from NetBird.
func (*IdentityProvider) Read(ctx context.Context, req infer.ReadRequest[IdentityProviderArgs, IdentityProviderState]) (infer.ReadResponse[IdentityProviderArgs, IdentityProviderState], error) {
	p.GetLogger(ctx).Debugf("Read:IdentityProvider[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[IdentityProviderArgs, IdentityProviderState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	idp, err := client.IdentityProviders.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[IdentityProviderArgs, IdentityProviderState]{
				ID:     "",
				Inputs: IdentityProviderArgs{},  //nolint:exhaustruct
				State:  IdentityProviderState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[IdentityProviderArgs, IdentityProviderState]{}, fmt.Errorf("reading identity provider failed: %w", err)
	}

	// The API never returns the client secret; preserve the configured value.
	return infer.ReadResponse[IdentityProviderArgs, IdentityProviderState]{
		ID: req.ID,
		Inputs: IdentityProviderArgs{
			Name:         idp.Name,
			Type:         IdentityProviderType(idp.Type),
			Issuer:       idp.Issuer,
			ClientID:     idp.ClientId,
			ClientSecret: req.Inputs.ClientSecret,
		},
		State: IdentityProviderState{
			Name:         idp.Name,
			Type:         IdentityProviderType(idp.Type),
			Issuer:       idp.Issuer,
			ClientID:     idp.ClientId,
			ClientSecret: req.State.ClientSecret,
		},
	}, nil
}

// Update updates an identity provider.
func (*IdentityProvider) Update(ctx context.Context, req infer.UpdateRequest[IdentityProviderArgs, IdentityProviderState]) (infer.UpdateResponse[IdentityProviderState], error) {
	p.GetLogger(ctx).Debugf("Update:IdentityProvider[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[IdentityProviderState]{
			Output: identityProviderStateFromArgs(req.Inputs),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[IdentityProviderState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	_, err = client.IdentityProviders.Update(ctx, req.ID, nbapi.IdentityProviderRequest{
		Name:         req.Inputs.Name,
		Type:         nbapi.IdentityProviderType(req.Inputs.Type),
		Issuer:       req.Inputs.Issuer,
		ClientId:     req.Inputs.ClientID,
		ClientSecret: req.Inputs.ClientSecret,
	})
	if err != nil {
		return infer.UpdateResponse[IdentityProviderState]{}, fmt.Errorf("updating identity provider failed: %w", err)
	}

	return infer.UpdateResponse[IdentityProviderState]{
		Output: identityProviderStateFromArgs(req.Inputs),
	}, nil
}

// Delete removes an identity provider from NetBird.
func (*IdentityProvider) Delete(ctx context.Context, req infer.DeleteRequest[IdentityProviderState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:IdentityProvider[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.IdentityProviders.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting identity provider failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*IdentityProvider) Diff(ctx context.Context, req infer.DiffRequest[IdentityProviderArgs, IdentityProviderState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:IdentityProvider[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Type != req.State.Type {
		diff["type"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Issuer != req.State.Issuer {
		diff["issuer"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.ClientID != req.State.ClientID {
		diff["clientId"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.ClientSecret != req.State.ClientSecret {
		diff["clientSecret"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for an identity provider.
func (*IdentityProvider) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[IdentityProviderArgs], error) {
	p.GetLogger(ctx).Debugf("Check:IdentityProvider old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[IdentityProviderArgs](ctx, req.NewInputs)

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{Property: "name", Reason: "name must not be empty"})
	}

	if isBlank(args.Issuer) {
		failures = append(failures, p.CheckFailure{Property: "issuer", Reason: "issuer must not be empty"})
	}

	if isBlank(args.ClientID) {
		failures = append(failures, p.CheckFailure{Property: "clientId", Reason: "clientId must not be empty"})
	}

	if isBlank(args.ClientSecret) {
		failures = append(failures, p.CheckFailure{Property: "clientSecret", Reason: "clientSecret must not be empty"})
	}

	return infer.CheckResponse[IdentityProviderArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*IdentityProvider) WireDependencies(field infer.FieldSelector, args *IdentityProviderArgs, state *IdentityProviderState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Type).DependsOn(field.InputField(&args.Type))
	field.OutputField(&state.Issuer).DependsOn(field.InputField(&args.Issuer))
	field.OutputField(&state.ClientID).DependsOn(field.InputField(&args.ClientID))
	field.OutputField(&state.ClientSecret).DependsOn(field.InputField(&args.ClientSecret))
}

// identityProviderStateFromArgs mirrors inputs into state (all fields are user-supplied).
func identityProviderStateFromArgs(args IdentityProviderArgs) IdentityProviderState {
	return IdentityProviderState(args)
}
