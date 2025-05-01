package provider

import (
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// Define provider-level configuration.
type Config struct {
	Scream       *bool  `pulumi:"itsasecret,optional"`
	NetBirdToken string `pulumi:"netbirdToken"`
	NetBirdUrl   string `pulumi:"netbirdUrl"`
}

// Annotate provider configuration.
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.NetBirdUrl, "URL to Netbird API, example: https://nb.mydomain:33073")
	a.Describe(&c.NetBirdToken, "Netbird API Token")

	a.SetDefault(&c.NetBirdUrl, "https://nb.mydomain:33073", "NETBIRD_URL")
	a.SetDefault(&c.NetBirdToken, "", "NETBIRD_TOKEN")
}

// Configure validates the provider configuration.
func (c *Config) Configure() error {
	if c.NetBirdToken == "" {
		return fmt.Errorf("netbirdToken must be set in provider configuration")
	}
	if c.NetBirdUrl == "" {
		return fmt.Errorf("netbirdUrl must be set in provider configuration")
	}
	return nil
}
