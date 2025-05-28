package resource

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// TEST: InputDiff: false

// DNS represents a DNS Group.
type DNS struct{}

// Annotate adds a description to the DNS resource type.
func (n *DNS) Annotate(a infer.Annotator) {
	a.Describe(&n, "A NetBird network.")
}

// DNSArgs defines input fields for creating or updating a network.
type DNSArgs struct {
	Name                 string       `pulumi:"name"`
	Description          string       `pulumi:"description"`
	Domains              []string     `pulumi:"domains"`
	Enabled              bool         `pulumi:"enabled"`
	Groups               []string     `pulumi:"groups"`
	Primary              bool         `pulumi:"primary"`
	Nameservers          []Nameserver `pulumi:"nameservers"`
	SearchDomainsEnabled bool         `pulumi:"search_domains_enabled"`
}

// Annotate provides documentation for DNSArgs fields.
func (n *DNSArgs) Annotate(a infer.Annotator) {
	a.Describe(&n.Name, "Name of nameserver group name")
	a.Describe(&n.Description, "Description of the nameserver group")
	a.Describe(&n.Domains, "Domains Match domain list. It should be empty only if primary is true.")
	a.Describe(&n.Enabled, "Enabled Nameserver group status")
	a.Describe(&n.Groups, "Groups Distribution group IDs that defines group of peers that will use this nameserver group")
	a.Describe(&n.Primary, "Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.")
	a.Describe(&n.Nameservers, "Nameservers Nameserver list")
	a.Describe(&n.SearchDomainsEnabled, "SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.")
}

// DNSState represents the output state of a network resource.
type DNSState struct {
	Name                 string       `pulumi:"name"`
	Description          string       `pulumi:"description"`
	Domains              []string     `pulumi:"domains"`
	Enabled              bool         `pulumi:"enabled"`
	Groups               []string     `pulumi:"groups"`
	Primary              bool         `pulumi:"primary"`
	Nameservers          []Nameserver `pulumi:"nameservers"`
	SearchDomainsEnabled bool         `pulumi:"search_domains_enabled"`
}

// Annotate provides documentation for DNSState fields.
func (n *DNSState) Annotate(a infer.Annotator) {
	a.Describe(&n.Name, "Name of nameserver group name")
	a.Describe(&n.Description, "Description of the nameserver group")
	a.Describe(&n.Domains, "Domains Match domain list. It should be empty only if primary is true.")
	a.Describe(&n.Enabled, "Enabled Nameserver group status")
	a.Describe(&n.Groups, "Groups Distribution group IDs that defines group of peers that will use this nameserver group")
	a.Describe(&n.Primary, "Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.")
	a.Describe(&n.Nameservers, "Nameservers Nameserver list")
	a.Describe(&n.SearchDomainsEnabled, "SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.")
}

// Nameserver defines model for Nameserver.
type Nameserver struct {
	Ip     string           `pulumi:"ip"`
	NsType NameserverNsType `pulumi:"type"`
	Port   int              `pulumi:"port"`
}

// Annotate provides documentation for DNSState fields.
func (n *Nameserver) Annotate(a infer.Annotator) {
	a.Describe(&n.Ip, "Ip Nameserver IP")
	a.Describe(&n.NsType, "NsType Nameserver Type")
	a.Describe(&n.Port, "Port Nameserver Port")
}

// NameserverNsType defines the allowed DNS types
// This wraps the nbapi type to allow method definitions (like Values()).
type NameserverNsType string

const (
	NameserverNsTypeUdp NameserverNsType = NameserverNsType(nbapi.NameserverNsTypeUdp)
)

// Values returns the valid enum values for NameserverNsType, used by Pulumi for schema generation and validation.
func (NameserverNsType) Values() []infer.EnumValue[NameserverNsType] {
	return []infer.EnumValue[NameserverNsType]{
		{Name: "udp", Value: NameserverNsTypeUdp, Description: "UDP type"},
	}
}

