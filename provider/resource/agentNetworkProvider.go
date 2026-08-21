package resource

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// AgentNetworkProvider represents a NetBird Agent Network (AI/LLM gateway) upstream provider.
type AgentNetworkProvider struct{}

// Annotate adds a description to the AgentNetworkProvider resource type.
func (a *AgentNetworkProvider) Annotate(ann infer.Annotator) {
	ann.Describe(&a, "A NetBird Agent Network provider: an upstream AI/LLM API or gateway "+
		"(OpenAI, Anthropic, Bedrock, LiteLLM, a custom OpenAI-compatible endpoint, ...) that "+
		"policies can route traffic to.")
}

// AgentNetworkProviderModelPricing overrides the catalog per-1k token price for one model exposed by a provider.
type AgentNetworkProviderModelPricing struct {
	ID                 string   `pulumi:"id"`
	InputPer1k         float64  `pulumi:"inputPer1k"`
	OutputPer1k        float64  `pulumi:"outputPer1k"`
	CacheReadPer1k     *float64 `pulumi:"cacheReadPer1k,optional"`
	CacheCreationPer1k *float64 `pulumi:"cacheCreationPer1k,optional"`
	CachedInputPer1k   *float64 `pulumi:"cachedInputPer1k,optional"`
}

// Annotate provides documentation for AgentNetworkProviderModelPricing fields.
func (m *AgentNetworkProviderModelPricing) Annotate(annotator infer.Annotator) {
	annotator.Describe(&m.ID, "Catalog model identifier (e.g. \"gpt-4o-mini\").")
	annotator.Describe(&m.InputPer1k, "Cost per 1k input tokens, in USD.")
	annotator.Describe(&m.OutputPer1k, "Cost per 1k output tokens, in USD.")
	annotator.Describe(&m.CacheReadPer1k, "Anthropic-shape cache rate: cost per 1k cache-read tokens, in USD.")
	annotator.Describe(&m.CacheCreationPer1k, "Anthropic-shape cache rate: cost per 1k cache-creation tokens, in USD.")
	annotator.Describe(&m.CachedInputPer1k, "OpenAI-shape cache rate: cost per 1k cached prompt tokens, in USD.")
}

// AgentNetworkProviderArgs defines input fields for an Agent Network provider.
type AgentNetworkProviderArgs struct {
	Name                 string                              `pulumi:"name"`
	ProviderID           string                              `pulumi:"providerId"`
	UpstreamURL          string                              `pulumi:"upstreamUrl"`
	APIKey               string                              `provider:"secret"                      pulumi:"apiKey"`
	Enabled              *bool                               `pulumi:"enabled,optional"`
	SkipTLSVerification  *bool                               `pulumi:"skipTlsVerification,optional"`
	MetadataDisabled     *bool                               `pulumi:"metadataDisabled,optional"`
	IdentityHeaderUserID *string                             `pulumi:"identityHeaderUserId,optional"`
	IdentityHeaderGroups *string                             `pulumi:"identityHeaderGroups,optional"`
	ExtraValues          *map[string]string                  `pulumi:"extraValues,optional"`
	Models               *[]AgentNetworkProviderModelPricing `pulumi:"models,optional"`
	BootstrapCluster     *string                             `pulumi:"bootstrapCluster,optional"`
}

// Annotate provides documentation for AgentNetworkProviderArgs fields.
func (a *AgentNetworkProviderArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name shown in the dashboard.")
	ann.Describe(&a.ProviderID, "Catalog identifier for the upstream AI provider (e.g. openai_api, "+
		"anthropic_api, azure_openai_api, bedrock_api, vertex_ai_api, mistral_api, custom). See "+
		"getAgentNetworkCatalogProviders. Changing this forces a replacement.")
	ann.Describe(&a.UpstreamURL, "Full upstream URL (with scheme) that NetBird forwards traffic to.")
	ann.Describe(&a.APIKey, "Upstream provider API key. Sealed at rest on the management server and never returned in responses.")
	ann.Describe(&a.Enabled, "Whether the provider is enabled.")
	ann.Describe(&a.SkipTLSVerification, "Skip upstream TLS certificate verification. For self-hosted / internal gateways with a private or self-signed certificate.")
	ann.Describe(&a.MetadataDisabled, "Disable identity metadata injection (caller's user + authorizing group) for this provider.")
	ann.Describe(&a.IdentityHeaderUserID, "Wire header name the proxy stamps with the caller's display identity, when the catalog entry supports customizable identity headers.")
	ann.Describe(&a.IdentityHeaderGroups, "Wire header name the proxy stamps with the caller's NetBird groups (comma-separated), when the catalog entry supports customizable identity headers.")
	ann.Describe(&a.ExtraValues, "Operator-typed values for catalog-declared extra headers (see getAgentNetworkCatalogProviders).")
	ann.Describe(&a.Models, "Models exposed through this endpoint, with operator per-1k price overrides. Empty means all catalog models are allowed at catalog prices.")
	ann.Describe(&a.BootstrapCluster, "Deprecated and ignored. NetBird removed bootstrap_cluster from the provider API; "+
		"bootstrap the account's gateway endpoint with netbird:resource:AgentNetworkSettings (proxyAddress or endpoint) instead.")
	ann.Deprecate(&a.BootstrapCluster, "bootstrapCluster is ignored. Bootstrap the account gateway with AgentNetworkSettings (proxyAddress or endpoint).")
}

