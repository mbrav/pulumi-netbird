package tests_test

import (
	"testing"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
)

func TestRouteLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("Route")
	inputs := routeInputs("net-1")

	created := create(t, server, urn, inputs)
	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)

	assertNoDiff(t, server, urn, created.ID, readResp.Properties, inputs)
	deleteResource(t, server, urn, created.ID, created.Properties)
	assert.Empty(t, read(t, server, urn, created.ID, created.Properties, inputs).ID)
}

func TestRouteNetworkIDChangeRequiresReplace(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("Route")
	inputs := routeInputs("net-1")

	created := create(t, server, urn, inputs)
	routeDiff := diff(t, server, urn, created.ID, created.Properties, routeInputs("net-2"), inputs)

	assert.True(t, routeDiff.HasChanges)
	assert.Equal(t, p.UpdateReplace, routeDiff.DetailedDiff["networkId"].Kind)
}

func routeInputs(networkID string) property.Map {
	return props(
		"networkId", networkID,
		"description", "test route",
		"enabled", true,
		"masquerade", false,
		"metric", float64(9999),
		"keepRoute", false,
		"network", "10.0.0.0/8",
		"groups", stringArray(),
		"peerGroups", stringArray("all"),
	)
}
