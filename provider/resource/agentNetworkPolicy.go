package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// AgentNetworkPolicy represents a NetBird Agent Network policy: which source groups may reach
// which providers, subject to attached guardrails and limits.
type AgentNetworkPolicy struct{}

// Annotate adds a description to the AgentNetworkPolicy resource type.
func (a *AgentNetworkPolicy) Annotate(ann infer.Annotator) {
	ann.Describe(&a, "A NetBird Agent Network policy: authorizes source groups to reach a set "+
		"of Agent Network providers, subject to attached guardrails and token/budget limits.")
}

// AgentNetworkPolicyArgs defines input fields for an Agent Network policy.
type AgentNetworkPolicyArgs struct {
	Name                   string              `pulumi:"name"`
	Description            *string             `pulumi:"description,optional"`
	Enabled                *bool               `pulumi:"enabled,optional"`
	SourceGroups           []string            `pulumi:"sourceGroups"`
	DestinationProviderIDs []string            `pulumi:"destinationProviderIds"`
	GuardrailIDs           *[]string           `pulumi:"guardrailIds,optional"`
	Limits                 *AgentNetworkLimits `pulumi:"limits,optional"`
}

// Annotate provides documentation for AgentNetworkPolicyArgs fields.
func (a *AgentNetworkPolicyArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name for the policy.")
	ann.Describe(&a.Description, "Optional human-readable description.")
	ann.Describe(&a.Enabled, "Whether the policy is enabled.")
	ann.Describe(&a.SourceGroups, "NetBird group IDs whose members are allowed to call the destination providers.")
	ann.Describe(&a.DestinationProviderIDs, "AgentNetworkProvider IDs the source groups can reach.")
	ann.Describe(&a.GuardrailIDs, "AgentNetworkGuardrail IDs to attach to this policy.")
	ann.Describe(&a.Limits, "Token and budget caps attached directly to the policy. These compose with any guardrail-level checks.")
}

// AgentNetworkPolicyState represents the output state of an Agent Network policy.
type AgentNetworkPolicyState struct {
	Name                   string              `pulumi:"name"`
	Description            *string             `pulumi:"description,optional"`
	Enabled                *bool               `pulumi:"enabled,optional"`
	SourceGroups           []string            `pulumi:"sourceGroups"`
	DestinationProviderIDs []string            `pulumi:"destinationProviderIds"`
	GuardrailIDs           *[]string           `pulumi:"guardrailIds,optional"`
	Limits                 *AgentNetworkLimits `pulumi:"limits,optional"`
	CreatedAt              *string             `pulumi:"createdAt,optional"`
	UpdatedAt              *string             `pulumi:"updatedAt,optional"`
}

// Annotate provides documentation for AgentNetworkPolicyState fields.
func (a *AgentNetworkPolicyState) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name for the policy.")
	ann.Describe(&a.Description, "Optional human-readable description.")
	ann.Describe(&a.Enabled, "Whether the policy is enabled.")
	ann.Describe(&a.SourceGroups, "NetBird group IDs whose members are allowed to call the destination providers.")
	ann.Describe(&a.DestinationProviderIDs, "AgentNetworkProvider IDs the source groups can reach.")
	ann.Describe(&a.GuardrailIDs, "AgentNetworkGuardrail IDs attached to this policy.")
	ann.Describe(&a.Limits, "Token and budget caps attached directly to the policy.")
	ann.Describe(&a.CreatedAt, "Timestamp when the policy was created.")
	ann.Describe(&a.UpdatedAt, "Timestamp when the policy was last updated.")
}

