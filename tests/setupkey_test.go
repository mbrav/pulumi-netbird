package tests_test

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
)

func TestSetupKeyLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("SetupKey")
	inputs := setupKeyInputs()

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("test-key"), created.Properties.Get("name"))
	assert.False(t, created.Properties.Get("key").IsNull())
	assert.Equal(t, property.New(true), created.Properties.Get("valid"))

	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)
	assert.False(t, readResp.Properties.Get("key").IsNull())

	deleteResource(t, server, urn, created.ID, created.Properties)
}

func setupKeyInputs() property.Map {
	return props(
		"name", "test-key",
		"type", "reusable",
		"expiresIn", float64(0),
		"usageLimit", float64(0),
		"autoGroups", stringArray(),
	)
}
