package tests_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/blang/semver"
	netbird "github.com/mbrav/pulumi-netbird/provider"
	"github.com/mbrav/pulumi-netbird/tests/mock"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	presource "github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/require"
)

// startMockServer starts an httptest.Server backed by the in-process mock and
// returns the server's base URL. The server is shut down when the test ends.
func startMockServer(t *testing.T) string {
	t.Helper()
	ts := httptest.NewServer(mock.NewServer())
	t.Cleanup(ts.Close)

	return ts.URL
}

// newProviderServer creates a pulumi-go-provider integration.Server configured
// to use the given mock server URL. Configure is called with a test bearer token
// so the mock accepts all requests.
func newProviderServer(t *testing.T, mockURL string) integration.Server {
	t.Helper()
	ctx := context.Background()
	server, err := integration.NewServer(
		ctx,
		netbird.Name,
		semver.MustParse(netbird.Version),
		integration.WithProvider(netbird.Provider()),
	)
	require.NoError(t, err)

	// infer's Configure reads config from Args (a property.Map), not Variables.
	// The keys match the pulumi struct tags on config.Config: "url" and "token".
	err = server.Configure(p.ConfigureRequest{
		Args: property.NewMap(map[string]property.Value{
			"url":   property.New(mockURL),
			"token": property.New("test-token"),
		}),
	})
	require.NoError(t, err)

	return server
}

// testURN builds a deterministic URN for the given resource type.
// Resources live in the "resource" module (derived from their Go package path).
func testURN(typ string) presource.URN {
	return presource.NewURN("stack", "proj", "", tokens.Type("netbird:resource:"+typ), "test")
}
