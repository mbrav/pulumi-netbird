package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

const (
	agentNetworkSettingsID = "agent-network-settings"

	// defaultAccessLogRetentionDays mirrors the API default for a bootstrap
	// request. It is applied in Check so that the update endpoint — which takes
	// a plain int, where 0 means "keep indefinitely" — never receives an
	// accidental 0 for an omitted input.
	defaultAccessLogRetentionDays = 30
)

// AgentNetworkSettings represents the per-account Agent Network gateway settings.
type AgentNetworkSettings struct{}

// Annotate adds a description to the AgentNetworkSettings resource type.
func (a *AgentNetworkSettings) Annotate(ann infer.Annotator) {
	ann.Describe(&a, "Per-account NetBird Agent Network gateway settings. This is a singleton "+
		"resource — only one instance exists per account. Creating it bootstraps the account's "+
		"gateway endpoint from exactly one of proxyAddress (the server allocates a label beneath "+
		"that cluster) or endpoint (the hostname is claimed verbatim as a dedicated endpoint). "+
		"The assigned endpoint is immutable; changing either field replaces the resource, which "+
		"releases the endpoint and allocates a new one.")
}

// AgentNetworkSettingsArgs defines input fields for Agent Network settings.
type AgentNetworkSettingsArgs struct {
	ProxyAddress           *string `pulumi:"proxyAddress,optional"`
	Endpoint               *string `pulumi:"endpoint,optional"`
	EnableLogCollection    bool    `pulumi:"enableLogCollection"`
	EnablePromptCollection bool    `pulumi:"enablePromptCollection"`
	RedactPii              bool    `pulumi:"redactPii"`
	AccessLogRetentionDays *int    `pulumi:"accessLogRetentionDays,optional"`
}

// Annotate provides documentation for AgentNetworkSettingsArgs fields.
func (a *AgentNetworkSettingsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProxyAddress, "Cluster address to allocate a labeled endpoint beneath: the server "+
		"assigns a label and the endpoint becomes `<label>.<proxyAddress>`. Mutually exclusive with "+
		"endpoint; exactly one of the two is required. Immutable — changing it replaces the resource.")
	ann.Describe(&a.Endpoint, "Hostname to claim as the account's self-addressed (dedicated) endpoint, "+
		"served only by a proxy declaring exactly that address. Mutually exclusive with proxyAddress; "+
		"exactly one of the two is required. Immutable — changing it replaces the resource.")
	ann.Describe(&a.EnableLogCollection, "Whether per-request access-log entries are collected for this account's agent-network traffic.")
	ann.Describe(&a.EnablePromptCollection, "Master switch for request/response prompt capture. Capture runs only when this is on AND a policy guardrail also enables it.")
	ann.Describe(&a.RedactPii, "Whether captured prompts have PII redacted.")
	ann.Describe(&a.AccessLogRetentionDays, "Days to retain full access-log rows; older rows are swept. "+
		"0 or less means keep indefinitely. Defaults to 30 when omitted.")
}

// AgentNetworkSettingsState represents the output state of Agent Network settings.
type AgentNetworkSettingsState struct {
	ProxyAddress           *string `pulumi:"proxyAddress,optional"`
	Endpoint               *string `pulumi:"endpoint,optional"`
	Dedicated              bool    `pulumi:"dedicated"`
	EnableLogCollection    bool    `pulumi:"enableLogCollection"`
	EnablePromptCollection bool    `pulumi:"enablePromptCollection"`
	RedactPii              bool    `pulumi:"redactPii"`
	AccessLogRetentionDays *int    `pulumi:"accessLogRetentionDays,optional"`
	CreatedAt              *string `pulumi:"createdAt,optional"`
	UpdatedAt              *string `pulumi:"updatedAt,optional"`
}

// Annotate provides documentation for AgentNetworkSettingsState fields.
func (a *AgentNetworkSettingsState) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProxyAddress, "Declared cluster address of the proxy serving this account's gateway. "+
		"Equal to endpoint when a dedicated proxy serves the account; otherwise the endpoint's immediate parent.")
	ann.Describe(&a.Endpoint, "Bare hostname agents call for this account. Empty until bootstrapped.")
	ann.Describe(&a.Dedicated, "Whether the account's gateway is served by a proxy dedicated to it (endpoint equals proxyAddress).")
	ann.Describe(&a.EnableLogCollection, "Whether per-request access-log entries are collected.")
	ann.Describe(&a.EnablePromptCollection, "Master switch for request/response prompt capture.")
	ann.Describe(&a.RedactPii, "Whether captured prompts have PII redacted.")
	ann.Describe(&a.AccessLogRetentionDays, "Days to retain full access-log rows.")
	ann.Describe(&a.CreatedAt, "Timestamp when the settings row was created. Absent until bootstrapped.")
	ann.Describe(&a.UpdatedAt, "Timestamp when the settings row was last updated. Absent until bootstrapped.")
}

