// Package function provides the NetBird Pulumi provider functions (invokes / data sources).
package function

import "github.com/pulumi/pulumi-go-provider/infer"

// All returns all registered provider functions.
func All() []infer.InferredFunction {
	return []infer.InferredFunction{
		infer.Function(&GetPeers{}),
		infer.Function(&LookupGroup{}),
		infer.Function(&LookupPeer{}),
		infer.Function(&LookupRoute{}),
		infer.Function(&LookupSetupKey{}),
		infer.Function(&LookupUser{}),
	}
}
