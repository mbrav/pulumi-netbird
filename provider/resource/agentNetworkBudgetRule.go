package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// AgentNetworkBudgetRule represents an account-level Agent Network budget rule: a limit-only
// rule bound to groups and/or users that applies across all policies as a min-wins ceiling.
type AgentNetworkBudgetRule struct{}

// Annotate adds a description to the AgentNetworkBudgetRule resource type.
func (a *AgentNetworkBudgetRule) Annotate(ann infer.Annotator) {
	ann.Describe(&a, "A NetBird Agent Network account-level budget rule: a limit-only rule "+
		"bound to groups and/or users that applies across all policies as a min-wins ceiling. "+
		"Empty targets means it applies to every caller.")
}

// AgentNetworkBudgetRuleArgs defines input fields for an Agent Network budget rule.
type AgentNetworkBudgetRuleArgs struct {
	Name         string             `pulumi:"name"`
	Enabled      *bool              `pulumi:"enabled,optional"`
	Limits       AgentNetworkLimits `pulumi:"limits"`
	TargetGroups *[]string          `pulumi:"targetGroups,optional"`
	TargetUsers  *[]string          `pulumi:"targetUsers,optional"`
}

// Annotate provides documentation for AgentNetworkBudgetRuleArgs fields.
func (a *AgentNetworkBudgetRuleArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name for the budget rule.")
	ann.Describe(&a.Enabled, "Whether the rule is enforced.")
	ann.Describe(&a.Limits, "Token and budget caps attached directly to the rule.")
	ann.Describe(&a.TargetGroups, "NetBird group IDs the rule binds. Empty plus empty targetUsers means account-wide.")
	ann.Describe(&a.TargetUsers, "NetBird user IDs the rule binds directly.")
}

// AgentNetworkBudgetRuleState represents the output state of an Agent Network budget rule.
type AgentNetworkBudgetRuleState struct {
	Name         string             `pulumi:"name"`
	Enabled      *bool              `pulumi:"enabled,optional"`
	Limits       AgentNetworkLimits `pulumi:"limits"`
	TargetGroups *[]string          `pulumi:"targetGroups,optional"`
	TargetUsers  *[]string          `pulumi:"targetUsers,optional"`
	CreatedAt    *string            `pulumi:"createdAt,optional"`
	UpdatedAt    *string            `pulumi:"updatedAt,optional"`
}

// Annotate provides documentation for AgentNetworkBudgetRuleState fields.
func (a *AgentNetworkBudgetRuleState) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name for the budget rule.")
	ann.Describe(&a.Enabled, "Whether the rule is enforced.")
	ann.Describe(&a.Limits, "Token and budget caps attached directly to the rule.")
	ann.Describe(&a.TargetGroups, "NetBird group IDs the rule binds.")
	ann.Describe(&a.TargetUsers, "NetBird user IDs the rule binds directly.")
	ann.Describe(&a.CreatedAt, "Timestamp when the budget rule was created.")
	ann.Describe(&a.UpdatedAt, "Timestamp when the budget rule was last updated.")
}

