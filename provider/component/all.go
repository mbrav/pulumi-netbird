package component

import "github.com/pulumi/pulumi-go-provider/infer"

// All returns all registered provider components.
func All() []infer.InferredComponent {
	return []infer.InferredComponent{
		infer.Component(&NetworkBundle{}),
		infer.Component(&DNSZoneBundle{}),
	}
}
