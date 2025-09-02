package config

import (
	"context"

	"github.com/netbirdio/netbird/shared/management/client/rest"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Config holds the provider configuration for NetBiryyd.
type Config struct {
	NetBirdURL   string `pulumi:"url"`
	NetBirdToken string `provider:"secret" pulumi:"token"`
}

// Annotate provider configuration.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.NetBirdURL, "URL to Netbird API, example: https://api.netbird.io")
	a.Describe(&c.NetBirdToken, "Netbird API Token")

	a.SetDefault(&c.NetBirdURL, "https://api.netbird.io", "NETBIRD_URL")
	a.SetDefault(&c.NetBirdToken, "", "NETBIRD_TOKEN")
}

// Configure validates the provider configuration.
func (c *Config) Configure(ctx context.Context) error {
	p.GetLogger(ctx).Debugf("Configure:Config")
	// p.GetLogger(ctx).Debugf("Config netbirdToken=%s, netbirdUrl=%s", c.NetBirdUrl, c.NetBirdToken)

	if c.NetBirdToken == "" {
		return ErrMissingNetBirdToken
	}

	if c.NetBirdURL == "" {
		return ErrMissingNetBirdURL
	}

	return nil
}

// GetNetBirdClient creates and returns a new NetBird REST client using the provider configuration from the given context.
func GetNetBirdClient(ctx context.Context) (*rest.Client, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return nil, ErrNilProviderConfig
	}

	if config.NetBirdToken == "" {
		return nil, ErrMissingNetBirdToken
	}

	if config.NetBirdURL == "" {
		return nil, ErrMissingNetBirdURL
	}

	client := rest.NewWithBearerToken(config.NetBirdURL, config.NetBirdToken)

	return client, nil
}

// GetNetBirdURL retrieves the NetBird URL from the provider configuration in the given context.
func GetNetBirdURL(ctx context.Context) (string, error) {
	config := infer.GetConfig[*Config](ctx)
	if config == nil {
		return "", ErrNilProviderConfig
	}

	if config.NetBirdURL == "" {
		return "", ErrMissingNetBirdURL
	}

	return config.NetBirdURL, nil
}
