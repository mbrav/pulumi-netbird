// Package component provides composite components for the NetBird Pulumi provider.
package component

import (
	"github.com/mbrav/pulumi-netbird/provider/resource"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Resource type tokens derived from the actual resource structs at init time.
// Using infer.Resource(...).GetToken() mirrors exactly how provider.go registers
// each resource, so the strings can never drift out of sync.
//nolint:gochecknoglobals // init-time constants; values are immutable after startup
var (
	tokenNetwork         = mustToken(infer.Resource(&resource.Network{}))
	tokenNetworkRouter   = mustToken(infer.Resource(&resource.NetworkRouter{}))
	tokenNetworkResource = mustToken(infer.Resource(&resource.NetworkResource{}))
	tokenDNSZone         = mustToken(infer.Resource(&resource.DNSZone{}))
	tokenDNSRecord       = mustToken(infer.Resource(&resource.DNSRecord{}))
)

// mustToken panics if the token cannot be derived — a programming error, not a
// runtime condition, so panic is appropriate here.
func mustToken(r infer.InferredResource) string {
	tok, err := r.GetToken()
	if err != nil {
		panic("component: failed to derive resource token: " + err.Error())
	}

	return tok.String()
}
