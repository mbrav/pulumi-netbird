package tests_test

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
)

func TestUserLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("User")
	inputs := userInputs(false)

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("admin"), created.Properties.Get("role"))
	assert.Equal(t, property.New("active"), created.Properties.Get("status"))

	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)
	assert.Equal(t, property.New("active"), readResp.Properties.Get("status"))

	deleteResource(t, server, urn, created.ID, created.Properties)
}

// TestUserUpdateRefreshesStatus is a regression test for a v0.5.3 bug: User.Update
// discarded the API's update response and echoed back the pre-update `status`,
// so blocking a user never showed status: blocked until the next Read/refresh.
func TestUserUpdateRefreshesStatus(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("User")
	inputs := userInputs(false)

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("active"), created.Properties.Get("status"))

	blockedInputs := userInputs(true)
	updated := update(t, server, urn, created.ID, created.Properties, blockedInputs, inputs)

	assert.Equal(t, property.New(true), updated.Properties.Get("blocked"))
	assert.Equal(t, property.New("blocked"), updated.Properties.Get("status"))
}

func userInputs(blocked bool) property.Map {
	return props(
		"role", "admin",
		"isServiceUser", true,
		"name", "ci-bot",
		"autoGroups", stringArray(),
		"blocked", blocked,
	)
}