// Create creates a new NetBird DNS (Nameserver Group).
func (*DNS) Create(ctx context.Context, req infer.CreateRequest[DNSArgs]) (infer.CreateResponse[DNSState], error) {
	p.GetLogger(ctx).Debugf("Create:DNS name=%s, description=%s", req.Inputs.Name, req.Inputs.Description)

	if req.DryRun {
		return infer.CreateResponse[DNSState]{
			ID: "preview",
			Output: DNSState{
				Name:                 req.Inputs.Name,
				Description:          req.Inputs.Description,
				Domains:              req.Inputs.Domains,
				Enabled:              req.Inputs.Enabled,
				Groups:               req.Inputs.Groups,
				Primary:              req.Inputs.Primary,
				Nameservers:          req.Inputs.Nameservers,
				SearchDomainsEnabled: req.Inputs.SearchDomainsEnabled,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[DNSState]{}, fmt.Errorf("error getting Netbird client: %w", err)
	}

	// Build request payload
	createReq := nbapi.NameserverGroupRequest{
		Name:                 req.Inputs.Name,
		Description:          req.Inputs.Description,
		Domains:              req.Inputs.Domains,
		Enabled:              req.Inputs.Enabled,
		Groups:               req.Inputs.Groups,
		Primary:              req.Inputs.Primary,
		Nameservers:          toAPINameservers(req.Inputs.Nameservers),
		SearchDomainsEnabled: req.Inputs.SearchDomainsEnabled,
	}

	// Call the API
	created, err := client.DNS.CreateNameserverGroup(ctx, createReq)
	if err != nil {
		return infer.CreateResponse[DNSState]{}, fmt.Errorf("creating DNS group failed: %w", err)
	}

	return infer.CreateResponse[DNSState]{
		ID: created.Id,
		Output: DNSState{
			Name:                 created.Name,
			Description:          created.Description,
			Domains:              created.Domains,
			Enabled:              created.Enabled,
			Groups:               created.Groups,
			Primary:              created.Primary,
			Nameservers:          fromAPINameservers(created.Nameservers),
			SearchDomainsEnabled: created.SearchDomainsEnabled,
		},
	}, nil
}

// Read reads a DNS (Nameserver Group) from NetBird.
func (*DNS) Read(ctx context.Context, req infer.ReadRequest[DNSArgs, DNSState]) (infer.ReadResponse[DNSArgs, DNSState], error) {
	p.GetLogger(ctx).Debugf("Read:DNS[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[DNSArgs, DNSState]{}, fmt.Errorf("error getting Netbird client: %w", err)
	}

	group, err := client.DNS.GetNameserverGroup(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[DNSArgs, DNSState]{}, fmt.Errorf("reading DNS group failed: %w", err)
	}

	// Convert API nameservers to Pulumi state format
	stateNameservers := make([]Nameserver, len(group.Nameservers))
	for i, ns := range group.Nameservers {
		stateNameservers[i] = Nameserver{
			Ip:     ns.Ip,
			NsType: NameserverNsType(ns.NsType),
			Port:   ns.Port,
		}
	}

	// Return response with both current Inputs and updated State
	return infer.ReadResponse[DNSArgs, DNSState]{
		ID: req.ID,
		Inputs: DNSArgs{
			Name:                 group.Name,
			Description:          group.Description,
			Domains:              group.Domains,
			Enabled:              group.Enabled,
			Groups:               group.Groups,
			Primary:              group.Primary,
			Nameservers:          stateNameservers,
			SearchDomainsEnabled: group.SearchDomainsEnabled,
		},
		State: DNSState{
			Name:                 group.Name,
			Description:          group.Description,
			Domains:              group.Domains,
			Enabled:              group.Enabled,
			Groups:               group.Groups,
			Primary:              group.Primary,
			Nameservers:          stateNameservers,
			SearchDomainsEnabled: group.SearchDomainsEnabled,
		},
	}, nil
}

// Update updates a DNS (Nameserver Group) from NetBird.
func (*DNS) Update(ctx context.Context, req infer.UpdateRequest[DNSArgs, DNSState]) (infer.UpdateResponse[DNSState], error) {
	p.GetLogger(ctx).Debugf("Update:DNS[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[DNSState]{
			Output: DNSState{
				Name:                 req.Inputs.Name,
				Description:          req.Inputs.Description,
				Domains:              req.Inputs.Domains,
				Enabled:              req.Inputs.Enabled,
				Groups:               req.Inputs.Groups,
				Primary:              req.Inputs.Primary,
				Nameservers:          req.Inputs.Nameservers,
				SearchDomainsEnabled: req.Inputs.SearchDomainsEnabled,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[DNSState]{}, fmt.Errorf("error getting Netbird client: %w", err)
	}

	updated, err := client.DNS.UpdateNameserverGroup(ctx, req.ID, nbapi.NameserverGroupRequest{
		Name:                 req.Inputs.Name,
		Description:          req.Inputs.Description,
		Domains:              req.Inputs.Domains,
		Enabled:              req.Inputs.Enabled,
		Groups:               req.Inputs.Groups,
		Primary:              req.Inputs.Primary,
		Nameservers:          toAPINameservers(req.Inputs.Nameservers),
		SearchDomainsEnabled: req.Inputs.SearchDomainsEnabled,
	})
	if err != nil {
		return infer.UpdateResponse[DNSState]{}, fmt.Errorf("updating DNS entry failed: %w", err)
	}

	return infer.UpdateResponse[DNSState]{
		Output: DNSState{
			Name:                 updated.Name,
			Description:          updated.Description,
			Enabled:              updated.Enabled,
			Domains:              updated.Domains,
			Groups:               updated.Groups,
			Primary:              updated.Primary,
			Nameservers:          fromAPINameservers(updated.Nameservers),
			SearchDomainsEnabled: updated.SearchDomainsEnabled,
		},
	}, nil
}

// Converts a slice of internal Nameserver to API Nameserver.
func toAPINameservers(in []Nameserver) []nbapi.Nameserver {
	apiNameservers := make([]nbapi.Nameserver, len(in))
	for i, ns := range in {
		apiNameservers[i] = nbapi.Nameserver{
			Ip:     ns.Ip,
			NsType: nbapi.NameserverNsType(ns.NsType),
			Port:   ns.Port,
		}
	}

	return apiNameservers
}

// Converts a slice of API Nameserver to internal Nameserver.
func fromAPINameservers(in []nbapi.Nameserver) []Nameserver {
	nameservers := make([]Nameserver, len(in))
	for i, ns := range in {
		nameservers[i] = Nameserver{
			Ip:     ns.Ip,
			NsType: NameserverNsType(ns.NsType),
			Port:   ns.Port,
		}
	}

	return nameservers
}

// Delete removes a DNS (Nameserver Group) from NetBird.
func (*DNS) Delete(ctx context.Context, req infer.DeleteRequest[DNSState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:DNS[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting Netbird client: %w", err)
	}

	err = client.DNS.DeleteNameserverGroup(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting DNS entry failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between DNSArgs and DNSState.
func (*DNS) Diff(ctx context.Context, req infer.DiffRequest[DNSArgs, DNSState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:DNS[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Description != req.State.Description {
		diff["description"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if !slices.Equal(req.Inputs.Domains, req.State.Domains) {
		diff["domains"] = p.PropertyDiff{
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

	if !slices.Equal(req.Inputs.Groups, req.State.Groups) {
		diff["groups"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Primary != req.State.Primary {
		diff["primary"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.SearchDomainsEnabled != req.State.SearchDomainsEnabled {
		diff["search_domains_enabled"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	// Compare nameservers
	if len(req.Inputs.Nameservers) != len(req.State.Nameservers) {
		diff["nameservers"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	} else {
		for i := range req.Inputs.Nameservers {
			in := req.Inputs.Nameservers[i]
			st := req.State.Nameservers[i]

			if in.Ip != st.Ip || in.NsType != st.NsType || in.Port != st.Port {
				diff["nameservers"] = p.PropertyDiff{
					InputDiff: false,
					Kind:      p.Update,
				}

				break
			}
		}
	}

	p.GetLogger(ctx).Debugf("Diff:DNS[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}
