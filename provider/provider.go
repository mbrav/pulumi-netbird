package provider

import (
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

var Version string

const Name string = "netbird"

// Define Provider.
func Provider() p.Provider {
	return infer.Provider(infer.Options{
		Metadata: schema.Metadata{
			DisplayName: Name,
			Description: "Manage NetBird resources declaratively using Pulumi's infrastructure-as-code framework.",
			Keywords: []string{
				"pulumi",
				"networking",
				"netbird",
				"security",
			},
			Homepage:          "https://pulumi.com",
			License:           "AGPL-3.0",
			Repository:        "https://github.com/mbrav/pulumi-netbird",
			PluginDownloadURL: "github://api.github.com/mbrav",
			Publisher:         "mbrav",
			LogoURL:           "https://raw.githubusercontent.com/mbrav/pulumi-netbird/master/assets/logo.webp",
			// This contains language specific details for generating the provider's SDKs
			LanguageMap: map[string]any{
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
			},
		},
		Resources: []infer.InferredResource{
			infer.Resource[Random](),
			infer.Resource[Network](),
		},
		Components: []infer.InferredComponent{
			infer.Component(NewRandomComponent),
		},
		Config: infer.Config[Config](),
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		},
	})
}

// Define provider-level configuration.
type Config struct {
	Scream       *bool  `pulumi:"itsasecret,optional"`
	NetBirdToken string `pulumi:"netbirdToken"`
	NetBirdURL   string `pulumi:"netbirdUrl"`
}

// Annotate provider configuration.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.NetBirdURL, "URL to Netbird API, example: https://nb.mydomain:33073")
	a.Describe(&c.NetBirdToken, "Netbird API Token")

	a.SetDefault(&c.NetBirdURL, "https://nb.mydomain:33073", "NETBIRD_URL")
	a.SetDefault(&c.NetBirdToken, "", "NETBIRD_TOKEN")
}

// Configure validates the provider configuration.
func (c *Config) Configure() error {
	if c.NetBirdToken == "" {
		return fmt.Errorf("netbirdToken must be set in provider configuration")
	}
	if c.NetBirdURL == "" {
		return fmt.Errorf("netbirdUrl must be set in provider configuration")
	}
	return nil
}
