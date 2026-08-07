package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

const agentNetworkSettingsID = "agent-network-settings"

// AgentNetworkSettings represents the per-account Agent Network gateway settings.
// This is a singleton resource — only one instance exists per account. Setting Cluster on
// create bootstraps the account's agent-network endpoint (assigning a subdomain); the cluster
// is immutable thereafter.
type AgentNetworkSettings struct{}

// Annotate adds a description to the AgentNetworkSettings resource type.
func (a *AgentNetworkSettings) Annotate(ann infer.Annotator) {
	ann.Describe(&a, "Per-account NetBird Agent Network gateway settings. This is a singleton "+
		"resource — only one instance exists per account. Setting cluster bootstraps the "+
		"account's agent-network endpoint; the cluster is immutable thereafter.")
}

// AgentNetworkSettingsArgs defines input fields for Agent Network settings.
type AgentNetworkSettingsArgs struct {
	Cluster                *string `pulumi:"cluster,optional"`
	EnableLogCollection    bool    `pulumi:"enableLogCollection"`
	EnablePromptCollection bool    `pulumi:"enablePromptCollection"`
	RedactPii              bool    `pulumi:"redactPii"`
	AccessLogRetentionDays *int    `pulumi:"accessLogRetentionDays,optional"`
}

// Annotate provides documentation for AgentNetworkSettingsArgs fields.
func (a *AgentNetworkSettingsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Cluster, "Proxy cluster fronting this account's agent-network endpoint. "+
		"Bootstraps the account's settings row when set for the first time; immutable thereafter — "+
		"later changes must omit it or send the assigned value.")
	ann.Describe(&a.EnableLogCollection, "Whether per-request access-log entries are collected for this account's agent-network traffic.")
	ann.Describe(&a.EnablePromptCollection, "Master switch for request/response prompt capture. Capture runs only when this is on AND a policy guardrail also enables it.")
	ann.Describe(&a.RedactPii, "Whether captured prompts have PII redacted.")
	ann.Describe(&a.AccessLogRetentionDays, "Days to retain full access-log rows; older rows are swept. 0 or less means keep indefinitely.")
}

// AgentNetworkSettingsState represents the output state of Agent Network settings.
type AgentNetworkSettingsState struct {
	Cluster                *string `pulumi:"cluster,optional"`
	EnableLogCollection    bool    `pulumi:"enableLogCollection"`
	EnablePromptCollection bool    `pulumi:"enablePromptCollection"`
	RedactPii              bool    `pulumi:"redactPii"`
	AccessLogRetentionDays *int    `pulumi:"accessLogRetentionDays,optional"`
	Subdomain              *string `pulumi:"subdomain,optional"`
	Endpoint               *string `pulumi:"endpoint,optional"`
	CreatedAt              *string `pulumi:"createdAt,optional"`
	UpdatedAt              *string `pulumi:"updatedAt,optional"`
}

// Annotate provides documentation for AgentNetworkSettingsState fields.
func (a *AgentNetworkSettingsState) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Cluster, "Proxy cluster fronting this account's agent-network endpoint.")
	ann.Describe(&a.EnableLogCollection, "Whether per-request access-log entries are collected.")
	ann.Describe(&a.EnablePromptCollection, "Master switch for request/response prompt capture.")
	ann.Describe(&a.RedactPii, "Whether captured prompts have PII redacted.")
	ann.Describe(&a.AccessLogRetentionDays, "Days to retain full access-log rows.")
	ann.Describe(&a.Subdomain, "Auto-generated DNS-safe label that prefixes the cluster to form the agent-network endpoint. Empty until bootstrapped.")
	ann.Describe(&a.Endpoint, "Bare hostname agents call for this account. Empty until bootstrapped.")
	ann.Describe(&a.CreatedAt, "Timestamp when the settings row was created. Absent until bootstrapped.")
	ann.Describe(&a.UpdatedAt, "Timestamp when the settings row was last updated. Absent until bootstrapped.")
}

