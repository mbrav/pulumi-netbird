package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Policy defines the Pulumi resource handler for NetBird policy resources.
type Policy struct{}

// PolicyArgs defines the user-supplied arguments for creating/updating a Policy resource.
type PolicyArgs struct {
	Name                string           `pulumi:"name"`                    // Policy name (required)
	Description         *string          `pulumi:"description,optional"`    // Optional description
	Enabled             bool             `pulumi:"enabled"`                 // Whether the policy is enabled
	Rules               []PolicyRuleArgs `pulumi:"rules"`                   // List of rules defined in the policy
	SourcePostureChecks *[]string        `pulumi:"posture_checks,optional"` // Optional list of posture check IDs
}

// PolicyState represents the state of a Policy resource stored in Pulumi state.
type PolicyState struct {
	Name                string            `pulumi:"name"`
	Description         *string           `pulumi:"description,optional"`
	Enabled             bool              `pulumi:"enabled"`
	Rules               []PolicyRuleState `pulumi:"rules"`
	SourcePostureChecks *[]string         `pulumi:"posture_checks,optional"`
}

// PolicyRuleArgs represents user input for an individual rule in a policy.
type PolicyRuleArgs struct {
	Id                  *string          `pulumi:"id,optional"`           // Optional rule ID (used for updates)
	Name                string           `pulumi:"name"`                  // Rule name
	Description         *string          `pulumi:"description,optional"`  // Optional rule description
	Bidirectional       bool             `pulumi:"bidirectional"`         // Whether the rule is bidirectional
	Action              RuleAction       `pulumi:"action"`                // Rule action (accept/drop)
	Enabled             bool             `pulumi:"enabled"`               // Whether the rule is enabled
	Protocol            Protocol         `pulumi:"protocol"`              // Network protocol
	Ports               *[]string        `pulumi:"ports,optional"`        // Optional list of specific ports
	PortRanges          *[]RulePortRange `pulumi:"portRanges,optional"`   // Optional list of port ranges
	Sources             *[]string        `pulumi:"sources,optional"`      // Optional list of source group IDs
	Destinations        *[]string        `pulumi:"destinations,optional"` // Optional list of destination group IDs
	SourceResource      *Resource        `pulumi:"source,optional"`       // Optional single source resource
	DestinationResource *Resource        `pulumi:"destination,optional"`  // Optional single destination resource
}

// PolicyRuleState represents the state of an individual rule within a policy.
type PolicyRuleState struct {
	Id                  *string          `pulumi:"id,optional"`
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
	SourceResource      *Resource        `pulumi:"source,optional"`
	DestinationResource *Resource        `pulumi:"destination,optional"`
}

// RuleAction defines the allowed actions for a rule (accept/drop).
// This wraps the nbapi type to allow method definitions (like Values()).
type RuleAction string

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
	ProtocolTcp  Protocol = Protocol(nbapi.PolicyRuleProtocolTcp)
	ProtocolUdp  Protocol = Protocol(nbapi.PolicyRuleProtocolUdp)
)

// Values returns valid protocol values for Pulumi enum support.
func (Protocol) Values() []infer.EnumValue[Protocol] {
	return []infer.EnumValue[Protocol]{
		{Name: "All", Value: ProtocolAll, Description: "All protocols"},
		{Name: "ICMP", Value: ProtocolIcmp, Description: "ICMP protocol"},
		{Name: "TCP", Value: ProtocolTcp, Description: "TCP protocol"},
		{Name: "UDP", Value: ProtocolUdp, Description: "UDP protocol"},
	}
}

// Resource represents a single NetBird resource used in a rule (e.g., domain, host, subnet).
type Resource struct {
	Id   string       `pulumi:"id"`   // The unique ID of the resource
	Type ResourceType `pulumi:"type"` // The type of the resource (domain, host, subnet)
}

func (r *Resource) Annotate(a infer.Annotator) {
	a.Describe(&r.Id, "The unique identifier of the resource.")
	a.Describe(&r.Type, "The type of resource: 'domain', 'host', or 'subnet'.")
}

// ResourceType defines the allowed resource types for a policy rule.
type ResourceType string

// Enum constants for resource types.
const (
	ResourceTypeDomain ResourceType = ResourceType(nbapi.ResourceTypeDomain)
	ResourceTypeHost   ResourceType = ResourceType(nbapi.ResourceTypeHost)
	ResourceTypeSubnet ResourceType = ResourceType(nbapi.ResourceTypeSubnet)
)

