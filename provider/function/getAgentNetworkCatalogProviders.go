package function

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// GetAgentNetworkCatalogProviders lists the catalog of supported upstream AI providers
// (openai_api, anthropic_api, bedrock_api, ...) with their default models and pricing, used to
// prefill AgentNetworkProvider create forms.
type GetAgentNetworkCatalogProviders struct{}

// Annotate describes the function.
func (f *GetAgentNetworkCatalogProviders) Annotate(a infer.Annotator) {
	a.Describe(f, "List the catalog of upstream AI providers supported by NetBird Agent "+
		"Network, with their default models and pricing. Useful for looking up a providerId "+
		"and model ids for AgentNetworkProvider.")
}

// GetAgentNetworkCatalogProvidersArgs are the inputs for GetAgentNetworkCatalogProviders (none).
type GetAgentNetworkCatalogProvidersArgs struct{}

// AgentNetworkCatalogModel is one model offered by a catalog provider, with default pricing.
type AgentNetworkCatalogModel struct {
	ID                 string   `pulumi:"id"`
	Label              string   `pulumi:"label"`
	ContextWindow      int      `pulumi:"contextWindow"`
	InputPer1k         float64  `pulumi:"inputPer1k"`
	OutputPer1k        float64  `pulumi:"outputPer1k"`
	CacheReadPer1k     *float64 `pulumi:"cacheReadPer1k,optional"`
	CacheCreationPer1k *float64 `pulumi:"cacheCreationPer1k,optional"`
	CachedInputPer1k   *float64 `pulumi:"cachedInputPer1k,optional"`
}

// Annotate provides field descriptions for AgentNetworkCatalogModel.
func (m *AgentNetworkCatalogModel) Annotate(annotator infer.Annotator) {
	annotator.Describe(&m.ID, "Catalog model identifier as exposed by the upstream provider.")
	annotator.Describe(&m.Label, "Human-friendly model name.")
	annotator.Describe(&m.ContextWindow, "Maximum context window in tokens.")
	annotator.Describe(&m.InputPer1k, "Default input token price per 1k tokens, in USD.")
	annotator.Describe(&m.OutputPer1k, "Default output token price per 1k tokens, in USD.")
	annotator.Describe(&m.CacheReadPer1k, "Anthropic-shape cache rate: default cost per 1k cache-read tokens, in USD.")
	annotator.Describe(&m.CacheCreationPer1k, "Anthropic-shape cache rate: default cost per 1k cache-creation tokens, in USD.")
	annotator.Describe(&m.CachedInputPer1k, "OpenAI-shape cache rate: default cost per 1k cached prompt tokens, in USD.")
}

// AgentNetworkCatalogProviderEntry is one entry in the Agent Network provider catalog.
type AgentNetworkCatalogProviderEntry struct {
	ID                 string                     `pulumi:"id"`
	Name               string                     `pulumi:"name"`
	Kind               string                     `pulumi:"kind"`
	Description        string                     `pulumi:"description"`
	DefaultHost        string                     `pulumi:"defaultHost"`
	DefaultContentType string                     `pulumi:"defaultContentType"`
	AuthHeaderTemplate string                     `pulumi:"authHeaderTemplate"`
	BrandColor         string                     `pulumi:"brandColor"`
	PricingSurfaces    *[]string                  `pulumi:"pricingSurfaces,optional"`
	Models             []AgentNetworkCatalogModel `pulumi:"models"`
}

// Annotate provides field descriptions for AgentNetworkCatalogProviderEntry.
func (e *AgentNetworkCatalogProviderEntry) Annotate(annotator infer.Annotator) {
	annotator.Describe(&e.ID, "Catalog provider identifier, used as AgentNetworkProvider.providerId.")
	annotator.Describe(&e.Name, "Display name for the provider.")
	annotator.Describe(&e.Kind, "Presentation grouping: \"provider\" (first-party vendor API), \"gateway\" (routing/aggregation layer), or \"custom\" (generic OpenAI-compatible endpoint).")
	annotator.Describe(&e.Description, "Short description shown in the provider picker.")
	annotator.Describe(&e.DefaultHost, "Default upstream host suggested when adding a provider of this type.")
	annotator.Describe(&e.DefaultContentType, "Default Content-Type for upstream requests.")
	annotator.Describe(&e.AuthHeaderTemplate, "Template the proxy uses to inject the API key.")
	annotator.Describe(&e.BrandColor, "Hex brand color used to render the provider badge in the dashboard.")
	annotator.Describe(&e.PricingSurfaces, "Cost-meter pricing surfaces this provider's traffic is metered under (\"openai\", \"anthropic\", \"bedrock\").")
	annotator.Describe(&e.Models, "Catalog models available for this provider, with default pricing.")
}

// GetAgentNetworkCatalogProvidersResult is the output of GetAgentNetworkCatalogProviders.
type GetAgentNetworkCatalogProvidersResult struct {
	Providers []AgentNetworkCatalogProviderEntry `pulumi:"providers"`
}

// Annotate provides field descriptions for GetAgentNetworkCatalogProvidersResult.
func (r *GetAgentNetworkCatalogProvidersResult) Annotate(a infer.Annotator) {
	a.Describe(&r.Providers, "The list of catalog providers.")
}

// Invoke lists the Agent Network provider catalog.
func (f *GetAgentNetworkCatalogProviders) Invoke(
	ctx context.Context,
	_ infer.FunctionRequest[GetAgentNetworkCatalogProvidersArgs],
) (infer.FunctionResponse[GetAgentNetworkCatalogProvidersResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[GetAgentNetworkCatalogProvidersResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiProviders, err := client.AgentNetwork.ListCatalogProviders(ctx)
	if err != nil {
		return infer.FunctionResponse[GetAgentNetworkCatalogProvidersResult]{}, fmt.Errorf("listing agent network catalog providers failed: %w", err)
	}

	providers := make([]AgentNetworkCatalogProviderEntry, 0, len(apiProviders))

	for _, provider := range apiProviders {
		models := make([]AgentNetworkCatalogModel, 0, len(provider.Models))

		for _, model := range provider.Models {
			models = append(models, AgentNetworkCatalogModel{
				ID:                 model.Id,
				Label:              model.Label,
				ContextWindow:      model.ContextWindow,
				InputPer1k:         model.InputPer1k,
				OutputPer1k:        model.OutputPer1k,
				CacheReadPer1k:     model.CacheReadPer1k,
				CacheCreationPer1k: model.CacheCreationPer1k,
				CachedInputPer1k:   model.CachedInputPer1k,
			})
		}

		providers = append(providers, AgentNetworkCatalogProviderEntry{
			ID:                 provider.Id,
			Name:               provider.Name,
			Kind:               string(provider.Kind),
			Description:        provider.Description,
			DefaultHost:        provider.DefaultHost,
			DefaultContentType: provider.DefaultContentType,
			AuthHeaderTemplate: provider.AuthHeaderTemplate,
			BrandColor:         provider.BrandColor,
			PricingSurfaces:    provider.PricingSurfaces,
			Models:             models,
		})
	}

	return infer.FunctionResponse[GetAgentNetworkCatalogProvidersResult]{
		Output: GetAgentNetworkCatalogProvidersResult{
			Providers: providers,
		},
	}, nil
}
