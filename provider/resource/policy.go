package resource

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// TEST: InputDiff: false
// TODO: Implement AuthorizedGroups

// Policy defines the Pulumi resource handler for NetBird policy resources.
type Policy struct{}

// Annotate adds a description annotation for the Policy type for generated SDKs.
func (policy *Policy) Annotate(annotator infer.Annotator) {
	annotator.Describe(policy, "A NetBird policy defining rules for communication between peers.")
}

// PolicyArgs defines the user-supplied arguments for creating/updating a Policy resource.
type PolicyArgs struct {
	Name                string           `pulumi:"name"`                   // Policy name (required)
	Description         *string          `pulumi:"description,optional"`   // Optional description
	Enabled             bool             `pulumi:"enabled"`                // Whether the policy is enabled
	Rules               []PolicyRuleArgs `pulumi:"rules"`                  // List of rules defined in the policy
	SourcePostureChecks *[]string        `pulumi:"postureChecks,optional"` // Optional list of posture check IDs
}

// Annotate adds field descriptions to PolicyArgs for generated SDKs.
func (policy *PolicyArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&policy.Name, "Name Policy name identifier")
	annotator.Describe(&policy.Description, "Description Policy friendly description, optional")
	annotator.Describe(&policy.Enabled, "Enabled Policy status")
	annotator.Describe(&policy.Rules, "Rules Policy rule object for policy UI editor")
	annotator.Describe(&policy.SourcePostureChecks, "SourcePostureChecks Posture checks ID's applied to policy source groups, optional")
}

// PolicyState represents the state of a Policy resource stored in Pulumi state.
type PolicyState struct {
	Name                string            `pulumi:"name"`
	Description         *string           `pulumi:"description,optional"`
	Enabled             bool              `pulumi:"enabled"`
	Rules               []PolicyRuleState `pulumi:"rules"`
	SourcePostureChecks *[]string         `pulumi:"postureChecks,optional"`
}

// Annotate adds descriptive annotations to the PolicyState fields for use in generated SDKs.
func (policy *PolicyState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&policy.Name, "Name Policy name identifier")
	annotator.Describe(&policy.Description, "Description Policy friendly description, optional")
	annotator.Describe(&policy.Enabled, "Enabled Policy status")
	annotator.Describe(&policy.Rules, "Rules Policy rule object for policy UI editor")
	annotator.Describe(&policy.SourcePostureChecks, "SourcePostureChecks Posture checks ID's applied to policy source groups, optional")
}

// PolicyRuleArgs represents user input for an individual rule in a policy.
type PolicyRuleArgs struct {
	ID                  *string          `pulumi:"id,optional"`                  // Optional rule ID (used for updates)
	Name                string           `pulumi:"name"`                         // Rule name
	Description         *string          `pulumi:"description,optional"`         // Optional rule description
	Bidirectional       bool             `pulumi:"bidirectional"`                // Whether the rule is bidirectional
	Action              RuleAction       `pulumi:"action"`                       // Rule action (accept/drop)
	Enabled             bool             `pulumi:"enabled"`                      // Whether the rule is enabled
	Protocol            Protocol         `pulumi:"protocol"`                     // Network protocol
	Ports               *[]string        `pulumi:"ports,optional"`               // Optional list of specific ports
	PortRanges          *[]RulePortRange `pulumi:"portRanges,optional"`          // Optional list of port ranges
	Sources             *[]string        `pulumi:"sources,optional"`             // Optional list of source group IDs
	Destinations        *[]string        `pulumi:"destinations,optional"`        // Optional list of destination group IDs
	SourceResource      *Resource        `pulumi:"sourceResource,optional"`      // Optional single source resource
	DestinationResource *Resource        `pulumi:"destinationResource,optional"` // Optional single destination resource
}

