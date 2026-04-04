package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ReverseProxyDomain represents a NetBird reverse proxy custom domain resource.
type ReverseProxyDomain struct{}

// Annotate adds a description to the ReverseProxyDomain resource type.
func (r *ReverseProxyDomain) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r, "A NetBird reverse proxy custom domain.")
}

// ReverseProxyDomainArgs defines input fields for creating a reverse proxy domain.
type ReverseProxyDomainArgs struct {
	Domain        string `pulumi:"domain"`
	TargetCluster string `pulumi:"targetCluster"`
}

// Annotate provides documentation for ReverseProxyDomainArgs fields.
func (r *ReverseProxyDomainArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Domain, "Domain name for the reverse proxy.")
	annotator.Describe(&r.TargetCluster, "The proxy cluster this domain should be validated against.")
}

// ReverseProxyDomainState represents the output state of a reverse proxy domain resource.
type ReverseProxyDomainState struct {
	Domain              string                 `pulumi:"domain"`
	TargetCluster       string                 `pulumi:"targetCluster"`
	Type                ReverseProxyDomainType `pulumi:"type"`
	Validated           bool                   `pulumi:"validated"`
	RequireSubdomain    *bool                  `pulumi:"requireSubdomain,optional"`
	SupportsCustomPorts *bool                  `pulumi:"supportsCustomPorts,optional"`
}

// Annotate provides documentation for ReverseProxyDomainState fields.
func (r *ReverseProxyDomainState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Domain, "Domain name for the reverse proxy.")
	annotator.Describe(&r.TargetCluster, "The proxy cluster this domain is validated against.")
	annotator.Describe(&r.Type, "Type of the reverse proxy domain (custom or free).")
	annotator.Describe(&r.Validated, "Whether the domain has been validated.")
	annotator.Describe(&r.RequireSubdomain, "Whether a subdomain label is required in front of this domain.")
	annotator.Describe(&r.SupportsCustomPorts, "Whether the cluster supports binding arbitrary TCP/UDP ports.")
}

// ReverseProxyDomainType defines the allowed domain types.
type ReverseProxyDomainType string

const (
	// ReverseProxyDomainTypeCustom represents a custom domain.
	ReverseProxyDomainTypeCustom ReverseProxyDomainType = ReverseProxyDomainType(nbapi.ReverseProxyDomainTypeCustom)
	// ReverseProxyDomainTypeFree represents a free (managed) domain.
	ReverseProxyDomainTypeFree ReverseProxyDomainType = ReverseProxyDomainType(nbapi.ReverseProxyDomainTypeFree)
)

// Values returns the valid enum values for ReverseProxyDomainType.
func (ReverseProxyDomainType) Values() []infer.EnumValue[ReverseProxyDomainType] {
	return []infer.EnumValue[ReverseProxyDomainType]{
		{Name: "custom", Value: ReverseProxyDomainTypeCustom, Description: "A custom domain managed by the user."},
		{Name: "free", Value: ReverseProxyDomainTypeFree, Description: "A free managed domain provided by NetBird."},
	}
}

