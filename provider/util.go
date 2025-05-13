package provider

import (
	"context"
	"slices"

	"github.com/netbirdio/netbird/management/client/rest"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Retrieve the NetBird client using the provider configuration.
func getNetBirdClient(ctx context.Context) (*rest.Client, error) {
	// Get the configuration from the provider's context
	config := infer.GetConfig[*Config](ctx)

	// BUG: Fix this workaround
	var nbToken string
	var nbURL string

	nbToken = config.NetBirdToken
	nbURL = config.NetBirdUrl

	// nbToken = ""
	// nbURL = ""

	// Create and return the client using the provided token and URL
	return rest.NewWithBearerToken(nbURL, nbToken), nil
}

// Helper to stringify a pointer safely
func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper to deref a string slice pointer safely
func strSlicePtr(s *[]string) []string {
	if s == nil {
		return []string{}
	}
	return *s
}

// Helper to compare string pointers safely
func equalPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// helper function to compare optional []string values
func equalSlicePtr(a, b *[]string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Copy and sort both before comparing
	aSorted := slices.Clone(*a)
	bSorted := slices.Clone(*b)
	slices.Sort(aSorted)
	slices.Sort(bSorted)

	if len(aSorted) != len(bSorted) {
		return false
	}
	for i := range aSorted {
		if aSorted[i] != bSorted[i] {
			return false
		}
	}
	return true
}