// Create creates a new Agent Network budget rule.
func (*AgentNetworkBudgetRule) Create(
	ctx context.Context, req infer.CreateRequest[AgentNetworkBudgetRuleArgs],
) (infer.CreateResponse[AgentNetworkBudgetRuleState], error) {
	p.GetLogger(ctx).Debugf("Create:AgentNetworkBudgetRule name=%s", req.Inputs.Name)

	if req.DryRun {
		return infer.CreateResponse[AgentNetworkBudgetRuleState]{
			ID:     "preview",
			Output: agentNetworkBudgetRuleStateFromArgs(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[AgentNetworkBudgetRuleState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	limits := toAPIAgentNetworkLimits(req.Inputs.Limits)

	created, err := client.AgentNetwork.CreateBudgetRule(ctx, nbapi.AgentNetworkBudgetRuleRequest{
		Name:         req.Inputs.Name,
		Enabled:      req.Inputs.Enabled,
		Limits:       limits,
		TargetGroups: req.Inputs.TargetGroups,
		TargetUsers:  req.Inputs.TargetUsers,
	})
	if err != nil {
		return infer.CreateResponse[AgentNetworkBudgetRuleState]{}, fmt.Errorf("creating agent network budget rule failed: %w", err)
	}

	return infer.CreateResponse[AgentNetworkBudgetRuleState]{
		ID:     created.Id,
		Output: agentNetworkBudgetRuleStateFromAPI(*created),
	}, nil
}

// Read fetches the current state of an Agent Network budget rule from NetBird.
func (*AgentNetworkBudgetRule) Read(
	ctx context.Context, req infer.ReadRequest[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState],
) (infer.ReadResponse[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState], error) {
	p.GetLogger(ctx).Debugf("Read:AgentNetworkBudgetRule[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	rule, err := client.AgentNetwork.GetBudgetRule(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState]{
				ID:     "",
				Inputs: AgentNetworkBudgetRuleArgs{},  //nolint:exhaustruct
				State:  AgentNetworkBudgetRuleState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState]{}, fmt.Errorf("reading agent network budget rule failed: %w", err)
	}

	state := agentNetworkBudgetRuleStateFromAPI(*rule)

	return infer.ReadResponse[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState]{
		ID: req.ID,
		Inputs: AgentNetworkBudgetRuleArgs{
			Name:         rule.Name,
			Enabled:      &rule.Enabled,
			Limits:       state.Limits,
			TargetGroups: &rule.TargetGroups,
			TargetUsers:  &rule.TargetUsers,
		},
		State: state,
	}, nil
}

// Update updates an Agent Network budget rule.
func (*AgentNetworkBudgetRule) Update(
	ctx context.Context, req infer.UpdateRequest[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState],
) (infer.UpdateResponse[AgentNetworkBudgetRuleState], error) {
	p.GetLogger(ctx).Debugf("Update:AgentNetworkBudgetRule[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[AgentNetworkBudgetRuleState]{
			Output: agentNetworkBudgetRuleStateFromArgs(req.Inputs, req.State.CreatedAt, req.State.UpdatedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AgentNetworkBudgetRuleState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	limits := toAPIAgentNetworkLimits(req.Inputs.Limits)

	updated, err := client.AgentNetwork.UpdateBudgetRule(ctx, req.ID, nbapi.AgentNetworkBudgetRuleRequest{
		Name:         req.Inputs.Name,
		Enabled:      req.Inputs.Enabled,
		Limits:       limits,
		TargetGroups: req.Inputs.TargetGroups,
		TargetUsers:  req.Inputs.TargetUsers,
	})
	if err != nil {
		return infer.UpdateResponse[AgentNetworkBudgetRuleState]{}, fmt.Errorf("updating agent network budget rule failed: %w", err)
	}

	return infer.UpdateResponse[AgentNetworkBudgetRuleState]{
		Output: agentNetworkBudgetRuleStateFromAPI(*updated),
	}, nil
}

// Delete removes an Agent Network budget rule from NetBird.
func (*AgentNetworkBudgetRule) Delete(ctx context.Context, req infer.DeleteRequest[AgentNetworkBudgetRuleState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:AgentNetworkBudgetRule[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.AgentNetwork.DeleteBudgetRule(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting agent network budget rule failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*AgentNetworkBudgetRule) Diff(
	ctx context.Context, req infer.DiffRequest[AgentNetworkBudgetRuleArgs, AgentNetworkBudgetRuleState],
) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:AgentNetworkBudgetRule[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if boolVal(req.Inputs.Enabled) != boolVal(req.State.Enabled) {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Limits != req.State.Limits {
		diff["limits"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.TargetGroups, req.State.TargetGroups) {
		diff["targetGroups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.TargetUsers, req.State.TargetUsers) {
		diff["targetUsers"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:AgentNetworkBudgetRule[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*AgentNetworkBudgetRule) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[AgentNetworkBudgetRuleArgs], error) {
	p.GetLogger(ctx).Debugf("Check:AgentNetworkBudgetRule old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[AgentNetworkBudgetRuleArgs](ctx, req.NewInputs)

	if args.Enabled == nil {
		enabled := true
		args.Enabled = &enabled
	}

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{Property: "name", Reason: "name must not be empty"})
	}

	if args.TargetGroups != nil {
		for i, groupID := range *args.TargetGroups {
			if isBlank(groupID) {
				failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("targetGroups[%d]", i), Reason: "group id must not be empty"})
			}
		}
	}

	if args.TargetUsers != nil {
		for i, userID := range *args.TargetUsers {
			if isBlank(userID) {
				failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("targetUsers[%d]", i), Reason: "user id must not be empty"})
			}
		}
	}

	return infer.CheckResponse[AgentNetworkBudgetRuleArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*AgentNetworkBudgetRule) WireDependencies(field infer.FieldSelector, args *AgentNetworkBudgetRuleArgs, state *AgentNetworkBudgetRuleState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.Limits).DependsOn(field.InputField(&args.Limits))
	field.OutputField(&state.TargetGroups).DependsOn(field.InputField(&args.TargetGroups))
	field.OutputField(&state.TargetUsers).DependsOn(field.InputField(&args.TargetUsers))
}

func agentNetworkBudgetRuleStateFromArgs(args AgentNetworkBudgetRuleArgs, createdAt, updatedAt *string) AgentNetworkBudgetRuleState {
	return AgentNetworkBudgetRuleState{
		Name:         args.Name,
		Enabled:      args.Enabled,
		Limits:       args.Limits,
		TargetGroups: args.TargetGroups,
		TargetUsers:  args.TargetUsers,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func agentNetworkBudgetRuleStateFromAPI(rule nbapi.AgentNetworkBudgetRule) AgentNetworkBudgetRuleState {
	var createdAt, updatedAt *string

	if rule.CreatedAt != nil {
		formatted := rule.CreatedAt.Format(idpTimeFormat)
		createdAt = &formatted
	}

	if rule.UpdatedAt != nil {
		formatted := rule.UpdatedAt.Format(idpTimeFormat)
		updatedAt = &formatted
	}

	return AgentNetworkBudgetRuleState{
		Name:         rule.Name,
		Enabled:      &rule.Enabled,
		Limits:       fromAPIAgentNetworkLimits(rule.Limits),
		TargetGroups: &rule.TargetGroups,
		TargetUsers:  &rule.TargetUsers,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
