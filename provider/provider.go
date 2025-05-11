package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Change to const to disable semver Version management
var (
	Name    string = "netbird"
	Version string = "0.0.9"
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
		// WithNamespace("nb").
		WithConfig(infer.Config[*Config]()).
		WithResources(
			infer.Resource[*Group](),
			infer.Resource[*Network](),
			infer.Resource[*NetworkResource](),
			infer.Resource[*NetworkRouter](),
			infer.Resource[*Peer](),
			infer.Resource[*Policy](),
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
				"importBasePath":                 "github.com/mbrav/pulumi-netbird/sdk/go/netbird",
			},
			"nodejs": map[string]any{
				"packageName": "@mbrav/pulumi-netbird",
				"dependencies": map[string]string{
					"@pulumi/pulumi": "^3.0.0",
				},
			},
			"python": map[string]any{
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