// Create bootstraps the account's Agent Network settings row. Since this is a singleton, it
// uses a fixed ID.
func (*AgentNetworkSettings) Create(
	ctx context.Context, req infer.CreateRequest[AgentNetworkSettingsArgs],
) (infer.CreateResponse[AgentNetworkSettingsState], error) {
	p.GetLogger(ctx).Debugf(
		"Create:AgentNetworkSettings proxyAddress=%s endpoint=%s",
		strPtr(req.Inputs.ProxyAddress), strPtr(req.Inputs.Endpoint),
	)

	if req.DryRun {
		return infer.CreateResponse[AgentNetworkSettingsState]{
			ID:     agentNetworkSettingsID,
			Output: agentNetworkSettingsStateFromArgs(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	enableLogCollection := req.Inputs.EnableLogCollection
	enablePromptCollection := req.Inputs.EnablePromptCollection
	redactPii := req.Inputs.RedactPii

	created, err := client.AgentNetwork.CreateSettings(ctx, nbapi.AgentNetworkSettingsCreateRequest{
		ProxyAddress:           req.Inputs.ProxyAddress,
		Endpoint:               req.Inputs.Endpoint,
		EnableLogCollection:    &enableLogCollection,
		EnablePromptCollection: &enablePromptCollection,
		RedactPii:              &redactPii,
		AccessLogRetentionDays: req.Inputs.AccessLogRetentionDays,
	})
	if err != nil {
		if isConflictErr(err) {
			return infer.CreateResponse[AgentNetworkSettingsState]{}, fmt.Errorf(
				"the account's agent network settings are already bootstrapped; import them instead "+
					"(pulumi import netbird:resource:AgentNetworkSettings <name> %s): %w",
				agentNetworkSettingsID, err,
			)
		}

		return infer.CreateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("creating agent network settings failed: %w", err)
	}

	return infer.CreateResponse[AgentNetworkSettingsState]{
		ID:     agentNetworkSettingsID,
		Output: agentNetworkSettingsStateFromAPI(*created),
	}, nil
}

// Read reads the current Agent Network settings from NetBird.
func (*AgentNetworkSettings) Read(
	ctx context.Context, req infer.ReadRequest[AgentNetworkSettingsArgs, AgentNetworkSettingsState],
) (infer.ReadResponse[AgentNetworkSettingsArgs, AgentNetworkSettingsState], error) {
	p.GetLogger(ctx).Debugf("Read:AgentNetworkSettings[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[AgentNetworkSettingsArgs, AgentNetworkSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	settings, err := client.AgentNetwork.GetSettings(ctx)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[AgentNetworkSettingsArgs, AgentNetworkSettingsState]{
				ID:     "",
				Inputs: AgentNetworkSettingsArgs{},  //nolint:exhaustruct,exhaustruct_v5
				State:  AgentNetworkSettingsState{}, //nolint:exhaustruct,exhaustruct_v5
			}, nil
		}

		return infer.ReadResponse[AgentNetworkSettingsArgs, AgentNetworkSettingsState]{}, fmt.Errorf("reading agent network settings failed: %w", err)
	}

	state := agentNetworkSettingsStateFromAPI(*settings)

	// Only one of the two bootstrap fields can have been configured, and the
	// dedicated flag says which: a dedicated endpoint was claimed verbatim,
	// otherwise the label was allocated beneath the proxy address. Echoing back
	// exactly one keeps a refreshed or imported resource diff-free.
	inputs := AgentNetworkSettingsArgs{
		ProxyAddress:           state.ProxyAddress,
		Endpoint:               nil,
		EnableLogCollection:    settings.EnableLogCollection,
		EnablePromptCollection: settings.EnablePromptCollection,
		RedactPii:              settings.RedactPii,
		AccessLogRetentionDays: settings.AccessLogRetentionDays,
	}

	if settings.Dedicated {
		inputs.ProxyAddress = nil
		inputs.Endpoint = state.Endpoint
	}

	return infer.ReadResponse[AgentNetworkSettingsArgs, AgentNetworkSettingsState]{
		ID:     req.ID,
		Inputs: inputs,
		State:  state,
	}, nil
}

// Update updates the mutable Agent Network settings. The endpoint and proxy address are
// assigned at bootstrap and immutable, so the request echoes the assigned values back
// unchanged — the API rejects anything else.
func (*AgentNetworkSettings) Update(
	ctx context.Context, req infer.UpdateRequest[AgentNetworkSettingsArgs, AgentNetworkSettingsState],
) (infer.UpdateResponse[AgentNetworkSettingsState], error) {
	p.GetLogger(ctx).Debugf("Update:AgentNetworkSettings[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[AgentNetworkSettingsState]{
			Output: agentNetworkSettingsStateFromArgs(req.Inputs, req.State.CreatedAt, req.State.UpdatedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.AgentNetwork.UpdateSettings(ctx, nbapi.AgentNetworkSettingsRequest{
		Endpoint:               agentNetworkSettingsIdentity(req.State.Endpoint, req.Inputs.Endpoint),
		ProxyAddress:           agentNetworkSettingsIdentity(req.State.ProxyAddress, req.Inputs.ProxyAddress),
		EnableLogCollection:    req.Inputs.EnableLogCollection,
		EnablePromptCollection: req.Inputs.EnablePromptCollection,
		RedactPii:              req.Inputs.RedactPii,
		AccessLogRetentionDays: agentNetworkSettingsRetentionDays(req.Inputs.AccessLogRetentionDays),
	})
	if err != nil {
		return infer.UpdateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("updating agent network settings failed: %w", err)
	}

	return infer.UpdateResponse[AgentNetworkSettingsState]{
		Output: agentNetworkSettingsStateFromAPI(*updated),
	}, nil
}

// Delete releases the account's Agent Network endpoint by deleting the settings row. The API
// refuses while any provider still exists for the account or while a proxy is actively
// serving the endpoint; bootstrapping again allocates a new endpoint.
func (*AgentNetworkSettings) Delete(
	ctx context.Context, req infer.DeleteRequest[AgentNetworkSettingsState],
) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:AgentNetworkSettings[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.AgentNetwork.DeleteSettings(ctx)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting agent network settings failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between AgentNetworkSettingsArgs and AgentNetworkSettingsState.
func (*AgentNetworkSettings) Diff(
	ctx context.Context, req infer.DiffRequest[AgentNetworkSettingsArgs, AgentNetworkSettingsState],
) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:AgentNetworkSettings[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	// Exactly one of the two is configured and the other is server-assigned, so
	// each is compared only when it was actually set: a labeled bootstrap leaves
	// endpoint nil while the API reports the allocated hostname, and a dedicated
	// one leaves proxyAddress nil while the API reports the claimed address.
	// Comparing those unconditionally would propose the same replacement forever.
	if !equalServerAssignedPtr(req.Inputs.ProxyAddress, req.State.ProxyAddress) {
		diff["proxyAddress"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if !equalServerAssignedPtr(req.Inputs.Endpoint, req.State.Endpoint) {
		diff["endpoint"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.EnableLogCollection != req.State.EnableLogCollection {
		diff["enableLogCollection"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.EnablePromptCollection != req.State.EnablePromptCollection {
		diff["enablePromptCollection"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.RedactPii != req.State.RedactPii {
		diff["redactPii"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	// Compared through the default both sides resolve to: Check pins an omitted
	// input to the API default, and the API always reports a value back.
	if agentNetworkSettingsRetentionDays(req.Inputs.AccessLogRetentionDays) !=
		agentNetworkSettingsRetentionDays(req.State.AccessLogRetentionDays) {
		diff["accessLogRetentionDays"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:AgentNetworkSettings[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates the bootstrap inputs and applies the retention default.
func (*AgentNetworkSettings) Check(
	ctx context.Context, req infer.CheckRequest,
) (infer.CheckResponse[AgentNetworkSettingsArgs], error) {
	p.GetLogger(ctx).Debugf("Check:AgentNetworkSettings old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[AgentNetworkSettingsArgs](ctx, req.NewInputs)

	hasProxyAddress := args.ProxyAddress != nil && !isBlank(*args.ProxyAddress)
	hasEndpoint := args.Endpoint != nil && !isBlank(*args.Endpoint)

	switch {
	case hasProxyAddress && hasEndpoint:
		failures = append(failures, p.CheckFailure{
			Property: "proxyAddress",
			Reason:   "proxyAddress and endpoint are mutually exclusive — set exactly one",
		})
	case !hasProxyAddress && !hasEndpoint:
		failures = append(failures, p.CheckFailure{
			Property: "proxyAddress",
			Reason:   "exactly one of proxyAddress or endpoint is required to bootstrap the account gateway",
		})
	}

	if args.ProxyAddress != nil && isBlank(*args.ProxyAddress) {
		failures = append(failures, p.CheckFailure{Property: "proxyAddress", Reason: "proxyAddress must not be empty"})
	}

	if args.Endpoint != nil && isBlank(*args.Endpoint) {
		failures = append(failures, p.CheckFailure{Property: "endpoint", Reason: "endpoint must not be empty"})
	}

	// The update endpoint takes a plain int where 0 means "keep indefinitely",
	// so an omitted input is pinned to the API's documented default instead of
	// silently becoming 0 on the next update.
	if args.AccessLogRetentionDays == nil {
		retentionDays := defaultAccessLogRetentionDays
		args.AccessLogRetentionDays = &retentionDays
	}

	return infer.CheckResponse[AgentNetworkSettingsArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*AgentNetworkSettings) WireDependencies(
	field infer.FieldSelector, args *AgentNetworkSettingsArgs, state *AgentNetworkSettingsState,
) {
	field.OutputField(&state.ProxyAddress).DependsOn(field.InputField(&args.ProxyAddress))
	field.OutputField(&state.Endpoint).DependsOn(field.InputField(&args.Endpoint))
	field.OutputField(&state.EnableLogCollection).DependsOn(field.InputField(&args.EnableLogCollection))
	field.OutputField(&state.EnablePromptCollection).DependsOn(field.InputField(&args.EnablePromptCollection))
	field.OutputField(&state.RedactPii).DependsOn(field.InputField(&args.RedactPii))
	field.OutputField(&state.AccessLogRetentionDays).DependsOn(field.InputField(&args.AccessLogRetentionDays))
}

// agentNetworkSettingsRetentionDays resolves the access-log retention days,
// applying the API's documented default for an omitted value.
func agentNetworkSettingsRetentionDays(days *int) int {
	if days == nil {
		return defaultAccessLogRetentionDays
	}

	return *days
}

// agentNetworkSettingsIdentity resolves an immutable identity field for an update request,
// preferring the value the API assigned over the configured one.
func agentNetworkSettingsIdentity(state, input *string) string {
	if assigned := strPtr(state); assigned != "" {
		return assigned
	}

	return strPtr(input)
}

func agentNetworkSettingsStateFromArgs(
	args AgentNetworkSettingsArgs, createdAt, updatedAt *string,
) AgentNetworkSettingsState {
	return AgentNetworkSettingsState{
		ProxyAddress:           args.ProxyAddress,
		Endpoint:               args.Endpoint,
		Dedicated:              strPtr(args.Endpoint) != "",
		EnableLogCollection:    args.EnableLogCollection,
		EnablePromptCollection: args.EnablePromptCollection,
		RedactPii:              args.RedactPii,
		AccessLogRetentionDays: args.AccessLogRetentionDays,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}

func agentNetworkSettingsStateFromAPI(settings nbapi.AgentNetworkSettings) AgentNetworkSettingsState {
	var proxyAddress, endpoint *string

	if settings.ProxyAddress != "" {
		proxyAddress = &settings.ProxyAddress
	}

	if settings.Endpoint != "" {
		endpoint = &settings.Endpoint
	}

	var createdAt, updatedAt *string

	if settings.CreatedAt != nil {
		formatted := settings.CreatedAt.Format(idpTimeFormat)
		createdAt = &formatted
	}

	if settings.UpdatedAt != nil {
		formatted := settings.UpdatedAt.Format(idpTimeFormat)
		updatedAt = &formatted
	}

	return AgentNetworkSettingsState{
		ProxyAddress:           proxyAddress,
		Endpoint:               endpoint,
		Dedicated:              settings.Dedicated,
		EnableLogCollection:    settings.EnableLogCollection,
		EnablePromptCollection: settings.EnablePromptCollection,
		RedactPii:              settings.RedactPii,
		AccessLogRetentionDays: settings.AccessLogRetentionDays,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}