// Values returns the list of supported ResourceType values for Pulumi enum generation.
func (ResourceType) Values() []infer.EnumValue[ResourceType] {
	return []infer.EnumValue[ResourceType]{
		{Name: "Domain", Value: ResourceTypeDomain, Description: "A domain resource (e.g., example.com)."},
		{Name: "Host", Value: ResourceTypeHost, Description: "A host resource (e.g., peer or device)."},
		{Name: "Subnet", Value: ResourceTypeSubnet, Description: "A subnet resource (e.g., 192.168.0.0/24)."},
	}
}

type RulePortRange struct {
	Start int `pulumi:"start"`
	End   int `pulumi:"end"`
}

type RuleGroup struct {
	Id   string `pulumi:"id"`
	Name string `pulumi:"name"`
}

func (Policy) Annotate(a infer.Annotator) {
	a.Describe(&Policy{}, "A NetBird policy defining rules for communication between peers.")
}

func (p *PolicyArgs) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the policy.")
	a.Describe(&p.Description, "An optional description of the policy.")
	a.Describe(&p.Enabled, "Whether the policy is currently active.")
	a.Describe(&p.Rules, "The list of rules defining the behavior of this policy.")
	a.Describe(&p.SourcePostureChecks, "Optional posture check IDs used as sources in policy rules.")
}

func (p *PolicyState) Annotate(a infer.Annotator) {
	a.Describe(&p.Name, "The name of the policy.")
	a.Describe(&p.Description, "An optional description of the policy.")
	a.Describe(&p.Enabled, "Whether the policy is currently active.")
	a.Describe(&p.Rules, "The list of rules defining the behavior of this policy.")
	a.Describe(&p.SourcePostureChecks, "Optional posture check IDs used as sources in policy rules.")
}

func (p *PolicyRuleArgs) Annotate(a infer.Annotator) {
	a.Describe(&p.Id, "Optional unique identifier for the policy rule.")
	a.Describe(&p.Name, "The name of the policy rule.")
	a.Describe(&p.Description, "An optional description of the policy rule.")
	a.Describe(&p.Bidirectional, "Whether the rule applies bidirectionally.")
	a.Describe(&p.Action, "The action to take: 'accept' or 'drop'.")
	a.Describe(&p.Enabled, "Whether the rule is active.")
	a.Describe(&p.Protocol, "The protocol: 'tcp', 'udp', 'icmp', or 'all'.")
	a.Describe(&p.Ports, "Optional list of ports.")
	a.Describe(&p.PortRanges, "Optional list of port ranges.")
	a.Describe(&p.Sources, "Optional list of source group IDs.")
	a.Describe(&p.Destinations, "Optional list of destination group IDs.")
	a.Describe(&p.SourceResource, "Optional source resource for the rule.")
	a.Describe(&p.DestinationResource, "Optional destination resource for the rule.")
}

func (p *PolicyRuleState) Annotate(a infer.Annotator) {
	a.Describe(&p.Id, "Optional unique identifier for the policy rule.")
	a.Describe(&p.Name, "The name of the policy rule.")
	a.Describe(&p.Description, "An optional description of the policy rule.")
	a.Describe(&p.Bidirectional, "Whether the rule applies bidirectionally.")
	a.Describe(&p.Action, "The action to take: 'accept' or 'drop'.")
	a.Describe(&p.Enabled, "Whether the rule is active.")
	a.Describe(&p.Protocol, "The protocol: 'tcp', 'udp', 'icmp', or 'all'.")
	a.Describe(&p.Ports, "Optional list of ports.")
	a.Describe(&p.PortRanges, "Optional list of port ranges.")
	a.Describe(&p.Sources, "Optional list of source groups.")
	a.Describe(&p.Destinations, "Optional list of destination groups.")
	a.Describe(&p.SourceResource, "Optional source resource for the rule.")
	a.Describe(&p.DestinationResource, "Optional destination resource for the rule.")
}

func (g *RuleGroup) Annotate(a infer.Annotator) {
	a.Describe(&g.Id, "The unique identifier of the group.")
	a.Describe(&g.Name, "The name of the group.")
}

