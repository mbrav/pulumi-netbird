package tests_test

import (
	"testing"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func create(t *testing.T, server integration.Server, urn resource.URN, inputs property.Map) p.CreateResponse {
	t.Helper()

	check, err := server.Check(p.CheckRequest{Urn: urn, Inputs: inputs})
	require.NoError(t, err)
	require.Empty(t, check.Failures)

	created, err := server.Create(p.CreateRequest{Urn: urn, Properties: inputs})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	return created
}

func read(
	t *testing.T,
	server integration.Server,
	urn resource.URN,
	id string,
	state property.Map,
	inputs property.Map,
) p.ReadResponse {
	t.Helper()

	read, err := server.Read(p.ReadRequest{
		ID:         id,
		Urn:        urn,
		Properties: state,
		Inputs:     inputs,
	})
	require.NoError(t, err)

	return read
}

func update(
	t *testing.T,
	server integration.Server,
	urn resource.URN,
	id string,
	state property.Map,
	inputs property.Map,
	oldInputs property.Map,
) p.UpdateResponse {
	t.Helper()

	updated, err := server.Update(p.UpdateRequest{
		ID:        id,
		Urn:       urn,
		State:     state,
		Inputs:    inputs,
		OldInputs: oldInputs,
	})
	require.NoError(t, err)

	return updated
}

func deleteResource(t *testing.T, server integration.Server, urn resource.URN, id string, state property.Map) {
	t.Helper()

	require.NoError(t, server.Delete(p.DeleteRequest{
		ID:         id,
		Urn:        urn,
		Properties: state,
	}))
}

func diff(
	t *testing.T,
	server integration.Server,
	urn resource.URN,
	id string,
	state property.Map,
	inputs property.Map,
	oldInputs property.Map,
) p.DiffResponse {
	t.Helper()

	diff, err := server.Diff(p.DiffRequest{
		ID:        id,
		Urn:       urn,
		State:     state,
		Inputs:    inputs,
		OldInputs: oldInputs,
	})
	require.NoError(t, err)

	return diff
}

func assertNoDiff(t *testing.T, server integration.Server, urn resource.URN, id string, state property.Map, inputs property.Map) {
	t.Helper()

	assert.False(t, diff(t, server, urn, id, state, inputs, inputs).HasChanges)
}

func props(kv ...any) property.Map {
	values := make(map[string]property.Value, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		values[kv[i].(string)] = prop(kv[i+1])
	}

	return property.NewMap(values)
}

func prop(v any) property.Value {
	switch v := v.(type) {
	case property.Value:
		return v
	case string:
		return property.New(v)
	case bool:
		return property.New(v)
	case float64:
		return property.New(v)
	case map[string][]string:
		values := make(map[string]property.Value, len(v))
		for key, items := range v {
			values[key] = stringArray(items...)
		}

		return property.New(property.NewMap(values))
	default:
		panic("unsupported property value")
	}
}

func object(kv ...any) property.Value {
	return property.New(props(kv...))
}

func array(values ...property.Value) property.Value {
	return property.New(property.NewArray(values))
}

func stringArray(values ...string) property.Value {
	out := make([]property.Value, len(values))
	for i, value := range values {
		out[i] = property.New(value)
	}

	return array(out...)
}
