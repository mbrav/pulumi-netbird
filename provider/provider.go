package provider

import (
	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/mbrav/pulumi-netbird/provider/resource"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Change to const to disable semver Version management
var (
	Name    string = "netbird"
	Version string = "0.0.12"
)

// Define Provider.
func Provider() p.Provider {
	// Build provider
	provider, err := infer.NewProviderBuilder().
		WithDisplayName(Name).
		WithDescription("Manage NetBird resources declaratively using Pulumi's infrastructure-as-code framework.").
		WithKeywords(
			"pulumi",
			"networking",
			"netbird",
			"security",
		).
		WithHomepage("https://pulumi.com").
		WithLicense("AGPL-3.0").
		WithRepository("https://github.com/mbrav/pulumi-netbird").
		WithPluginDownloadURL("github://api.github.com/mbrav").
		WithPublisher("mbrav").
		WithLogoURL("https://raw.githubusercontent.com/mbrav/pulumi-netbird/master/assets/logo.webp").
		WithNamespace("pulumi").
		// WithWrapped(provider p.Provider),
		WithConfig(infer.Config(&config.Config{})).
		WithResources(
			infer.Resource(&resource.DNS{}),
			infer.Resource(&resource.Group{}),
			infer.Resource(&resource.Network{}),
			infer.Resource(&resource.NetworkResource{}),
			infer.Resource(&resource.NetworkRouter{}),
			infer.Resource(&resource.Peer{}),
			infer.Resource(&resource.Policy{}),
		).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"auto-naming": "index",
		}).
		WithLanguageMap(map[string]any{
			"csharp": map[string]any{
				"packageReferences": map[string]string{
					"Pulumi": "3.*",
				},
			},
			"go": map[string]any{
				"generateResourceContainerTypes": true,
				"respectSchemaVersion":           true,
				"importBasePath":                 "github.com/mbrav/pulumi-netbird/sdk/go/netbird",
			},
			"nodejs": map[string]any{
				"packageName": "@mbrav/pulumi-netbird",
				"dependencies": map[string]string{
					"@pulumi/pulumi": "^3.0.0",
				},
			},
			"python": map[string]any{
				"respectSchemaVersion": true,
				"pyproject": map[string]bool{
					"enabled": true,
				},
				"requires": map[string]string{
					"pulumi": ">=3.0.0,<4.0.0",
				},
			},
			"java": map[string]any{
				"buildFiles":                      "gradle",
				"gradleNexusPublishPluginVersion": "2.0.0",
				"dependencies": map[string]any{
					"com.pulumi:pulumi": "1.10.0",
				},
			},
		}).
		Build()
		// Check error
	if err != nil {
		panic("failed to build provider: " + err.Error())
	}

	return provider
}