// Create creates a new NetBird policy.
func (*Policy) Create(ctx context.Context, req infer.CreateRequest[PolicyArgs]) (infer.CreateResponse[PolicyState], error) {
	p.GetLogger(ctx).Debugf("Create:Policy")
	if req.DryRun {
		// Convert PolicyRuleArgs to PolicyRuleState for preview
		rules := make([]PolicyRuleState, len(req.Inputs.Rules))
		for i, rule := range req.Inputs.Rules {
			// Construct sources and destination groups
			var sources, destinations *[]RuleGroup
			if rule.Sources != nil {
				groups := make([]RuleGroup, len(*rule.Sources))
				for j, g := range *rule.Sources {
					groups[j] = RuleGroup{Name: "preview", Id: g}
				}
				sources = &groups
			}
			if rule.Destinations != nil {
				groups := make([]RuleGroup, len(*rule.Destinations))
				for j, g := range *rule.Sources {
					groups[j] = RuleGroup{Name: "preview", Id: g}
				}
				destinations = &groups
			}

			rules[i] = PolicyRuleState{
				Id:                  rule.Id,
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

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[PolicyState]{}, err
	}

	// Convert input rules to nbapi.PolicyRuleUpdate
	apiRules := make([]nbapi.PolicyRuleUpdate, len(req.Inputs.Rules))
	for i, rule := range req.Inputs.Rules {
		apiRules[i] = nbapi.PolicyRuleUpdate{
			Id:                  rule.Id,
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

	// Convert created rules to PolicyRuleState
	rules := make([]PolicyRuleState, len(created.Rules))
	for i, rule := range created.Rules {
		rules[i] = PolicyRuleState{
			Id:                  rule.Id,
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

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[PolicyArgs, PolicyState]{}, err
	}

	policy, err := client.Policies.Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[PolicyArgs, PolicyState]{}, fmt.Errorf("reading policy failed: %w", err)
	}

	rules := make([]PolicyRuleState, len(policy.Rules))
	for i, rule := range policy.Rules {
		rules[i] = PolicyRuleState{
			Id:                  rule.Id,
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

	return infer.ReadResponse[PolicyArgs, PolicyState]{
		ID: req.ID,
		Inputs: PolicyArgs{
			Name:                policy.Name,
			Description:         policy.Description,
			Enabled:             policy.Enabled,
			Rules:               req.Inputs.Rules, // Not strictly accurate; optional improvement: reconstruct from `policy.Rules`
			SourcePostureChecks: &policy.SourcePostureChecks,
		},
		State: PolicyState{
			Name:                policy.Name,
			Description:         policy.Description,
			Enabled:             policy.Enabled,
			Rules:               rules,
			SourcePostureChecks: &policy.SourcePostureChecks,
		},
	}, nil
}

// Update updates an existing NetBird policy.
func (*Policy) Update(ctx context.Context, req infer.UpdateRequest[PolicyArgs, PolicyState]) (infer.UpdateResponse[PolicyState], error) {
	p.GetLogger(ctx).Debugf("Update:Policy[%s]", req.ID)

	if req.DryRun {
		// Construct PolicyRuleState for preview output
		rules := make([]PolicyRuleState, len(req.Inputs.Rules))
		for i, rule := range req.Inputs.Rules {
			var sources, destinations *[]RuleGroup
			if rule.Sources != nil {
				groups := make([]RuleGroup, len(*rule.Sources))
				for j, g := range *rule.Sources {
					groups[j] = RuleGroup{Name: "preview", Id: g}
				}
				sources = &groups
			}
			if rule.Destinations != nil {
				groups := make([]RuleGroup, len(*rule.Destinations))
				for j, g := range *rule.Destinations {
					groups[j] = RuleGroup{Name: "preview", Id: g}
				}
				destinations = &groups
			}

			rules[i] = PolicyRuleState{
				Id:                  rule.Id,
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

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[PolicyState]{}, err
	}

	// Convert input rules to nbapi.PolicyRuleUpdate
	apiRules := make([]nbapi.PolicyRuleUpdate, len(req.Inputs.Rules))
	for i, rule := range req.Inputs.Rules {
		apiRules[i] = nbapi.PolicyRuleUpdate{
			Id:                  rule.Id,
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
		return infer.UpdateResponse[PolicyState]{}, fmt.Errorf("updating policy failed: %w", err)
	}

	// Convert updated rules to PolicyRuleState
	rules := make([]PolicyRuleState, len(updated.Rules))
	for i, rule := range updated.Rules {
		rules[i] = PolicyRuleState{
			Id:                  rule.Id,
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

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	err = client.Policies.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting policy failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*Policy) Diff(ctx context.Context, req infer.DiffRequest[PolicyArgs, PolicyState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Policy[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{Kind: p.Update}
	}
	if strPtr(req.Inputs.Description) != strPtr(req.State.Description) {
		diff["description"] = p.PropertyDiff{Kind: p.Update}
	}
	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{Kind: p.Update}
	}
	// Rules Diff
	if len(req.Inputs.Rules) != len(req.State.Rules) {
		diff["rules"] = p.PropertyDiff{Kind: p.Update}
	} else {
		equal := true
		for i := range req.Inputs.Rules {
			in := req.Inputs.Rules[i]
			st := req.State.Rules[i]

			p.GetLogger(ctx).Debugf("Diff:Policy[%s]:Rules[%d] a=%+v b=%+v", req.ID, i, in, st)

			if in.Name != st.Name ||
				!equalPtr(in.Description, st.Description) ||
				in.Bidirectional != st.Bidirectional ||
				in.Action != RuleAction(st.Action) ||
				in.Enabled != st.Enabled ||
				in.Protocol != Protocol(st.Protocol) ||
				!equalSlicePtr(in.Ports, st.Ports) ||
				!equalPortRangePtr(in.PortRanges, st.PortRanges) ||
				// !equalSlicePtr(in.Sources, toGroupIds(st.Sources)) ||
				// !equalSlicePtr(in.Destinations, toGroupIds(st.Destinations)) ||
				!equalResourcePtr(in.SourceResource, st.SourceResource) ||
				!equalResourcePtr(in.DestinationResource, st.DestinationResource) {
				equal = false
				break
			}
		}
		if !equal {
			diff["rules"] = p.PropertyDiff{Kind: p.Update}
		}
	}
	if !equalSlicePtr(req.Inputs.SourcePostureChecks, req.State.SourcePostureChecks) {
		diff["posture_checks"] = p.PropertyDiff{Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:Policy[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Converts a slice of RulePortRange from state model to API model
func toAPIPortRanges(in *[]RulePortRange) *[]nbapi.RulePortRange {
	if in == nil {
		return nil
	}
	out := make([]nbapi.RulePortRange, len(*in))
	for i, r := range *in {
		out[i] = nbapi.RulePortRange{Start: r.Start, End: r.End}
	}
	return &out
}

// Converts a slice of API RulePortRange to state model
func fromAPIPortRanges(in *[]nbapi.RulePortRange) *[]RulePortRange {
	if in == nil {
		return nil
	}
	out := make([]RulePortRange, len(*in))
	for i, r := range *in {
		out[i] = RulePortRange{Start: r.Start, End: r.End}
	}
	return &out
}

// Converts a slice of nbapi.GroupMinimum to state RuleGroup
func fromAPIGroupMinimums(in *[]nbapi.GroupMinimum) *[]RuleGroup {
	if in == nil {
		return nil
	}
	out := make([]RuleGroup, len(*in))
	for i, r := range *in {
		out[i] = RuleGroup{Id: r.Id, Name: r.Name}
	}
	return &out
}

// Converts a single Resource to nbapi.Resource
func toAPIResource(in *Resource) *nbapi.Resource {
	if in == nil {
		return nil
	}
	return &nbapi.Resource{
		Id:   in.Id,
		Type: nbapi.ResourceType(in.Type),
	}
}

// Converts a single nbapi.Resource to Resource
func fromAPIResource(in *nbapi.Resource) *Resource {
	if in == nil {
		return nil
	}
	return &Resource{
		Id:   in.Id,
		Type: ResourceType(in.Type),
	}
}

func toGroupIds(groups *[]RuleGroup) *[]string {
	if groups == nil {
		return nil
	}
	ids := make([]string, len(*groups))
	for i, g := range *groups {
		ids[i] = g.Id
	}
	return &ids
}

func equalPortRangePtr(a, b *[]RulePortRange) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil || len(*a) != len(*b) {
		return false
	}
	for i := range *a {
		if (*a)[i] != (*b)[i] {
			return false
		}
	}
	return true
}

func equalResourcePtr(a, b *Resource) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Type == b.Type && a.Id == b.Id
}