// Annotate adds descriptive annotations to the PolicyRuleArgs fields for use in generated SDKs.
func (policy *PolicyRuleArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&policy.ID, "ID Policy rule.")
	annotator.Describe(&policy.Name, "Name Policy rule name identifier")
	annotator.Describe(&policy.Description, "Description Policy rule friendly description")
	annotator.Describe(&policy.Bidirectional, "Bidirectional Define if the rule is applicable in both directions, sources, and destinations.")
	annotator.Describe(&policy.Action, "Action Policy rule accept or drops packets")
	annotator.Describe(&policy.Enabled, "Enabled Policy rule status")
	annotator.Describe(&policy.Protocol, "Protocol Policy rule type of the traffic")
	annotator.Describe(&policy.Ports, "Ports Policy rule affected ports")
	annotator.Describe(&policy.PortRanges, "PortRanges Policy rule affected ports ranges list")
	annotator.Describe(&policy.Sources, "Sources Policy rule source group IDs")
	annotator.Describe(&policy.Destinations, "Destinations Policy rule destination group IDs")
	annotator.Describe(&policy.SourceResource, "SourceResource for the rule")
	annotator.Describe(&policy.DestinationResource, "DestinationResource for the rule ")
}

// PolicyRuleState represents the state of an individual rule within a policy.
type PolicyRuleState struct {
	ID                  *string          `pulumi:"id,optional"`
	Name                string           `pulumi:"name"`
	Description         *string          `pulumi:"description,optional"`
	Bidirectional       bool             `pulumi:"bidirectional"`
	Action              RuleAction       `pulumi:"action"`
	Enabled             bool             `pulumi:"enabled"`
	Protocol            Protocol         `pulumi:"protocol"`
	Ports               *[]string        `pulumi:"ports,optional"`
	PortRanges          *[]RulePortRange `pulumi:"portRanges,optional"`
	Sources             *[]RuleGroup     `pulumi:"sources,optional"` // Fully-resolved group info (not just IDs)
	Destinations        *[]RuleGroup     `pulumi:"destinations,optional"`
	SourceResource      *Resource        `pulumi:"sourceResource,optional"`
	DestinationResource *Resource        `pulumi:"destinationResource,optional"`
}

// Annotate adds descriptive annotations to the PolicyRuleState fields for use in generated SDKs.
func (policy *PolicyRuleState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&policy.ID, "ID Policy rule.")
	annotator.Describe(&policy.Name, "Name Policy rule name identifier")
	annotator.Describe(&policy.Description, "Description Policy rule friendly description")
	annotator.Describe(&policy.Bidirectional, "Bidirectional Define if the rule is applicable in both directions, sources, and destinations.")
	annotator.Describe(&policy.Action, "Action Policy rule accept or drops packets")
	annotator.Describe(&policy.Enabled, "Enabled Policy rule status")
	annotator.Describe(&policy.Protocol, "Protocol Policy rule type of the traffic")
	annotator.Describe(&policy.Ports, "Ports Policy rule affected ports")
	annotator.Describe(&policy.PortRanges, "PortRanges Policy rule affected ports ranges list")
	annotator.Describe(&policy.Sources, "Sources Policy rule source group IDs")
	annotator.Describe(&policy.Destinations, "Destinations Policy rule destination group IDs")
	annotator.Describe(&policy.SourceResource, "SourceResource for the rule")
	annotator.Describe(&policy.DestinationResource, "DestinationResource for the rule ")
}

// RulePortRange type.
type RulePortRange struct {
	Start int `pulumi:"start"`
	End   int `pulumi:"end"`
}

// Annotate adds descriptive annotations to the RulePortRange fields for use in generated SDKs.
func (r *RulePortRange) Annotate(a infer.Annotator) {
	a.Describe(&r.Start, "Start of port range")
	a.Describe(&r.End, "End of port range")
}

// RuleGroup type.
type RuleGroup struct {
	ID   string `pulumi:"id"`
	Name string `pulumi:"name"`
}

// Annotate adds descriptive annotations to the RuleGroup fields for use in generated SDKs.
func (g *RuleGroup) Annotate(a infer.Annotator) {
	a.Describe(&g.ID, "The unique identifier of the group.")
	a.Describe(&g.Name, "The name of the group.")
}

// RuleAction defines the allowed actions for a rule (accept/drop).
// This wraps the nbapi type to allow method definitions (like Values()).
type RuleAction string

