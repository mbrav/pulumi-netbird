package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

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
	a.Describe(&n.Name, "Name Name of nameserver group name")
	a.Describe(&n.Description, "Description Description of the nameserver group")
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
	a.Describe(&n.Name, "Name Name of nameserver group name")
	a.Describe(&n.Description, "Description Description of the nameserver group")
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
		{Name: "UDP", Value: NameserverNsTypeUdp, Description: "UDP type"},
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
		return infer.CreateResponse[DNSState]{}, err
	}

	// Convert to API Nameserver slice
	apiNameservers := make([]nbapi.Nameserver, len(req.Inputs.Nameservers))
	for i, ns := range req.Inputs.Nameservers {
		apiNameservers[i] = nbapi.Nameserver{
			Ip:     ns.Ip,
			NsType: nbapi.NameserverNsType(ns.NsType),
			Port:   ns.Port,
		}
	}

	// Build request payload
	createReq := nbapi.NameserverGroupRequest{
		Name:                 req.Inputs.Name,
		Description:          req.Inputs.Description,
		Domains:              req.Inputs.Domains,
		Enabled:              req.Inputs.Enabled,
		Groups:               req.Inputs.Groups,
		Primary:              req.Inputs.Primary,
		Nameservers:          apiNameservers,
		SearchDomainsEnabled: req.Inputs.SearchDomainsEnabled,
	}

	// Call the API
	created, err := client.DNS.CreateNameserverGroup(ctx, createReq)
	if err != nil {
		return infer.CreateResponse[DNSState]{}, fmt.Errorf("creating DNS group failed: %w", err)
	}

	// Convert API response back to internal state
	stateNameservers := make([]Nameserver, len(created.Nameservers))
	for i, ns := range created.Nameservers {
		stateNameservers[i] = Nameserver{
			Ip:     ns.Ip,
			NsType: NameserverNsType(ns.NsType),
			Port:   ns.Port,
		}
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
			Nameservers:          stateNameservers,
			SearchDomainsEnabled: created.SearchDomainsEnabled,
		},
	}, nil
}

// // Read fetches the current state of a network from NetBird.
// func (*DNS) Read(ctx context.Context, req infer.ReadRequest[DNSArgs, DNSState]) (infer.ReadResponse[DNSArgs, DNSState], error) {
// 	p.GetLogger(ctx).Debugf("Read:DNSArgs[%s] name=%s", req.ID, req.Inputs.Name)
// 	p.GetLogger(ctx).Debugf("Read:DNSState[%s] name=%s, id=%s", req.ID, req.State.Name, req.ID)
//
// 	client, err := config.GetNetBirdClient(ctx)
// 	if err != nil {
// 		return infer.ReadResponse[DNSArgs, DNSState]{}, err
// 	}
//
// 	net, err := client.DNSs.Get(ctx, req.ID)
// 	if err != nil {
// 		return infer.ReadResponse[DNSArgs, DNSState]{}, fmt.Errorf("reading network failed: %w", err)
// 	}
//
// 	p.GetLogger(ctx).Debugf("Read:DNSAPI[%s] name=%s", net.Id, net.Name)
//
// 	return infer.ReadResponse[DNSArgs, DNSState]{
// 		ID: req.ID,
// 		Inputs: DNSArgs{
// 			Name:        net.Name,
// 			Description: net.Description,
// 		},
// 		State: DNSState{
// 			Name:        net.Name,
// 			Description: net.Description,
// 		},
// 	}, nil
// }
//
// // Update updates the state of the network if needed.
// func (*DNS) Update(ctx context.Context, req infer.UpdateRequest[DNSArgs, DNSState]) (infer.UpdateResponse[DNSState], error) {
// 	p.GetLogger(ctx).Debugf("Update:DNS[%s] name=%s", req.ID, req.Inputs.Name)
//
// 	if req.DryRun {
// 		return infer.UpdateResponse[DNSState]{
// 			Output: DNSState{
// 				Name:        req.Inputs.Name,
// 				Description: req.Inputs.Description,
// 			},
// 		}, nil
// 	}
//
// 	client, err := config.GetNetBirdClient(ctx)
// 	if err != nil {
// 		return infer.UpdateResponse[DNSState]{}, err
// 	}
//
// 	_, err = client.DNSs.Update(ctx, req.ID, nbapi.DNSRequest{
// 		Name:        req.Inputs.Name,
// 		Description: req.Inputs.Description,
// 	})
// 	if err != nil {
// 		return infer.UpdateResponse[DNSState]{}, fmt.Errorf("updating network failed: %w", err)
// 	}
//
// 	return infer.UpdateResponse[DNSState]{
// 		Output: DNSState{
// 			Name:        req.Inputs.Name,
// 			Description: req.Inputs.Description,
// 		},
// 	}, nil
// }
//
// // Delete removes a network from NetBird.
// func (*DNS) Delete(ctx context.Context, req infer.DeleteRequest[DNSState]) (infer.DeleteResponse, error) {
// 	p.GetLogger(ctx).Debugf("Delete:DNS[%s]", req.ID)
//
// 	client, err := config.GetNetBirdClient(ctx)
// 	if err != nil {
// 		return infer.DeleteResponse{}, err
// 	}
//
// 	err = client.DNSs.Delete(ctx, req.ID)
// 	if err != nil {
// 		return infer.DeleteResponse{}, fmt.Errorf("deleting network failed: %w", err)
// 	}
//
// 	return infer.DeleteResponse{}, nil
// }
//
// // Diff detects changes between inputs and prior state.
// func (*DNS) Diff(ctx context.Context, req infer.DiffRequest[DNSArgs, DNSState]) (infer.DiffResponse, error) {
// 	p.GetLogger(ctx).Debugf("Diff:DNS[%s]", req.ID)
//
// 	diff := map[string]p.PropertyDiff{}
//
// 	if req.Inputs.Name != req.State.Name {
// 		diff["name"] = p.PropertyDiff{Kind: p.Update}
// 	}
//
// 	if !equalPtr(req.Inputs.Description, req.State.Description) {
// 		diff["description"] = p.PropertyDiff{Kind: p.Update}
// 	}
//
// 	p.GetLogger(ctx).Debugf("Diff:DNS[%s] diff=%d", req.ID, len(diff))
//
// 	return infer.DiffResponse{
// 		DeleteBeforeReplace: false,
// 		HasChanges:          len(diff) > 0,
// 		DetailedDiff:        diff,
// 	}, nil
// }
//
// // Check provides input validation and default setting.
// func (*DNS) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[DNSArgs], error) {
// 	p.GetLogger(ctx).Debugf("Check:DNS old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())
// 	args, failures, err := infer.DefaultCheck[DNSArgs](ctx, req.NewInputs)
//
// 	return infer.CheckResponse[DNSArgs]{
// 		Inputs:   args,
// 		Failures: failures,
// 	}, err
// }
//
// // WireDependencies explicitly defines input/output relationships.
// func (*DNS) WireDependencies(f infer.FieldSelector, args *DNSArgs, state *DNSState) {
// 	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
// 	f.OutputField(&state.Description).DependsOn(f.InputField(&args.Description))
// }
