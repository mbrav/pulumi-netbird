// Package resource provides the NetBird resource types
package resource

import "github.com/pulumi/pulumi-go-provider/infer"

// All returns all registered provider resources.
func All() []infer.InferredResource {
	return []infer.InferredResource{
		infer.Resource(&DNS{}),
		infer.Resource(&DNSRecord{}),
		infer.Resource(&DNSSettings{}),
		infer.Resource(&DNSZone{}),
		infer.Resource(&Group{}),
		infer.Resource(&Network{}),
		infer.Resource(&NetworkResource{}),
		infer.Resource(&NetworkRouter{}),
		infer.Resource(&Peer{}),
		infer.Resource(&Policy{}),
		infer.Resource(&PostureCheck{}),
		infer.Resource(&ReverseProxyDomain{}),
		infer.Resource(&ReverseProxyService{}),
		infer.Resource(&Route{}),
		infer.Resource(&SetupKey{}),
		infer.Resource(&User{}),
	}
}