// RuleActi2yyonAccept and RuleActionDrop represent possible actions for a policy rule.
// RuleActionAccept allows traffic, while RuleActionDrop blocks traffic.
const (
	RuleActionAccept RuleAction = RuleAction(nbapi.PolicyRuleActionAccept)
	RuleActionDrop   RuleAction = RuleAction(nbapi.PolicyRuleActionDrop)
)

// Values returns the valid enum values for RuleAction, used by Pulumi for schema generation and validation.
func (RuleAction) Values() []infer.EnumValue[RuleAction] {
	return []infer.EnumValue[RuleAction]{
		{Name: "Accept", Value: RuleActionAccept, Description: "Accept action"},
		{Name: "Drop", Value: RuleActionDrop, Description: "Drop action"},
	}
}

// Protocol defines the allowed network protocols for a policy rule.
type Protocol string

// Enum constants for supported network protocols.
const (
	ProtocolAll  Protocol = Protocol(nbapi.PolicyRuleProtocolAll)
	ProtocolIcmp Protocol = Protocol(nbapi.PolicyRuleProtocolIcmp)
	ProtocolTCP  Protocol = Protocol(nbapi.PolicyRuleProtocolTcp)
	ProtocolUDP  Protocol = Protocol(nbapi.PolicyRuleProtocolUdp)
)

// Values returns valid protocol values for Pulumi enum support.
func (Protocol) Values() []infer.EnumValue[Protocol] {
	return []infer.EnumValue[Protocol]{
		{Name: "All", Value: ProtocolAll, Description: "All protocols"},
		{Name: "ICMP", Value: ProtocolIcmp, Description: "ICMP protocol"},
		{Name: "TCP", Value: ProtocolTCP, Description: "TCP protocol"},
		{Name: "UDP", Value: ProtocolUDP, Description: "UDP protocol"},
	}
}

