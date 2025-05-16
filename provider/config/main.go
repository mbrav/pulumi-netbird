package config

import (
	"context"
	"errors"

	"github.com/netbirdio/netbird/management/client/rest"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Define provider-level configuration.
type Config struct {
	NetBirdUrl   string `pulumi:"netbirdUrl"`
	NetBirdToken string `pulumi:"netbirdToken" provider:"secret"`
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
	p.GetLogger(ctx).Debugf("Config netbirdToken=%s, netbirdUrl=%s", c.NetBirdUrl, c.NetBirdToken)

	if c.NetBirdToken == "" {
		return errors.New("netbirdToken must be set in provider configuration")
	}

	if c.NetBirdUrl == "" {
		return errors.New("netbirdUrl must be set in provider configuration")
	}

	return nil
}

// Retrieve the NetBird client using the provider configuration.
func GetNetBirdClient(ctx context.Context) (*rest.Client, error) {
	// Get the configuration from the provider's context
	config := infer.GetConfig[*Config](ctx)

	nbToken := config.NetBirdToken
	nbURL := config.NetBirdUrl

	// Create and return the client using the provided token and URL
	return rest.NewWithBearerToken(nbURL, nbToken), nil
}