// AgentNetworkProviderState represents the output state of an Agent Network provider.
type AgentNetworkProviderState struct {
	Name                 string                              `pulumi:"name"`
	ProviderID           string                              `pulumi:"providerId"`
	UpstreamURL          string                              `pulumi:"upstreamUrl"`
	APIKey               string                              `provider:"secret"                      pulumi:"apiKey"`
	Enabled              *bool                               `pulumi:"enabled,optional"`
	SkipTLSVerification  *bool                               `pulumi:"skipTlsVerification,optional"`
	MetadataDisabled     *bool                               `pulumi:"metadataDisabled,optional"`
	IdentityHeaderUserID *string                             `pulumi:"identityHeaderUserId,optional"`
	IdentityHeaderGroups *string                             `pulumi:"identityHeaderGroups,optional"`
	ExtraValues          *map[string]string                  `pulumi:"extraValues,optional"`
	Models               *[]AgentNetworkProviderModelPricing `pulumi:"models,optional"`
	CreatedAt            *string                             `pulumi:"createdAt,optional"`
	UpdatedAt            *string                             `pulumi:"updatedAt,optional"`
}

// Annotate provides documentation for AgentNetworkProviderState fields.
func (a *AgentNetworkProviderState) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name shown in the dashboard.")
	ann.Describe(&a.ProviderID, "Catalog identifier for the upstream AI provider.")
	ann.Describe(&a.UpstreamURL, "Full upstream URL (with scheme) that NetBird forwards traffic to.")
	ann.Describe(&a.APIKey, "Upstream provider API key. Not returned by the API; preserved from configuration.")
	ann.Describe(&a.Enabled, "Whether the provider is enabled.")
	ann.Describe(&a.SkipTLSVerification, "Whether upstream TLS certificate verification is skipped.")
	ann.Describe(&a.MetadataDisabled, "Whether identity metadata injection is disabled for this provider.")
	ann.Describe(&a.IdentityHeaderUserID, "Wire header name for the caller's display identity.")
	ann.Describe(&a.IdentityHeaderGroups, "Wire header name for the caller's NetBird groups.")
	ann.Describe(&a.ExtraValues, "Operator-typed values for catalog-declared extra headers.")
	ann.Describe(&a.Models, "Models exposed through this endpoint, with operator per-1k price overrides.")
	ann.Describe(&a.CreatedAt, "Timestamp when the provider was created.")
	ann.Describe(&a.UpdatedAt, "Timestamp when the provider was last updated.")
}