// Create creates a new NetBird policy.
func (*Policy) Create(ctx context.Context, req infer.CreateRequest[PolicyArgs]) (infer.CreateResponse[PolicyState], error) {
	p.GetLogger(ctx).Debugf("Create:Policy")

	// Handle dry-run (preview) mode by constructing a preview PolicyState.
	if req.DryRun {
		// Convert PolicyRuleArgs to PolicyRuleState for preview
		rules := make([]PolicyRuleState, len(req.Inputs.Rules))

		for ruleIndex, rule := range req.Inputs.Rules {
			// Construct sources and destination groups
			var sources, destinations *[]RuleGroup

			if rule.Sources != nil {
				groups := make([]RuleGroup, len(*rule.Sources))
				for j, g := range *rule.Sources {
					groups[j] = RuleGroup{Name: "preview", ID: g}
				}
				// Assign the address of the groups slice to sources
				sources = &groups
			}

			if rule.Destinations != nil {
				groups := make([]RuleGroup, len(*rule.Destinations))
				for j, g := range *rule.Destinations {
					groups[j] = RuleGroup{Name: "preview", ID: g}
				}
				// Assign the address of the groups slice to destinations
				destinations = &groups
			}

			rules[ruleIndex] = PolicyRuleState{
				ID:                  rule.ID,
				Name:                rule.Name,
				Description:         rule.Description,
				Bidirectional:       rule.Bidirectional,
				Action:              rule.Action,
				Enabled:             rule.Enabled,
				Protocol:            rule.Protocol,
				Ports:               rule.Ports,
				PortRanges:          rule.PortRanges,
				Sources:             sources,
				Destinations:        destinations,
				SourceResource:      rule.SourceResource,
				DestinationResource: rule.DestinationResource,
			}
		}

		return infer.CreateResponse[PolicyState]{
			ID: "preview",
			Output: PolicyState{
				Name:                req.Inputs.Name,
				Description:         req.Inputs.Description,
				Enabled:             req.Inputs.Enabled,
				Rules:               rules,
				SourcePostureChecks: req.Inputs.SourcePostureChecks,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[PolicyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// Convert input rules to nbapi.PolicyRuleUpdate
	apiRules := make([]nbapi.PolicyRuleUpdate, len(req.Inputs.Rules))
	for i, rule := range req.Inputs.Rules {
		apiRules[i] = nbapi.PolicyRuleUpdate{
			Id:                  rule.ID,
			Name:                rule.Name,
			Description:         rule.Description,
			Bidirectional:       rule.Bidirectional,
			Action:              nbapi.PolicyRuleUpdateAction(rule.Action),
			Enabled:             rule.Enabled,
			Protocol:            nbapi.PolicyRuleUpdateProtocol(rule.Protocol),
			Ports:               rule.Ports,
			PortRanges:          toAPIPortRanges(rule.PortRanges),
			Sources:             rule.Sources,
			Destinations:        rule.Destinations,
			AuthorizedGroups:    nil, // TODO: Implement AuthorizedGroups
			SourceResource:      toAPIResource(rule.SourceResource),
			DestinationResource: toAPIResource(rule.DestinationResource),
		}
	}

	created, err := client.Policies.Create(ctx, nbapi.PolicyUpdate{
		Name:                req.Inputs.Name,
		Description:         req.Inputs.Description,
		Enabled:             req.Inputs.Enabled,
		Rules:               apiRules,
		SourcePostureChecks: req.Inputs.SourcePostureChecks,
	})
	if err != nil {
		return infer.CreateResponse[PolicyState]{}, fmt.Errorf("creating policy failed: %w", err)
	}

	if created.Id == nil || *created.Id == "" {
		return infer.CreateResponse[PolicyState]{}, errors.New("policy create response did not include an id")
	}

	// Convert created rules to PolicyRuleState
	rules := make([]PolicyRuleState, len(created.Rules))
	for ruleIndex, rule := range created.Rules {
		rules[ruleIndex] = PolicyRuleState{
			ID:                  rule.Id,
			Name:                rule.Name,
			Description:         rule.Description,
			Bidirectional:       rule.Bidirectional,
			Action:              RuleAction(rule.Action),
			Enabled:             rule.Enabled,
			Protocol:            Protocol(rule.Protocol),
			Ports:               rule.Ports,
			PortRanges:          fromAPIPortRanges(rule.PortRanges),
			Sources:             fromAPIGroupMinimums(rule.Sources),
			Destinations:        fromAPIGroupMinimums(rule.Destinations),
			SourceResource:      fromAPIResource(rule.SourceResource),
			DestinationResource: fromAPIResource(rule.DestinationResource),
		}
	}

	return infer.CreateResponse[PolicyState]{
		ID: *created.Id,
		Output: PolicyState{
			Name:                created.Name,
			Description:         created.Description,
			Enabled:             created.Enabled,
			Rules:               rules,
			SourcePostureChecks: &created.SourcePostureChecks,
		},
	}, nil
}

// Read reads a Policy from NetBird.
func (*Policy) Read(ctx context.Context, req infer.ReadRequest[PolicyArgs, PolicyState]) (infer.ReadResponse[PolicyArgs, PolicyState], error) {
	p.GetLogger(ctx).Debugf("Read:Policy[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[PolicyArgs, PolicyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	policy, err := client.Policies.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[PolicyArgs, PolicyState]{
				ID:     "",
				Inputs: PolicyArgs{},  //nolint:exhaustruct
				State:  PolicyState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[PolicyArgs, PolicyState]{}, fmt.Errorf("reading policy failed: %w", err)
	}

	rules := make([]PolicyRuleState, len(policy.Rules))
	for ruleIndex, rule := range policy.Rules {
		rules[ruleIndex] = PolicyRuleState{
			ID:                  rule.Id,
			Name:                rule.Name,
			Description:         rule.Description,
			Bidirectional:       rule.Bidirectional,
			Action:              RuleAction(rule.Action),
			Enabled:             rule.Enabled,
			Protocol:            Protocol(rule.Protocol),
			Ports:               rule.Ports,
			PortRanges:          fromAPIPortRanges(rule.PortRanges),
			Sources:             fromAPIGroupMinimums(rule.Sources),
			Destinations:        fromAPIGroupMinimums(rule.Destinations),
			SourceResource:      fromAPIResource(rule.SourceResource),
			DestinationResource: fromAPIResource(rule.DestinationResource),
		}
	}

	inputRules := req.Inputs.Rules
	if len(inputRules) == 0 && len(policy.Rules) > 0 {
		inputRules = policyRulesToArgs(policy.Rules)
	}

	postureChecks := slices.Clone(policy.SourcePostureChecks)
	slices.Sort(postureChecks)

	var stateDescription *string
	if req.Inputs.Description != nil {
		stateDescription = policy.Description
	}

	return infer.ReadResponse[PolicyArgs, PolicyState]{
		ID: req.ID,
		Inputs: PolicyArgs{
			Name:                policy.Name,
			Description:         req.Inputs.Description,
			Enabled:             policy.Enabled,
			Rules:               inputRules,
			SourcePostureChecks: &postureChecks,
		},
		State: PolicyState{
			Name:                policy.Name,
			Description:         stateDescription,
			Enabled:             policy.Enabled,
			Rules:               rules,
			SourcePostureChecks: &postureChecks,
		},
	}, nil
}

// Update updates an existing NetBird policy.
func (*Policy) Update(ctx context.Context, req infer.UpdateRequest[PolicyArgs, PolicyState]) (infer.UpdateResponse[PolicyState], error) {
	p.GetLogger(ctx).Debugf("Update:Policy[%s]", req.ID)

	if req.DryRun {
		// Construct PolicyRuleState for preview output
		rules := make([]PolicyRuleState, len(req.Inputs.Rules))

		for ruleIndex, rule := range req.Inputs.Rules {
			var sources, destinations *[]RuleGroup

			if rule.Sources != nil {
				groups := make([]RuleGroup, len(*rule.Sources))
				for j, g := range *rule.Sources {
					groups[j] = RuleGroup{Name: "preview", ID: g}
				}

				sources = &groups
			}

			if rule.Destinations != nil {
				groups := make([]RuleGroup, len(*rule.Destinations))
				for j, g := range *rule.Destinations {
					groups[j] = RuleGroup{Name: "preview", ID: g}
				}

				destinations = &groups
			}

			rules[ruleIndex] = PolicyRuleState{
				ID:                  rule.ID,
				Name:                rule.Name,
				Description:         rule.Description,
				Bidirectional:       rule.Bidirectional,
				Action:              rule.Action,
				Enabled:             rule.Enabled,
				Protocol:            rule.Protocol,
				Ports:               rule.Ports,
				PortRanges:          rule.PortRanges,
				Sources:             sources,
				Destinations:        destinations,
				SourceResource:      rule.SourceResource,
				DestinationResource: rule.DestinationResource,
			}
		}

		return infer.UpdateResponse[PolicyState]{
			Output: PolicyState{
				Name:                req.Inputs.Name,
				Description:         req.Inputs.Description,
				Enabled:             req.Inputs.Enabled,
				Rules:               rules,
				SourcePostureChecks: req.Inputs.SourcePostureChecks,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[PolicyState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// Convert input rules to nbapi.PolicyRuleUpdate
	apiRules := make([]nbapi.PolicyRuleUpdate, len(req.Inputs.Rules))
	for ruleIndex, rule := range req.Inputs.Rules {
		apiRules[ruleIndex] = nbapi.PolicyRuleUpdate{
			Id:                  rule.ID,
			Name:                rule.Name,
			Description:         rule.Description,
			Bidirectional:       rule.Bidirectional,
			Action:              nbapi.PolicyRuleUpdateAction(rule.Action),
			Enabled:             rule.Enabled,
			Protocol:            nbapi.PolicyRuleUpdateProtocol(rule.Protocol),
			Ports:               rule.Ports,
			PortRanges:          toAPIPortRanges(rule.PortRanges),
			Sources:             rule.Sources,
			Destinations:        rule.Destinations,
			AuthorizedGroups:    nil, // TODO: Implement AuthorizedGroups
			SourceResource:      toAPIResource(rule.SourceResource),
			DestinationResource: toAPIResource(rule.DestinationResource),
		}
	}

	updated, err := client.Policies.Update(ctx, req.ID, nbapi.PolicyCreate{
		Name:                req.Inputs.Name,
		Description:         req.Inputs.Description,
		Enabled:             req.Inputs.Enabled,
		Rules:               apiRules,
		SourcePostureChecks: req.Inputs.SourcePostureChecks,
	})
	if err != nil {
		return infer.UpdateResponse[PolicyState]{}, fmt.Errorf("updating policy %s failed: %w", req.Inputs.Name, err)
	}

	// Convert updated rules to PolicyRuleState
	rules := make([]PolicyRuleState, len(updated.Rules))
	for ruleIndex, rule := range updated.Rules {
		rules[ruleIndex] = PolicyRuleState{
			ID:                  rule.Id,
			Name:                rule.Name,
			Description:         rule.Description,
			Bidirectional:       rule.Bidirectional,
			Action:              RuleAction(rule.Action),
			Enabled:             rule.Enabled,
			Protocol:            Protocol(rule.Protocol),
			Ports:               rule.Ports,
			PortRanges:          fromAPIPortRanges(rule.PortRanges),
			Sources:             fromAPIGroupMinimums(rule.Sources),
			Destinations:        fromAPIGroupMinimums(rule.Destinations),
			SourceResource:      fromAPIResource(rule.SourceResource),
			DestinationResource: fromAPIResource(rule.DestinationResource),
		}
	}

	return infer.UpdateResponse[PolicyState]{
		Output: PolicyState{
			Name:                updated.Name,
			Description:         updated.Description,
			Enabled:             updated.Enabled,
			Rules:               rules,
			SourcePostureChecks: &updated.SourcePostureChecks,
		},
	}, nil
}

// Delete removes a Policy from NetBird.
func (*Policy) Delete(ctx context.Context, req infer.DeleteRequest[PolicyState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:Policy[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Policies.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting policy failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*Policy) Diff(ctx context.Context, req infer.DiffRequest[PolicyArgs, PolicyState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Policy[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Description != nil && !equalPtr(req.Inputs.Description, req.State.Description) {
		diff["description"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}
	// Rules Diff
	if len(req.Inputs.Rules) != len(req.State.Rules) {
		diff["rules"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	} else {
		equal := true

		for ruleIndex := range req.Inputs.Rules {
			input := req.Inputs.Rules[ruleIndex]
			state := req.State.Rules[ruleIndex]

			p.GetLogger(ctx).Debugf("Diff:Policy[%s]:Rules[%d] a=%+v b=%+v", req.ID, ruleIndex, input, state)

			if input.Name != state.Name ||
				(input.Description != nil && !equalPtr(input.Description, state.Description)) ||
				input.Bidirectional != state.Bidirectional ||
				input.Action != state.Action ||
				input.Enabled != state.Enabled ||
				input.Protocol != state.Protocol ||
				!equalSlicePtr(input.Ports, state.Ports) ||
				!equalPortRangePtr(input.PortRanges, state.PortRanges) ||
				!equalSlicePtr(input.Sources, toGroupIDs(state.Sources)) ||
				!equalSlicePtr(input.Destinations, toGroupIDs(state.Destinations)) ||
				!equalResourcePtr(input.SourceResource, state.SourceResource) ||
				!equalResourcePtr(input.DestinationResource, state.DestinationResource) {
				equal = false

				break
			}
		}

		if !equal {
			diff["rules"] = p.PropertyDiff{
				InputDiff: false,
				Kind:      p.Update,
			}
		}
	}

	if !equalSlicePtr(req.Inputs.SourcePostureChecks, req.State.SourcePostureChecks) {
		diff["postureChecks"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	p.GetLogger(ctx).Debugf("Diff:Policy[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*Policy) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[PolicyArgs], error) { //nolint:gocognit,gocyclo
	p.GetLogger(ctx).Debugf("Check:Policy old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[PolicyArgs](ctx, req.NewInputs)
	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{
			Property: "name",
			Reason:   "name must not be empty",
		})
	}

	if len(args.Rules) == 0 {
		failures = append(failures, p.CheckFailure{
			Property: "rules",
			Reason:   "at least one rule is required",
		})
	}

	for ruleIndex, rule := range args.Rules {
		if isBlank(rule.Name) {
			failures = append(failures, p.CheckFailure{
				Property: fmt.Sprintf("rules[%d].name", ruleIndex),
				Reason:   "rule name must not be empty",
			})
		}

		if rule.Ports != nil {
			for portIndex, port := range *rule.Ports {
				if isBlank(port) {
					failures = append(failures, p.CheckFailure{
						Property: fmt.Sprintf("rules[%d].ports[%d]", ruleIndex, portIndex),
						Reason:   "port must not be empty",
					})
				}
			}
		}

		if rule.PortRanges != nil {
			for portRangeIndex, portRange := range *rule.PortRanges {
				if portRange.Start < 1 || portRange.Start > 65535 {
					failures = append(failures, p.CheckFailure{
						Property: fmt.Sprintf("rules[%d].portRanges[%d].start", ruleIndex, portRangeIndex),
						Reason:   "start must be between 1 and 65535",
					})
				}

				if portRange.End < 1 || portRange.End > 65535 {
					failures = append(failures, p.CheckFailure{
						Property: fmt.Sprintf("rules[%d].portRanges[%d].end", ruleIndex, portRangeIndex),
						Reason:   "end must be between 1 and 65535",
					})
				}

				if portRange.Start > portRange.End {
					failures = append(failures, p.CheckFailure{
						Property: fmt.Sprintf("rules[%d].portRanges[%d]", ruleIndex, portRangeIndex),
						Reason:   "start must be less than or equal to end",
					})
				}
			}
		}

		hasSources := rule.Sources != nil && len(*rule.Sources) > 0

		hasSourceResource := rule.SourceResource != nil
		if !hasSources && !hasSourceResource {
			failures = append(failures, p.CheckFailure{
				Property: fmt.Sprintf("rules[%d].sources", ruleIndex),
				Reason:   "at least one source or sourceResource is required",
			})
		}

		hasDestinations := rule.Destinations != nil && len(*rule.Destinations) > 0

		hasDestinationResource := rule.DestinationResource != nil
		if !hasDestinations && !hasDestinationResource {
			failures = append(failures, p.CheckFailure{
				Property: fmt.Sprintf("rules[%d].destinations", ruleIndex),
				Reason:   "at least one destination or destinationResource is required",
			})
		}

		if rule.Sources != nil {
			for sourceIndex, source := range *rule.Sources {
				if isBlank(source) {
					failures = append(failures, p.CheckFailure{
						Property: fmt.Sprintf("rules[%d].sources[%d]", ruleIndex, sourceIndex),
						Reason:   "source id must not be empty",
					})
				}
			}
		}

		if rule.Destinations != nil {
			for destinationIndex, destination := range *rule.Destinations {
				if isBlank(destination) {
					failures = append(failures, p.CheckFailure{
						Property: fmt.Sprintf("rules[%d].destinations[%d]", ruleIndex, destinationIndex),
						Reason:   "destination id must not be empty",
					})
				}
			}
		}

		if rule.SourceResource != nil {
			if isBlank(rule.SourceResource.ID) {
				failures = append(failures, p.CheckFailure{
					Property: fmt.Sprintf("rules[%d].sourceResource.id", ruleIndex),
					Reason:   "sourceResource.id must not be empty",
				})
			}

			if isBlank(string(rule.SourceResource.Type)) {
				failures = append(failures, p.CheckFailure{
					Property: fmt.Sprintf("rules[%d].sourceResource.type", ruleIndex),
					Reason:   "sourceResource.type must not be empty",
				})
			}
		}

		if rule.DestinationResource != nil {
			if isBlank(rule.DestinationResource.ID) {
				failures = append(failures, p.CheckFailure{
					Property: fmt.Sprintf("rules[%d].destinationResource.id", ruleIndex),
					Reason:   "destinationResource.id must not be empty",
				})
			}

			if isBlank(string(rule.DestinationResource.Type)) {
				failures = append(failures, p.CheckFailure{
					Property: fmt.Sprintf("rules[%d].destinationResource.type", ruleIndex),
					Reason:   "destinationResource.type must not be empty",
				})
			}
		}
	}

	return infer.CheckResponse[PolicyArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*Policy) WireDependencies(f infer.FieldSelector, args *PolicyArgs, state *PolicyState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.Description).DependsOn(f.InputField(&args.Description))
	f.OutputField(&state.Enabled).DependsOn(f.InputField(&args.Enabled))
	f.OutputField(&state.Rules).DependsOn(f.InputField(&args.Rules))
	f.OutputField(&state.SourcePostureChecks).DependsOn(f.InputField(&args.SourcePostureChecks))
}

// Converts a slice of RulePortRange from state model to API model.
func toAPIPortRanges(rulePortRange *[]RulePortRange) *[]nbapi.RulePortRange {
	if rulePortRange == nil {
		return nil
	}

	out := make([]nbapi.RulePortRange, len(*rulePortRange))
	for rulePRIndex, rulePR := range *rulePortRange {
		out[rulePRIndex] = nbapi.RulePortRange{Start: rulePR.Start, End: rulePR.End}
	}

	return &out
}

// Converts a slice of API RulePortRange to state model.
func fromAPIPortRanges(reulePortRangeAPI *[]nbapi.RulePortRange) *[]RulePortRange {
	if reulePortRangeAPI == nil {
		return nil
	}

	out := make([]RulePortRange, len(*reulePortRangeAPI))
	for rulePRIndex, rulePR := range *reulePortRangeAPI {
		out[rulePRIndex] = RulePortRange{Start: rulePR.Start, End: rulePR.End}
	}

	return &out
}

// Converts a slice of nbapi.GroupMinimum to state RuleGroup, sorted by ID.
func fromAPIGroupMinimums(group *[]nbapi.GroupMinimum) *[]RuleGroup {
	if group == nil {
		return nil
	}

	out := make([]RuleGroup, len(*group))
	for groupIndex, group := range *group {
		out[groupIndex] = RuleGroup{ID: group.Id, Name: group.Name}
	}

	slices.SortFunc(out, func(a, b RuleGroup) int {
		if a.ID < b.ID {
			return -1
		}

		if a.ID > b.ID {
			return 1
		}

		return 0
	})

	return &out
}

func toGroupIDs(groups *[]RuleGroup) *[]string {
	if groups == nil {
		return nil
	}

	iDs := make([]string, len(*groups))
	for i, g := range *groups {
		iDs[i] = g.ID
	}

	return &iDs
}

func policyRulesToArgs(rules []nbapi.PolicyRule) []PolicyRuleArgs {
	out := make([]PolicyRuleArgs, len(rules))
	for i, rule := range rules {
		out[i] = PolicyRuleArgs{
			ID:                  rule.Id,
			Name:                rule.Name,
			Description:         nil,
			Bidirectional:       rule.Bidirectional,
			Action:              RuleAction(rule.Action),
			Enabled:             rule.Enabled,
			Protocol:            Protocol(rule.Protocol),
			Ports:               rule.Ports,
			PortRanges:          fromAPIPortRanges(rule.PortRanges),
			Sources:             groupMinimumIDs(rule.Sources),
			Destinations:        groupMinimumIDs(rule.Destinations),
			SourceResource:      fromAPIResource(rule.SourceResource),
			DestinationResource: fromAPIResource(rule.DestinationResource),
		}
	}

	return out
}

func groupMinimumIDs(groups *[]nbapi.GroupMinimum) *[]string {
	if groups == nil {
		return nil
	}

	out := make([]string, len(*groups))
	for i, group := range *groups {
		out[i] = group.Id
	}

	slices.Sort(out)

	return &out
}

func equalPortRangePtr(portRangeA, portRangeB *[]RulePortRange) bool {
	aLen := 0
	if portRangeA != nil {
		aLen = len(*portRangeA)
	}

	bLen := 0
	if portRangeB != nil {
		bLen = len(*portRangeB)
	}

	if aLen == 0 && bLen == 0 {
		return true
	}

	if portRangeA == nil || portRangeB == nil || aLen != bLen {
		return false
	}

	for i := range *portRangeA {
		if (*portRangeA)[i] != (*portRangeB)[i] {
			return false
		}
	}

	return true
}