// Create creates a new Agent Network policy.
func (*AgentNetworkPolicy) Create(
	ctx context.Context, req infer.CreateRequest[AgentNetworkPolicyArgs],
) (infer.CreateResponse[AgentNetworkPolicyState], error) {
	p.GetLogger(ctx).Debugf("Create:AgentNetworkPolicy name=%s", req.Inputs.Name)

	if req.DryRun {
		return infer.CreateResponse[AgentNetworkPolicyState]{
			ID:     "preview",
			Output: agentNetworkPolicyStateFromArgs(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[AgentNetworkPolicyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// The API always echoes description on the wire and treats "" as unset, so an
	// omitted input travels as "" rather than being left out: that keeps inputs and
	// state in agreement and makes removing the field clear it server-side.
	description := strPtr(req.Inputs.Description)

	created, err := client.AgentNetwork.CreatePolicy(ctx, nbapi.AgentNetworkPolicyRequest{
		Name:                   req.Inputs.Name,
		Description:            &description,
		Enabled:                req.Inputs.Enabled,
		SourceGroups:           req.Inputs.SourceGroups,
		DestinationProviderIds: req.Inputs.DestinationProviderIDs,
		GuardrailIds:           req.Inputs.GuardrailIDs,
		Limits:                 toAPIAgentNetworkLimitsPtr(req.Inputs.Limits),
	})
	if err != nil {
		return infer.CreateResponse[AgentNetworkPolicyState]{}, fmt.Errorf("creating agent network policy failed: %w", err)
	}

	return infer.CreateResponse[AgentNetworkPolicyState]{
		ID:     created.Id,
		Output: agentNetworkPolicyStateFromAPI(*created),
	}, nil
}

// Read fetches the current state of an Agent Network policy from NetBird.
func (*AgentNetworkPolicy) Read(
	ctx context.Context, req infer.ReadRequest[AgentNetworkPolicyArgs, AgentNetworkPolicyState],
) (infer.ReadResponse[AgentNetworkPolicyArgs, AgentNetworkPolicyState], error) {
	p.GetLogger(ctx).Debugf("Read:AgentNetworkPolicy[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[AgentNetworkPolicyArgs, AgentNetworkPolicyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	policy, err := client.AgentNetwork.GetPolicy(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[AgentNetworkPolicyArgs, AgentNetworkPolicyState]{
				ID:     "",
				Inputs: AgentNetworkPolicyArgs{},  //nolint:exhaustruct
				State:  AgentNetworkPolicyState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[AgentNetworkPolicyArgs, AgentNetworkPolicyState]{}, fmt.Errorf("reading agent network policy failed: %w", err)
	}

	state := agentNetworkPolicyStateFromAPI(*policy)

	return infer.ReadResponse[AgentNetworkPolicyArgs, AgentNetworkPolicyState]{
		ID: req.ID,
		Inputs: AgentNetworkPolicyArgs{
			Name:                   policy.Name,
			Description:            &policy.Description,
			Enabled:                &policy.Enabled,
			SourceGroups:           policy.SourceGroups,
			DestinationProviderIDs: policy.DestinationProviderIds,
			GuardrailIDs:           &policy.GuardrailIds,
			Limits:                 state.Limits,
		},
		State: state,
	}, nil
}

// Update updates an Agent Network policy.
func (*AgentNetworkPolicy) Update(
	ctx context.Context, req infer.UpdateRequest[AgentNetworkPolicyArgs, AgentNetworkPolicyState],
) (infer.UpdateResponse[AgentNetworkPolicyState], error) {
	p.GetLogger(ctx).Debugf("Update:AgentNetworkPolicy[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[AgentNetworkPolicyState]{
			Output: agentNetworkPolicyStateFromArgs(req.Inputs, req.State.CreatedAt, req.State.UpdatedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AgentNetworkPolicyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// The API always echoes description on the wire and treats "" as unset, so an
	// omitted input travels as "" rather than being left out: that keeps inputs and
	// state in agreement and makes removing the field clear it server-side.
	description := strPtr(req.Inputs.Description)

	updated, err := client.AgentNetwork.UpdatePolicy(ctx, req.ID, nbapi.AgentNetworkPolicyRequest{
		Name:                   req.Inputs.Name,
		Description:            &description,
		Enabled:                req.Inputs.Enabled,
		SourceGroups:           req.Inputs.SourceGroups,
		DestinationProviderIds: req.Inputs.DestinationProviderIDs,
		GuardrailIds:           req.Inputs.GuardrailIDs,
		Limits:                 toAPIAgentNetworkLimitsPtr(req.Inputs.Limits),
	})
	if err != nil {
		return infer.UpdateResponse[AgentNetworkPolicyState]{}, fmt.Errorf("updating agent network policy failed: %w", err)
	}

	return infer.UpdateResponse[AgentNetworkPolicyState]{
		Output: agentNetworkPolicyStateFromAPI(*updated),
	}, nil
}

// Delete removes an Agent Network policy from NetBird.
func (*AgentNetworkPolicy) Delete(ctx context.Context, req infer.DeleteRequest[AgentNetworkPolicyState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:AgentNetworkPolicy[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.AgentNetwork.DeletePolicy(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting agent network policy failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*AgentNetworkPolicy) Diff(
	ctx context.Context, req infer.DiffRequest[AgentNetworkPolicyArgs, AgentNetworkPolicyState],
) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:AgentNetworkPolicy[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalOptionalStr(req.Inputs.Description, req.State.Description) {
		diff["description"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if boolVal(req.Inputs.Enabled) != boolVal(req.State.Enabled) {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlice(req.Inputs.SourceGroups, req.State.SourceGroups) {
		diff["sourceGroups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlice(req.Inputs.DestinationProviderIDs, req.State.DestinationProviderIDs) {
		diff["destinationProviderIds"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.GuardrailIDs, req.State.GuardrailIDs) {
		diff["guardrailIds"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalAgentNetworkLimitsPtr(req.Inputs.Limits, req.State.Limits) {
		diff["limits"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:AgentNetworkPolicy[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*AgentNetworkPolicy) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[AgentNetworkPolicyArgs], error) {
	p.GetLogger(ctx).Debugf("Check:AgentNetworkPolicy old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[AgentNetworkPolicyArgs](ctx, req.NewInputs)

	if args.Enabled == nil {
		enabled := true
		args.Enabled = &enabled
	}

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{Property: "name", Reason: "name must not be empty"})
	}

	for i, groupID := range args.SourceGroups {
		if isBlank(groupID) {
			failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("sourceGroups[%d]", i), Reason: "group id must not be empty"})
		}
	}

	for i, providerID := range args.DestinationProviderIDs {
		if isBlank(providerID) {
			failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("destinationProviderIds[%d]", i), Reason: "provider id must not be empty"})
		}
	}

	if args.GuardrailIDs != nil {
		for i, guardrailID := range *args.GuardrailIDs {
			if isBlank(guardrailID) {
				failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("guardrailIds[%d]", i), Reason: "guardrail id must not be empty"})
			}
		}
	}

	return infer.CheckResponse[AgentNetworkPolicyArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*AgentNetworkPolicy) WireDependencies(field infer.FieldSelector, args *AgentNetworkPolicyArgs, state *AgentNetworkPolicyState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Description).DependsOn(field.InputField(&args.Description))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.SourceGroups).DependsOn(field.InputField(&args.SourceGroups))
	field.OutputField(&state.DestinationProviderIDs).DependsOn(field.InputField(&args.DestinationProviderIDs))
	field.OutputField(&state.GuardrailIDs).DependsOn(field.InputField(&args.GuardrailIDs))
	field.OutputField(&state.Limits).DependsOn(field.InputField(&args.Limits))
}

func agentNetworkPolicyStateFromArgs(args AgentNetworkPolicyArgs, createdAt, updatedAt *string) AgentNetworkPolicyState {
	return AgentNetworkPolicyState{
		Name:                   args.Name,
		Description:            args.Description,
		Enabled:                args.Enabled,
		SourceGroups:           args.SourceGroups,
		DestinationProviderIDs: args.DestinationProviderIDs,
		GuardrailIDs:           args.GuardrailIDs,
		Limits:                 args.Limits,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}

func agentNetworkPolicyStateFromAPI(policy nbapi.AgentNetworkPolicy) AgentNetworkPolicyState {
	var createdAt, updatedAt *string

	if policy.CreatedAt != nil {
		formatted := policy.CreatedAt.Format(idpTimeFormat)
		createdAt = &formatted
	}

	if policy.UpdatedAt != nil {
		formatted := policy.UpdatedAt.Format(idpTimeFormat)
		updatedAt = &formatted
	}

	limits := fromAPIAgentNetworkLimits(policy.Limits)

	return AgentNetworkPolicyState{
		Name:                   policy.Name,
		Description:            &policy.Description,
		Enabled:                &policy.Enabled,
		SourceGroups:           policy.SourceGroups,
		DestinationProviderIDs: policy.DestinationProviderIds,
		GuardrailIDs:           &policy.GuardrailIds,
		Limits:                 &limits,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}

func toAPIAgentNetworkLimitsPtr(limits *AgentNetworkLimits) *nbapi.AgentNetworkPolicyLimits {
	if limits == nil {
		return nil
	}

	converted := toAPIAgentNetworkLimits(*limits)

	return &converted
}