// Create creates a new Agent Network provider.
func (*AgentNetworkProvider) Create(
	ctx context.Context, req infer.CreateRequest[AgentNetworkProviderArgs],
) (infer.CreateResponse[AgentNetworkProviderState], error) {
	p.GetLogger(ctx).Debugf("Create:AgentNetworkProvider name=%s providerId=%s", req.Inputs.Name, req.Inputs.ProviderID)

	if req.DryRun {
		return infer.CreateResponse[AgentNetworkProviderState]{
			ID:     "preview",
			Output: agentNetworkProviderStateFromArgs(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[AgentNetworkProviderState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiKey := req.Inputs.APIKey
	// The identity headers always travel explicitly: the API treats an omitted
	// header as "keep" and an empty one as "do not stamp", and the response
	// always echoes both. Sending "" for an unset input keeps inputs and state
	// in agreement instead of diffing forever, and lets removing the field from
	// the configuration actually clear the header.
	identityHeaderUserID := strPtr(req.Inputs.IdentityHeaderUserID)
	identityHeaderGroups := strPtr(req.Inputs.IdentityHeaderGroups)

	created, err := client.AgentNetwork.CreateProvider(ctx, nbapi.AgentNetworkProviderRequest{
		Name:                 req.Inputs.Name,
		ProviderId:           req.Inputs.ProviderID,
		UpstreamUrl:          req.Inputs.UpstreamURL,
		ApiKey:               &apiKey,
		Enabled:              req.Inputs.Enabled,
		SkipTlsVerification:  req.Inputs.SkipTLSVerification,
		MetadataDisabled:     req.Inputs.MetadataDisabled,
		IdentityHeaderUserId: &identityHeaderUserID,
		IdentityHeaderGroups: &identityHeaderGroups,
		ExtraValues:          req.Inputs.ExtraValues,
		Models:               toAPIAgentNetworkProviderModels(req.Inputs.Models),
	})
	if err != nil {
		return infer.CreateResponse[AgentNetworkProviderState]{}, fmt.Errorf("creating agent network provider failed: %w", err)
	}

	return infer.CreateResponse[AgentNetworkProviderState]{
		ID:     created.Id,
		Output: agentNetworkProviderStateFromAPI(req.Inputs.APIKey, *created),
	}, nil
}

// Read fetches the current state of an Agent Network provider from NetBird.
func (*AgentNetworkProvider) Read(
	ctx context.Context, req infer.ReadRequest[AgentNetworkProviderArgs, AgentNetworkProviderState],
) (infer.ReadResponse[AgentNetworkProviderArgs, AgentNetworkProviderState], error) {
	p.GetLogger(ctx).Debugf("Read:AgentNetworkProvider[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[AgentNetworkProviderArgs, AgentNetworkProviderState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	provider, err := client.AgentNetwork.GetProvider(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[AgentNetworkProviderArgs, AgentNetworkProviderState]{
				ID:     "",
				Inputs: AgentNetworkProviderArgs{},  //nolint:exhaustruct
				State:  AgentNetworkProviderState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[AgentNetworkProviderArgs, AgentNetworkProviderState]{}, fmt.Errorf("reading agent network provider failed: %w", err)
	}

	state := agentNetworkProviderStateFromAPI(req.State.APIKey, *provider)

	return infer.ReadResponse[AgentNetworkProviderArgs, AgentNetworkProviderState]{
		ID: req.ID,
		Inputs: AgentNetworkProviderArgs{
			Name:                 provider.Name,
			ProviderID:           provider.ProviderId,
			UpstreamURL:          provider.UpstreamUrl,
			APIKey:               req.Inputs.APIKey,
			Enabled:              &provider.Enabled,
			SkipTLSVerification:  &provider.SkipTlsVerification,
			MetadataDisabled:     &provider.MetadataDisabled,
			IdentityHeaderUserID: &provider.IdentityHeaderUserId,
			IdentityHeaderGroups: &provider.IdentityHeaderGroups,
			ExtraValues:          provider.ExtraValues,
			Models:               fromAPIAgentNetworkProviderModels(provider.Models),
			BootstrapCluster:     req.Inputs.BootstrapCluster,
		},
		State: state,
	}, nil
}

// Update updates an Agent Network provider.
func (*AgentNetworkProvider) Update(
	ctx context.Context, req infer.UpdateRequest[AgentNetworkProviderArgs, AgentNetworkProviderState],
) (infer.UpdateResponse[AgentNetworkProviderState], error) {
	p.GetLogger(ctx).Debugf("Update:AgentNetworkProvider[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[AgentNetworkProviderState]{
			Output: agentNetworkProviderStateFromArgs(req.Inputs, req.State.CreatedAt, req.State.UpdatedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AgentNetworkProviderState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiKey := req.Inputs.APIKey
	identityHeaderUserID := strPtr(req.Inputs.IdentityHeaderUserID)
	identityHeaderGroups := strPtr(req.Inputs.IdentityHeaderGroups)

	updated, err := client.AgentNetwork.UpdateProvider(ctx, req.ID, nbapi.AgentNetworkProviderRequest{
		Name:                 req.Inputs.Name,
		ProviderId:           req.Inputs.ProviderID,
		UpstreamUrl:          req.Inputs.UpstreamURL,
		ApiKey:               &apiKey,
		Enabled:              req.Inputs.Enabled,
		SkipTlsVerification:  req.Inputs.SkipTLSVerification,
		MetadataDisabled:     req.Inputs.MetadataDisabled,
		IdentityHeaderUserId: &identityHeaderUserID,
		IdentityHeaderGroups: &identityHeaderGroups,
		ExtraValues:          req.Inputs.ExtraValues,
		Models:               toAPIAgentNetworkProviderModels(req.Inputs.Models),
	})
	if err != nil {
		return infer.UpdateResponse[AgentNetworkProviderState]{}, fmt.Errorf("updating agent network provider failed: %w", err)
	}

	return infer.UpdateResponse[AgentNetworkProviderState]{
		Output: agentNetworkProviderStateFromAPI(req.Inputs.APIKey, *updated),
	}, nil
}

// Delete removes an Agent Network provider from NetBird.
func (*AgentNetworkProvider) Delete(ctx context.Context, req infer.DeleteRequest[AgentNetworkProviderState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:AgentNetworkProvider[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.AgentNetwork.DeleteProvider(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting agent network provider failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*AgentNetworkProvider) Diff(
	ctx context.Context, req infer.DiffRequest[AgentNetworkProviderArgs, AgentNetworkProviderState],
) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:AgentNetworkProvider[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.ProviderID != req.State.ProviderID {
		diff["providerId"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.UpstreamURL != req.State.UpstreamURL {
		diff["upstreamUrl"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.APIKey != req.State.APIKey {
		diff["apiKey"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if boolVal(req.Inputs.Enabled) != boolVal(req.State.Enabled) {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if boolVal(req.Inputs.SkipTLSVerification) != boolVal(req.State.SkipTLSVerification) {
		diff["skipTlsVerification"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if boolVal(req.Inputs.MetadataDisabled) != boolVal(req.State.MetadataDisabled) {
		diff["metadataDisabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalOptionalStr(req.Inputs.IdentityHeaderUserID, req.State.IdentityHeaderUserID) {
		diff["identityHeaderUserId"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalOptionalStr(req.Inputs.IdentityHeaderGroups, req.State.IdentityHeaderGroups) {
		diff["identityHeaderGroups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalStringMapPtr(req.Inputs.ExtraValues, req.State.ExtraValues) {
		diff["extraValues"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalAgentNetworkProviderModels(req.Inputs.Models, req.State.Models) {
		diff["models"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:AgentNetworkProvider[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*AgentNetworkProvider) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[AgentNetworkProviderArgs], error) {
	p.GetLogger(ctx).Debugf("Check:AgentNetworkProvider old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[AgentNetworkProviderArgs](ctx, req.NewInputs)

	if args.Enabled == nil {
		enabled := true
		args.Enabled = &enabled
	}

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{Property: "name", Reason: "name must not be empty"})
	}

	if isBlank(args.ProviderID) {
		failures = append(failures, p.CheckFailure{Property: "providerId", Reason: "providerId must not be empty"})
	}

	if isBlank(args.UpstreamURL) {
		failures = append(failures, p.CheckFailure{Property: "upstreamUrl", Reason: "upstreamUrl must not be empty"})
	}

	if isBlank(args.APIKey) {
		failures = append(failures, p.CheckFailure{Property: "apiKey", Reason: "apiKey must not be empty"})
	}

	if args.Models != nil {
		for i, model := range *args.Models {
			if isBlank(model.ID) {
				failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("models[%d].id", i), Reason: "model id must not be empty"})
			}
		}
	}

	return infer.CheckResponse[AgentNetworkProviderArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*AgentNetworkProvider) WireDependencies(field infer.FieldSelector, args *AgentNetworkProviderArgs, state *AgentNetworkProviderState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.ProviderID).DependsOn(field.InputField(&args.ProviderID))
	field.OutputField(&state.UpstreamURL).DependsOn(field.InputField(&args.UpstreamURL))
	field.OutputField(&state.APIKey).DependsOn(field.InputField(&args.APIKey))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.SkipTLSVerification).DependsOn(field.InputField(&args.SkipTLSVerification))
	field.OutputField(&state.MetadataDisabled).DependsOn(field.InputField(&args.MetadataDisabled))
	field.OutputField(&state.IdentityHeaderUserID).DependsOn(field.InputField(&args.IdentityHeaderUserID))
	field.OutputField(&state.IdentityHeaderGroups).DependsOn(field.InputField(&args.IdentityHeaderGroups))
	field.OutputField(&state.ExtraValues).DependsOn(field.InputField(&args.ExtraValues))
	field.OutputField(&state.Models).DependsOn(field.InputField(&args.Models))
}

func agentNetworkProviderStateFromArgs(args AgentNetworkProviderArgs, createdAt, updatedAt *string) AgentNetworkProviderState {
	return AgentNetworkProviderState{
		Name:                 args.Name,
		ProviderID:           args.ProviderID,
		UpstreamURL:          args.UpstreamURL,
		APIKey:               args.APIKey,
		Enabled:              args.Enabled,
		SkipTLSVerification:  args.SkipTLSVerification,
		MetadataDisabled:     args.MetadataDisabled,
		IdentityHeaderUserID: args.IdentityHeaderUserID,
		IdentityHeaderGroups: args.IdentityHeaderGroups,
		ExtraValues:          args.ExtraValues,
		Models:               args.Models,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}

func agentNetworkProviderStateFromAPI(apiKey string, provider nbapi.AgentNetworkProvider) AgentNetworkProviderState {
	var createdAt, updatedAt *string

	if provider.CreatedAt != nil {
		formatted := provider.CreatedAt.Format(idpTimeFormat)
		createdAt = &formatted
	}

	if provider.UpdatedAt != nil {
		formatted := provider.UpdatedAt.Format(idpTimeFormat)
		updatedAt = &formatted
	}

	return AgentNetworkProviderState{
		Name:                 provider.Name,
		ProviderID:           provider.ProviderId,
		UpstreamURL:          provider.UpstreamUrl,
		APIKey:               apiKey,
		Enabled:              &provider.Enabled,
		SkipTLSVerification:  &provider.SkipTlsVerification,
		MetadataDisabled:     &provider.MetadataDisabled,
		IdentityHeaderUserID: &provider.IdentityHeaderUserId,
		IdentityHeaderGroups: &provider.IdentityHeaderGroups,
		ExtraValues:          provider.ExtraValues,
		Models:               fromAPIAgentNetworkProviderModels(provider.Models),
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}

func toAPIAgentNetworkProviderModels(models *[]AgentNetworkProviderModelPricing) *[]nbapi.AgentNetworkProviderModel {
	if models == nil {
		return nil
	}

	out := make([]nbapi.AgentNetworkProviderModel, len(*models))

	for i, model := range *models {
		out[i] = nbapi.AgentNetworkProviderModel{
			Id:                 model.ID,
			InputPer1k:         model.InputPer1k,
			OutputPer1k:        model.OutputPer1k,
			CacheReadPer1k:     model.CacheReadPer1k,
			CacheCreationPer1k: model.CacheCreationPer1k,
			CachedInputPer1k:   model.CachedInputPer1k,
		}
	}

	return &out
}

func fromAPIAgentNetworkProviderModels(models []nbapi.AgentNetworkProviderModel) *[]AgentNetworkProviderModelPricing {
	if len(models) == 0 {
		return nil
	}

	out := make([]AgentNetworkProviderModelPricing, len(models))

	for i, model := range models {
		out[i] = AgentNetworkProviderModelPricing{
			ID:                 model.Id,
			InputPer1k:         model.InputPer1k,
			OutputPer1k:        model.OutputPer1k,
			CacheReadPer1k:     model.CacheReadPer1k,
			CacheCreationPer1k: model.CacheCreationPer1k,
			CachedInputPer1k:   model.CachedInputPer1k,
		}
	}

	return &out
}

func equalAgentNetworkProviderModels(modelsA, modelsB *[]AgentNetworkProviderModelPricing) bool {
	var listA, listB []AgentNetworkProviderModelPricing

	if modelsA != nil {
		listA = *modelsA
	}

	if modelsB != nil {
		listB = *modelsB
	}

	if len(listA) != len(listB) {
		return false
	}

	sortedA := slices.Clone(listA)
	sortedB := slices.Clone(listB)

	sortByID := func(a, b AgentNetworkProviderModelPricing) int {
		if a.ID < b.ID {
			return -1
		}

		if a.ID > b.ID {
			return 1
		}

		return 0
	}

	slices.SortFunc(sortedA, sortByID)
	slices.SortFunc(sortedB, sortByID)

	for i := range sortedA {
		if !equalAgentNetworkProviderModel(sortedA[i], sortedB[i]) {
			return false
		}
	}

	return true
}

func equalAgentNetworkProviderModel(modelA, modelB AgentNetworkProviderModelPricing) bool {
	return modelA.ID == modelB.ID &&
		modelA.InputPer1k == modelB.InputPer1k &&
		modelA.OutputPer1k == modelB.OutputPer1k &&
		equalPtr(modelA.CacheReadPer1k, modelB.CacheReadPer1k) &&
		equalPtr(modelA.CacheCreationPer1k, modelB.CacheCreationPer1k) &&
		equalPtr(modelA.CachedInputPer1k, modelB.CachedInputPer1k)
}
