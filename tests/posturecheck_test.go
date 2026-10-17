package tests_test

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostureCheckCertificateCheckLifecycle covers the Certificate posture check
// netbird v0.81.0 added: the CA list must survive a create/read round trip.
func TestPostureCheckCertificateCheckLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("PostureCheck")
	inputs := certificatePostureCheckInputs(stringArray(caCertA, caCertB))

	created := create(t, server, urn, inputs)

	certs := certificateCheckCAs(t, created.Properties)
	assert.Len(t, certs, 2)
	assert.Contains(t, certs, caCertA)
	assert.Contains(t, certs, caCertB)

	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)
	assert.Len(t, certificateCheckCAs(t, readResp.Properties), 2)

	deleteResource(t, server, urn, created.ID, created.Properties)
}

// TestPostureCheckCertificateCheckNoDiffOnReorder pins that the CA list is
// compared order-insensitively, so a NetBird server returning the certificates
// in a different order than the program declared them does not plan an update
// on every preview.
func TestPostureCheckCertificateCheckNoDiffOnReorder(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("PostureCheck")
	inputs := certificatePostureCheckInputs(stringArray(caCertA, caCertB))

	created := create(t, server, urn, inputs)

	reordered := certificatePostureCheckInputs(stringArray(caCertB, caCertA))

	diffResp := diff(t, server, urn, created.ID, created.Properties, reordered, inputs)
	assert.False(t, diffResp.HasChanges, "reordering the CA list must not plan an update")

	deleteResource(t, server, urn, created.ID, created.Properties)
}

func TestPostureCheckCertificateCheckValidation(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("PostureCheck")

	assertCheckFails(t, server, urn,
		certificatePostureCheckInputs(stringArray()),
		"at least one CA certificate is required when certificateCheck is set")

	assertCheckFails(t, server, urn,
		certificatePostureCheckInputs(stringArray(caCertA, "  ")),
		"CA certificate must not be empty")
}

const (
	caCertA = "-----BEGIN CERTIFICATE-----\nAAAA\n-----END CERTIFICATE-----"
	caCertB = "-----BEGIN CERTIFICATE-----\nBBBB\n-----END CERTIFICATE-----"
)

func certificatePostureCheckInputs(cas property.Value) property.Map {
	return props(
		"name", "cert-posture",
		"checks", object(
			"certificateCheck", object("caCertificates", cas),
		),
	)
}

// certificateCheckCAs digs the CA list out of a posture check's property map.
func certificateCheckCAs(t *testing.T, properties property.Map) []string {
	t.Helper()

	checks := properties.Get("checks")
	require.True(t, checks.IsMap(), "checks must be a map")

	certCheck := checks.AsMap().Get("certificateCheck")
	require.True(t, certCheck.IsMap(), "certificateCheck must be a map")

	cas := certCheck.AsMap().Get("caCertificates")
	require.True(t, cas.IsArray(), "caCertificates must be an array")

	out := make([]string, 0, cas.AsArray().Len())
	for i := range cas.AsArray().Len() {
		out = append(out, cas.AsArray().Get(i).AsString())
	}

	return out
}
