package tests_test

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
)

// TestAgentNetworkProviderLifecycle covers create, read and delete for a
// provider configured with only its required fields.
func TestAgentNetworkProviderLifecycle(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("AgentNetworkProvider")
	inputs := agentNetworkProviderInputs()

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("Anthropic"), created.Properties.Get("name"))
	assert.Equal(t, property.New("anthropic_api"), created.Properties.Get("providerId"))

	readResp := read(t, server, urn, created.ID, created.Properties, inputs)
	assert.Equal(t, created.ID, readResp.ID)

	deleteResource(t, server, urn, created.ID, created.Properties)
}

// TestAgentNetworkProviderNoDiffWithoutIdentityHeaders is a regression test for
// the perpetual diff on identityHeaderUserId / identityHeaderGroups: the API
// always reports both headers on the wire, empty when unset, so comparing an
// omitted input against a pointer to "" proposed the same update on every
// preview and never converged.
func TestAgentNetworkProviderNoDiffWithoutIdentityHeaders(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("AgentNetworkProvider")
	inputs := agentNetworkProviderInputs()

	created := create(t, server, urn, inputs)
	assertNoDiff(t, server, urn, created.ID, created.Properties, inputs)

	// An update driven by an unrelated field must not reintroduce the diff.
	changed := withProps(inputs, "name", "Anthropic renamed")
	updated := update(t, server, urn, created.ID, created.Properties, changed, inputs)
	assertNoDiff(t, server, urn, created.ID, updated.Properties, changed)

	// Refreshing from the API must agree with the configuration too.
	readResp := read(t, server, urn, created.ID, updated.Properties, changed)
	assertNoDiff(t, server, urn, created.ID, readResp.Properties, readResp.Inputs)
}

// TestAgentNetworkProviderIdentityHeaderRoundTrip checks that setting a header
// diffs, applies and settles, and that removing it from the configuration
// clears it server-side rather than being silently kept.
func TestAgentNetworkProviderIdentityHeaderRoundTrip(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("AgentNetworkProvider")
	inputs := agentNetworkProviderInputs()

	created := create(t, server, urn, inputs)

	withHeader := withProps(inputs, "identityHeaderUserId", "x-nb-user")
	assert.True(t, diff(t, server, urn, created.ID, created.Properties, withHeader, inputs).HasChanges)

	set := update(t, server, urn, created.ID, created.Properties, withHeader, inputs)
	assert.Equal(t, property.New("x-nb-user"), set.Properties.Get("identityHeaderUserId"))
	assertNoDiff(t, server, urn, created.ID, set.Properties, withHeader)

	// Removing the field must be detected and must actually clear the header.
	assert.True(t, diff(t, server, urn, created.ID, set.Properties, inputs, withHeader).HasChanges)

	cleared := update(t, server, urn, created.ID, set.Properties, inputs, withHeader)
	assert.True(t, cleared.Properties.Get("identityHeaderUserId").IsNull() ||
		cleared.Properties.Get("identityHeaderUserId").AsString() == "")
	assertNoDiff(t, server, urn, created.ID, cleared.Properties, inputs)
}

// TestAgentNetworkSettingsLabeledBootstrap covers the current bootstrap model:
// proxyAddress allocates a labeled endpoint, both identity fields come back
// server-assigned, and neither shows up as a diff afterwards.
func TestAgentNetworkSettingsLabeledBootstrap(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("AgentNetworkSettings")
	inputs := agentNetworkSettingsInputs("proxyAddress", "eu.proxy.netbird.io")

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("brave-otter.eu.proxy.netbird.io"), created.Properties.Get("endpoint"))
	assert.Equal(t, property.New("eu.proxy.netbird.io"), created.Properties.Get("proxyAddress"))
	assert.Equal(t, property.New(false), created.Properties.Get("dedicated"))
	assertNoDiff(t, server, urn, created.ID, created.Properties, inputs)

	// The immutable endpoint must be echoed back on update, and toggling a
	// mutable field must not disturb it.
	changed := withProps(inputs, "enablePromptCollection", true)
	updated := update(t, server, urn, created.ID, created.Properties, changed, inputs)
	assert.Equal(t, property.New(true), updated.Properties.Get("enablePromptCollection"))
	assert.Equal(t, property.New("brave-otter.eu.proxy.netbird.io"), updated.Properties.Get("endpoint"))
	assertNoDiff(t, server, urn, created.ID, updated.Properties, changed)

	readResp := read(t, server, urn, created.ID, updated.Properties, changed)
	assertNoDiff(t, server, urn, created.ID, readResp.Properties, readResp.Inputs)

	deleteResource(t, server, urn, created.ID, updated.Properties)
}

// TestAgentNetworkSettingsDedicatedBootstrap covers claiming a hostname
// verbatim, where the proxy address is the one the server derives.
func TestAgentNetworkSettingsDedicatedBootstrap(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("AgentNetworkSettings")
	inputs := agentNetworkSettingsInputs("endpoint", "gw.example.com")

	created := create(t, server, urn, inputs)
	assert.Equal(t, property.New("gw.example.com"), created.Properties.Get("endpoint"))
	assert.Equal(t, property.New("gw.example.com"), created.Properties.Get("proxyAddress"))
	assert.Equal(t, property.New(true), created.Properties.Get("dedicated"))
	assertNoDiff(t, server, urn, created.ID, created.Properties, inputs)

	// The endpoint is assigned at bootstrap and immutable: changing it replaces.
	moved := withProps(inputs, "endpoint", "gw2.example.com")
	assert.True(t, diff(t, server, urn, created.ID, created.Properties, moved, inputs).HasChanges)
}

// TestAgentNetworkSettingsRequiresExactlyOneBootstrapField checks the
// mutually-exclusive Check validation.
func TestAgentNetworkSettingsRequiresExactlyOneBootstrapField(t *testing.T) {
	t.Parallel()

	server := newProviderServer(t, startMockServer(t))
	urn := testURN("AgentNetworkSettings")

	both := withProps(agentNetworkSettingsInputs("proxyAddress", "eu.proxy.netbird.io"), "endpoint", "gw.example.com")
	assertCheckFails(t, server, urn, both, "mutually exclusive")

	neither := props(
		"enableLogCollection", true,
		"enablePromptCollection", false,
		"redactPii", false,
	)
	assertCheckFails(t, server, urn, neither, "exactly one")
}

func agentNetworkProviderInputs() property.Map {
	return props(
		"name", "Anthropic",
		"providerId", "anthropic_api",
		"upstreamUrl", "https://api.anthropic.com",
		"apiKey", "sk-test",
		"enabled", true,
	)
}

func agentNetworkSettingsInputs(bootstrapField, bootstrapValue string) property.Map {
	return props(
		bootstrapField, bootstrapValue,
		"enableLogCollection", true,
		"enablePromptCollection", false,
		"redactPii", false,
	)
}
