package provider

import (
	"context"

	"github.com/netbirdio/netbird/management/client/rest"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Retrieve the NetBird client using the provider configuration.
func getNetBirdClient(ctx context.Context) (*rest.Client, error) {
	// Get the configuration from the provider's context

	nbToken := infer.GetConfig[Config](ctx).NetBirdToken
	nbURL := infer.GetConfig[Config](ctx).NetBirdUrl

	// Create and return the client using the provided token and URL
	return rest.NewWithBearerToken(nbURL, nbToken), nil
}

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
