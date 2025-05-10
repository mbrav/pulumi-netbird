package provider

import (
	"context"
	"fmt"
	"os"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

var _ = (infer.CustomConfigure)((*Config)(nil))

// Define provider-level configuration.
type Config struct {
	NetBirdUrl   string `pulumi:"netbirdUrl"`
	NetBirdToken string `pulumi:"netbirdToken"`
}

// Annotate provider configuration.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.NetBirdUrl, "URL to Netbird API, example: https://nb.mydomain:33073")
	a.Describe(&c.NetBirdToken, "Netbird API Token")

	a.SetDefault(&c.NetBirdUrl, "https://nb.mydomain:33073", "NETBIRD_URL")
	a.SetDefault(&c.NetBirdToken, "", "NETBIRD_TOKEN")
}

// Configure validates the provider configuration.
func (c *Config) Configure(ctx context.Context) error {
	if envVal, exists := os.LookupEnv("NETBIRD_URL"); exists {
		c.NetBirdUrl = envVal
	}
	if envVal, exists := os.LookupEnv("NETBIRD_TOKEN"); exists {
		c.NetBirdToken = envVal
	}
	p.GetLogger(ctx).Debugf("Config netbirdToken=%s, netbirdUrl=%s", c.NetBirdUrl, c.NetBirdToken)
	if c.NetBirdToken == "" {
		return fmt.Errorf("netbirdToken must be set in provider configuration")
	}
	if c.NetBirdUrl == "" {
		return fmt.Errorf("netbirdUrl must be set in provider configuration")
	}
	return nil
}
