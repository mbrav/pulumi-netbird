package tests_test

import (
	"testing"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReverseProxyDomainLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("ReverseProxyDomain")
	inputs := reverseProxyDomainInputs("corp.example.com")

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("corp.example.com"), created.Properties.Get("domain"))

	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)

	deleteResource(t, server, urn, created.ID, created.Properties)
}

// TestReverseProxyDomainDeleteIsIdempotent pins the 404-as-success contract every
// Delete in this provider follows: a domain already gone from NetBird must not
// fail the destroy. ReverseProxyDomain.Delete returned the raw error before
// 0.5.12, so a stack that lost the domain out-of-band could never be destroyed.
func TestReverseProxyDomainDeleteIsIdempotent(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("ReverseProxyDomain")

	created := create(t, server, urn, reverseProxyDomainInputs("corp.example.com"))

	deleteResource(t, server, urn, created.ID, created.Properties)
	// Second delete hits a 404 and must still succeed.
	deleteResource(t, server, urn, created.ID, created.Properties)
}

// TestReverseProxyDomainDeleteInUseFails pins the opposite case: netbird v0.80.0
// answers 412 while a service still uses the domain, and that must surface as an
// error rather than being swallowed like a 404 — otherwise Pulumi would drop the
// domain from state while it still exists.
func TestReverseProxyDomainDeleteInUseFails(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("ReverseProxyDomain")

	created := create(t, server, urn, reverseProxyDomainInputs("inuse.example.com"))

	err := server.Delete(p.DeleteRequest{
		ID:         created.ID,
		Urn:        urn,
		Properties: created.Properties,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "still used by a service")
}

func reverseProxyDomainInputs(domain string) property.Map {
	return props(
		"domain", domain,
		"targetCluster", "cluster-1",
	)
}