// Create creates a new reverse proxy custom domain.
func (*ReverseProxyDomain) Create(ctx context.Context, req infer.CreateRequest[ReverseProxyDomainArgs]) (infer.CreateResponse[ReverseProxyDomainState], error) {
	p.GetLogger(ctx).Debugf("Create:ReverseProxyDomain domain=%s, cluster=%s", req.Inputs.Domain, req.Inputs.TargetCluster)

	if req.DryRun {
		return infer.CreateResponse[ReverseProxyDomainState]{
			ID: "preview",
			Output: ReverseProxyDomainState{
				Domain:              req.Inputs.Domain,
				TargetCluster:       req.Inputs.TargetCluster,
				Type:                "",
				Validated:           false,
				RequireSubdomain:    nil,
				SupportsCustomPorts: nil,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[ReverseProxyDomainState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	domain, err := client.ReverseProxyDomains.Create(ctx, nbapi.ReverseProxyDomainRequest{
		Domain:        req.Inputs.Domain,
		TargetCluster: req.Inputs.TargetCluster,
	})
	if err != nil {
		return infer.CreateResponse[ReverseProxyDomainState]{}, fmt.Errorf("creating reverse proxy domain failed: %w", err)
	}

	targetCluster := ""
	if domain.TargetCluster != nil {
		targetCluster = *domain.TargetCluster
	}

	return infer.CreateResponse[ReverseProxyDomainState]{
		ID: domain.Id,
		Output: ReverseProxyDomainState{
			Domain:              domain.Domain,
			TargetCluster:       targetCluster,
			Type:                ReverseProxyDomainType(domain.Type),
			Validated:           domain.Validated,
			RequireSubdomain:    domain.RequireSubdomain,
			SupportsCustomPorts: domain.SupportsCustomPorts,
		},
	}, nil
}

// Read reads a reverse proxy domain by listing all and finding by ID.
// The API does not provide a Get-by-ID endpoint for domains.
func (*ReverseProxyDomain) Read(ctx context.Context, req infer.ReadRequest[ReverseProxyDomainArgs, ReverseProxyDomainState]) (infer.ReadResponse[ReverseProxyDomainArgs, ReverseProxyDomainState], error) {
	p.GetLogger(ctx).Debugf("Read:ReverseProxyDomain[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[ReverseProxyDomainArgs, ReverseProxyDomainState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	domains, err := client.ReverseProxyDomains.List(ctx)
	if err != nil {
		return infer.ReadResponse[ReverseProxyDomainArgs, ReverseProxyDomainState]{}, fmt.Errorf("reading reverse proxy domains failed: %w", err)
	}

	var found *nbapi.ReverseProxyDomain

	for i := range domains {
		if domains[i].Id == req.ID {
			found = &domains[i]

			break
		}
	}

	if found == nil {
		return infer.ReadResponse[ReverseProxyDomainArgs, ReverseProxyDomainState]{}, fmt.Errorf("reverse proxy domain with ID %s not found", req.ID)
	}

	targetCluster := ""
	if found.TargetCluster != nil {
		targetCluster = *found.TargetCluster
	}

	return infer.ReadResponse[ReverseProxyDomainArgs, ReverseProxyDomainState]{
		ID: req.ID,
		Inputs: ReverseProxyDomainArgs{
			Domain:        found.Domain,
			TargetCluster: targetCluster,
		},
		State: ReverseProxyDomainState{
			Domain:              found.Domain,
			TargetCluster:       targetCluster,
			Type:                ReverseProxyDomainType(found.Type),
			Validated:           found.Validated,
			RequireSubdomain:    found.RequireSubdomain,
			SupportsCustomPorts: found.SupportsCustomPorts,
		},
	}, nil
}

// Delete removes a reverse proxy domain.
func (*ReverseProxyDomain) Delete(ctx context.Context, req infer.DeleteRequest[ReverseProxyDomainState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:ReverseProxyDomain[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.ReverseProxyDomains.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting reverse proxy domain failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between ReverseProxyDomainArgs and ReverseProxyDomainState.
// Since the API has no Update endpoint, any input change forces a replacement.
func (*ReverseProxyDomain) Diff(ctx context.Context, req infer.DiffRequest[ReverseProxyDomainArgs, ReverseProxyDomainState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:ReverseProxyDomain[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Domain != req.State.Domain {
		diff["domain"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.TargetCluster != req.State.TargetCluster {
		diff["targetCluster"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	p.GetLogger(ctx).Debugf("Diff:ReverseProxyDomain[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a reverse proxy domain.
func (*ReverseProxyDomain) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[ReverseProxyDomainArgs], error) {
	p.GetLogger(ctx).Debugf("Check:ReverseProxyDomain old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[ReverseProxyDomainArgs](ctx, req.NewInputs)

	if isBlank(args.Domain) {
		failures = append(failures, p.CheckFailure{
			Property: "domain",
			Reason:   "domain must not be empty",
		})
	}

	if isBlank(args.TargetCluster) {
		failures = append(failures, p.CheckFailure{
			Property: "targetCluster",
			Reason:   "targetCluster must not be empty",
		})
	}

	return infer.CheckResponse[ReverseProxyDomainArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*ReverseProxyDomain) WireDependencies(field infer.FieldSelector, args *ReverseProxyDomainArgs, state *ReverseProxyDomainState) {
	field.OutputField(&state.Domain).DependsOn(field.InputField(&args.Domain))
	field.OutputField(&state.TargetCluster).DependsOn(field.InputField(&args.TargetCluster))
}
