package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
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
	Id                  *string                `pulumi:"id,optional"`                  // Optional rule ID (used for updates)
	Name                string                 `pulumi:"name"`                         // Rule name
	Description         *string                `pulumi:"description,optional"`         // Optional rule description
	Bidirectional       bool                   `pulumi:"bidirectional"`                // Whether the rule is bidirectional
	Action              RuleAction             `pulumi:"action"`                       // Rule action (accept/drop)
	Enabled             bool                   `pulumi:"enabled"`                      // Whether the rule is enabled
	Protocol            Protocol               `pulumi:"protocol"`                     // Network protocol
	Ports               *[]string              `pulumi:"ports,optional"`               // Optional list of specific ports
	PortRanges          *[]PolicyRulePortRange `pulumi:"portRanges,optional"`          // Optional list of port ranges
	Sources             *[]string              `pulumi:"sources,optional"`             // Optional list of source group IDs
	Destinations        *[]string              `pulumi:"destinations,optional"`        // Optional list of destination group IDs
	SourceResource      *Resource              `pulumi:"sourceResource,optional"`      // Optional single source resource
	DestinationResource *Resource              `pulumi:"destinationResource,optional"` // Optional single destination resource
}

// PolicyRuleState represents the state of an individual rule within a policy.
type PolicyRuleState struct {
	Id                  *string                `pulumi:"id,optional"`
	Name                string                 `pulumi:"name"`
	Description         *string                `pulumi:"description,optional"`
	Bidirectional       bool                   `pulumi:"bidirectional"`
	Action              RuleAction             `pulumi:"action"`
	Enabled             bool                   `pulumi:"enabled"`
	Protocol            Protocol               `pulumi:"protocol"`
	Ports               *[]string              `pulumi:"ports,optional"`
	PortRanges          *[]PolicyRulePortRange `pulumi:"portRanges,optional"`
	Sources             *[]PolicyRuleGroup     `pulumi:"sources,optional"` // Fully-resolved group info (not just IDs)
	Destinations        *[]PolicyRuleGroup     `pulumi:"destinations,optional"`
	SourceResource      *Resource              `pulumi:"sourceResource,optional"`
	DestinationResource *Resource              `pulumi:"destinationResource,optional"`
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

type PolicyRulePortRange struct {
	Start int `pulumi:"start"`
	End   int `pulumi:"end"`
}

type PolicyRuleGroup struct {
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

func (g *PolicyRuleGroup) Annotate(a infer.Annotator) {
	a.Describe(&g.Id, "The unique identifier of the group.")
	a.Describe(&g.Name, "The name of the group.")
}

// Create creates a new NetBird policy.
func (p *Policy) Create(ctx context.Context, req infer.CreateRequest[PolicyArgs]) (infer.CreateResponse[PolicyState], error) {
	if req.DryRun {
		// Convert PolicyRuleArgs to PolicyRuleState for preview
		rules := make([]PolicyRuleState, len(req.Inputs.Rules))
		for i, rule := range req.Inputs.Rules {
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
				Sources:             nil, // Cannot populate GroupMinimums in dry run
				Destinations:        nil,
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
			PortRanges:          convertRulePortRangesToAPI(rule.PortRanges),
			Sources:             rule.Sources,
			Destinations:        rule.Destinations,
			SourceResource:      convertResourceToAPI(rule.SourceResource),
			DestinationResource: convertResourceToAPI(rule.DestinationResource),
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
			Id:            rule.Id,
			Name:          rule.Name,
			Description:   rule.Description,
			Bidirectional: rule.Bidirectional,
			Action:        RuleAction(rule.Action),
			Enabled:       rule.Enabled,
			Protocol:      Protocol(rule.Protocol),
			Ports:         rule.Ports,
			// PortRanges:          convertRulePortRangesToAPI(rule.PortRanges),
			// Sources:             rule.Sources,
			// Destinations:        rule.Destinations,
			// SourceResource:      rule.SourceResource,
			// DestinationResource: rule.DestinationResource,
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

func convertRulePortRangesToAPI(in *[]PolicyRulePortRange) *[]nbapi.RulePortRange {
	if in == nil {
		return nil
	}
	out := make([]nbapi.RulePortRange, len(*in))
	for i, r := range *in {
		out[i] = nbapi.RulePortRange{
			End:   r.End,
			Start: r.Start,
		}
	}
	return &out
}

func convertResourceToAPI(in *Resource) *nbapi.Resource {
	if in == nil {
		return nil
	}
	return &nbapi.Resource{
		Id:   in.Id,
		Type: nbapi.ResourceType(in.Type), // Cast to API enum
	}
}