// Create initialises the Agent Network settings resource. Since this is a singleton, Create
// calls UpdateSettings and uses a fixed ID.
func (*AgentNetworkSettings) Create(
	ctx context.Context, req infer.CreateRequest[AgentNetworkSettingsArgs],
) (infer.CreateResponse[AgentNetworkSettingsState], error) {
	p.GetLogger(ctx).Debugf("Create:AgentNetworkSettings cluster=%v", req.Inputs.Cluster)

	if req.DryRun {
		return infer.CreateResponse[AgentNetworkSettingsState]{
			ID:     agentNetworkSettingsID,
			Output: agentNetworkSettingsStateFromArgs(req.Inputs, nil, nil, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.AgentNetwork.UpdateSettings(ctx, nbapi.AgentNetworkSettingsRequest{
		Cluster:                req.Inputs.Cluster,
		EnableLogCollection:    req.Inputs.EnableLogCollection,
		EnablePromptCollection: req.Inputs.EnablePromptCollection,
		RedactPii:              req.Inputs.RedactPii,
		AccessLogRetentionDays: req.Inputs.AccessLogRetentionDays,
	})
	if err != nil {
		return infer.CreateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("creating agent network settings failed: %w", err)
	}

	return infer.CreateResponse[AgentNetworkSettingsState]{
		ID:     agentNetworkSettingsID,
		Output: agentNetworkSettingsStateFromAPI(*updated),
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
				Inputs: AgentNetworkSettingsArgs{},  //nolint:exhaustruct
				State:  AgentNetworkSettingsState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[AgentNetworkSettingsArgs, AgentNetworkSettingsState]{}, fmt.Errorf("reading agent network settings failed: %w", err)
	}

	state := agentNetworkSettingsStateFromAPI(*settings)

	return infer.ReadResponse[AgentNetworkSettingsArgs, AgentNetworkSettingsState]{
		ID: req.ID,
		Inputs: AgentNetworkSettingsArgs{
			Cluster:                state.Cluster,
			EnableLogCollection:    settings.EnableLogCollection,
			EnablePromptCollection: settings.EnablePromptCollection,
			RedactPii:              settings.RedactPii,
			AccessLogRetentionDays: settings.AccessLogRetentionDays,
		},
		State: state,
	}, nil
}

// Update updates the Agent Network settings.
func (*AgentNetworkSettings) Update(
	ctx context.Context, req infer.UpdateRequest[AgentNetworkSettingsArgs, AgentNetworkSettingsState],
) (infer.UpdateResponse[AgentNetworkSettingsState], error) {
	p.GetLogger(ctx).Debugf("Update:AgentNetworkSettings[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[AgentNetworkSettingsState]{
			Output: agentNetworkSettingsStateFromArgs(
				req.Inputs, req.State.Subdomain, req.State.Endpoint, req.State.CreatedAt, req.State.UpdatedAt,
			),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.AgentNetwork.UpdateSettings(ctx, nbapi.AgentNetworkSettingsRequest{
		Cluster:                req.Inputs.Cluster,
		EnableLogCollection:    req.Inputs.EnableLogCollection,
		EnablePromptCollection: req.Inputs.EnablePromptCollection,
		RedactPii:              req.Inputs.RedactPii,
		AccessLogRetentionDays: req.Inputs.AccessLogRetentionDays,
	})
	if err != nil {
		return infer.UpdateResponse[AgentNetworkSettingsState]{}, fmt.Errorf("updating agent network settings failed: %w", err)
	}

	return infer.UpdateResponse[AgentNetworkSettingsState]{
		Output: agentNetworkSettingsStateFromAPI(*updated),
	}, nil
}

// Delete is a no-op because Agent Network settings are a singleton and cannot be deleted.
func (*AgentNetworkSettings) Delete(ctx context.Context, req infer.DeleteRequest[AgentNetworkSettingsState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:AgentNetworkSettings[%s] (no-op, singleton resource)", req.ID)

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between AgentNetworkSettingsArgs and AgentNetworkSettingsState.
func (*AgentNetworkSettings) Diff(
	ctx context.Context, req infer.DiffRequest[AgentNetworkSettingsArgs, AgentNetworkSettingsState],
) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:AgentNetworkSettings[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if !equalPtr(req.Inputs.Cluster, req.State.Cluster) {
		diff["cluster"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
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

	if !equalPtr(req.Inputs.AccessLogRetentionDays, req.State.AccessLogRetentionDays) {
		diff["accessLogRetentionDays"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:AgentNetworkSettings[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// WireDependencies explicitly defines input/output relationships.
func (*AgentNetworkSettings) WireDependencies(field infer.FieldSelector, args *AgentNetworkSettingsArgs, state *AgentNetworkSettingsState) {
	field.OutputField(&state.Cluster).DependsOn(field.InputField(&args.Cluster))
	field.OutputField(&state.EnableLogCollection).DependsOn(field.InputField(&args.EnableLogCollection))
	field.OutputField(&state.EnablePromptCollection).DependsOn(field.InputField(&args.EnablePromptCollection))
	field.OutputField(&state.RedactPii).DependsOn(field.InputField(&args.RedactPii))
	field.OutputField(&state.AccessLogRetentionDays).DependsOn(field.InputField(&args.AccessLogRetentionDays))
}

func agentNetworkSettingsStateFromArgs(
	args AgentNetworkSettingsArgs, subdomain, endpoint, createdAt, updatedAt *string,
) AgentNetworkSettingsState {
	return AgentNetworkSettingsState{
		Cluster:                args.Cluster,
		EnableLogCollection:    args.EnableLogCollection,
		EnablePromptCollection: args.EnablePromptCollection,
		RedactPii:              args.RedactPii,
		AccessLogRetentionDays: args.AccessLogRetentionDays,
		Subdomain:              subdomain,
		Endpoint:               endpoint,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}

func agentNetworkSettingsStateFromAPI(settings nbapi.AgentNetworkSettings) AgentNetworkSettingsState {
	var cluster, subdomain, endpoint *string

	if settings.Cluster != "" {
		cluster = &settings.Cluster
	}

	if settings.Subdomain != "" {
		subdomain = &settings.Subdomain
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
		Cluster:                cluster,
		EnableLogCollection:    settings.EnableLogCollection,
		EnablePromptCollection: settings.EnablePromptCollection,
		RedactPii:              settings.RedactPii,
		AccessLogRetentionDays: settings.AccessLogRetentionDays,
		Subdomain:              subdomain,
		Endpoint:               endpoint,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}
