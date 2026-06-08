package tests_test

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
)

func TestPolicyLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("Policy")
	inputs := policyInputs()

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("test-policy"), created.Properties.Get("name"))

	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)

	assertNoDiff(t, server, urn, created.ID, readResp.Properties, inputs)

	importRead := read(t, server, urn, created.ID, created.Properties, property.Map{})
	assert.Equal(t, created.ID, importRead.ID)
	assert.False(t, importRead.Inputs.Get("rules").IsNull())

	deleteResource(t, server, urn, created.ID, created.Properties)
	assert.Empty(t, read(t, server, urn, created.ID, property.Map{}, inputs).ID)
}

func policyInputs() property.Map {
	return props(
		"name", "test-policy",
		"enabled", true,
		"postureChecks", stringArray("check-b", "check-a"),
		"rules", array(object(
			"name", "allow-all",
			"enabled", true,
			"bidirectional", true,
			"action", "accept",
			"protocol", "all",
			"sources", stringArray("all"),
			"destinations", stringArray("all"),
			"authorizedGroups", map[string][]string{
				"all": {"root"},
			},
		)),
	)
}
