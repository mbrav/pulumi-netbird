package main

import (
	"github.com/mbrav/pulumi-netbird/provider"

	p "github.com/pulumi/pulumi-go-provider"
)

func main() {
	p.RunProvider("netbird", provider.Version, provider.Provider())
}
