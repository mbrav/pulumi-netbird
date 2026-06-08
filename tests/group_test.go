package tests_test

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
)

func TestGroupLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("Group")
	inputs := groupInputs("test-group")

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("test-group"), created.Properties.Get("name"))

	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)

	updatedInputs := groupInputs("updated-group")
	updated := update(t, server, urn, created.ID, created.Properties, updatedInputs, inputs)
	assert.Equal(t, property.New("updated-group"), updated.Properties.Get("name"))

	assertNoDiff(t, server, urn, created.ID, updated.Properties, updatedInputs)
	assert.True(t, diff(t, server, urn, created.ID, updated.Properties, groupInputs("new-name"), updatedInputs).HasChanges)

	deleteResource(t, server, urn, created.ID, updated.Properties)
	assert.Empty(t, read(t, server, urn, created.ID, property.Map{}, inputs).ID)
}

func groupInputs(name string) property.Map {
	return props("name", name)
}
