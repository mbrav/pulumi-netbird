package main

import (
	"log"

	"github.com/mbrav/pulumi-netbird/provider"

	p "github.com/pulumi/pulumi-go-provider"
)

func main() {
	log.Printf("Starting provider %s v%s", provider.Name, provider.Version)
	err := p.RunProvider(provider.Name, provider.Version, provider.Provider())
	if err != nil {
		log.Fatalf("Provider failed: %v", err)
	}
}
