package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

const Name string = "netbird"

func Provider() p.Provider {
	// We tell the provider what resources it needs to support.
	// In this case, a single resource and component
	return infer.Provider(infer.Options{
		Resources: []infer.InferredResource{
			infer.Resource[Random](),
		},
		Components: []infer.InferredComponent{
			infer.Component(NewRandomComponent),
		},
		Config: infer.Config[Config](),
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		},
	})
}

// Define some provider-level configuration
type Config struct {
	Scream *bool `pulumi:"itsasecret,optional"`
}
