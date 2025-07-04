// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package config

import (
	"github.com/mbrav/pulumi-netbird/sdk/go/netbird/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

var _ = internal.GetEnvOrDefault

// Netbird API Token
func GetToken(ctx *pulumi.Context) string {
	v, err := config.Try(ctx, "netbird:token")
	if err == nil {
		return v
	}
	var value string
	if d := internal.GetEnvOrDefault("", nil, "NETBIRD_TOKEN"); d != nil {
		value = d.(string)
	}
	return value
}

// URL to Netbird API, example: https://api.netbird.io
func GetUrl(ctx *pulumi.Context) string {
	v, err := config.Try(ctx, "netbird:url")
	if err == nil {
		return v
	}
	var value string
	if d := internal.GetEnvOrDefault("https://api.netbird.io", nil, "NETBIRD_URL"); d != nil {
		value = d.(string)
	}
	return value
}
