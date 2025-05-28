package config

import (
	"context"

	"github.com/netbirdio/netbird/management/client/rest"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Define provider-level configuration.
type Config struct {
	NetBirdUrl   string `pulumi:"netbirdUrl"`
	NetBirdToken string `provider:"secret"   pulumi:"netbirdToken"`
}

// Annotate provider configuration.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.NetBirdUrl, "URL to Netbird API, example: https://api.netbird.io")
	a.Describe(&c.NetBirdToken, "Netbird API Token")

	a.SetDefault(&c.NetBirdUrl, "https://api.netbird.io", "NETBIRD_URL")
	a.SetDefault(&c.NetBirdToken, "", "NETBIRD_TOKEN")
}

// Configure validates the provider configuration.
func (c *Config) Configure(ctx context.Context) error {
	p.GetLogger(ctx).Debugf("Configure:Config")
	// p.GetLogger(ctx).Debugf("Config netbirdToken=%s, netbirdUrl=%s", c.NetBirdUrl, c.NetBirdToken)

	if c.NetBirdToken == "" {
		return ErrMissingNetBirdToken
	}

	if c.NetBirdUrl == "" {
		return ErrMissingNetBirdURL
	}

	return nil
}

// Retrieve the NetBird client using the provider configuration.
func GetNetBirdClient(ctx context.Context) (*rest.Client, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return nil, ErrNilProviderConfig
	}

	if config.NetBirdToken == "" {
		return nil, ErrMissingNetBirdToken
	}

	if config.NetBirdUrl == "" {
		return nil, ErrMissingNetBirdURL
	}

	client := rest.NewWithBearerToken(config.NetBirdUrl, config.NetBirdToken)

	return client, nil
}

// Retrieve the NetBird URL using the provider configuration.
func GetNetBirdURL(ctx context.Context) (string, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", ErrNilProviderConfig
	}

	if config.NetBirdUrl == "" {
		return "", ErrMissingNetBirdURL
	}

	return config.NetBirdUrl, nil
}
