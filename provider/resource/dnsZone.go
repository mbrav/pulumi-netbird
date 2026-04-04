package resource

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// DNSZone represents a NetBird DNS zone resource.
type DNSZone struct{}

// Annotate adds a description to the DNSZone resource type.
func (d *DNSZone) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d, "A NetBird DNS zone.")
}

// DNSZoneArgs defines input fields for creating or updating a DNS zone.
type DNSZoneArgs struct {
	Name               string   `pulumi:"name"`
	Domain             string   `pulumi:"domain"`
	Enabled            bool     `pulumi:"enabled"`
	EnableSearchDomain bool     `pulumi:"enableSearchDomain"`
	DistributionGroups []string `pulumi:"distributionGroups"`
}

// Annotate provides documentation for DNSZoneArgs fields.
func (d *DNSZoneArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d.Name, "Zone name identifier.")
	annotator.Describe(&d.Domain, "Zone domain (FQDN).")
	annotator.Describe(&d.Enabled, "Zone status.")
	annotator.Describe(&d.EnableSearchDomain, "Enable this zone as a search domain.")
	annotator.Describe(&d.DistributionGroups, "Group IDs that define groups of peers that will resolve this zone.")
}

// DNSZoneState represents the output state of a DNS zone resource.
type DNSZoneState struct {
	Name               string   `pulumi:"name"`
	Domain             string   `pulumi:"domain"`
	Enabled            bool     `pulumi:"enabled"`
	EnableSearchDomain bool     `pulumi:"enableSearchDomain"`
	DistributionGroups []string `pulumi:"distributionGroups"`
}

// Annotate provides documentation for DNSZoneState fields.
func (d *DNSZoneState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d.Name, "Zone name identifier.")
	annotator.Describe(&d.Domain, "Zone domain (FQDN).")
	annotator.Describe(&d.Enabled, "Zone status.")
	annotator.Describe(&d.EnableSearchDomain, "Enable this zone as a search domain.")
	annotator.Describe(&d.DistributionGroups, "Group IDs that define groups of peers that will resolve this zone.")
}

// Create creates a new NetBird DNS zone.
func (*DNSZone) Create(ctx context.Context, req infer.CreateRequest[DNSZoneArgs]) (infer.CreateResponse[DNSZoneState], error) {
	p.GetLogger(ctx).Debugf("Create:DNSZone name=%s, domain=%s", req.Inputs.Name, req.Inputs.Domain)

	slices.Sort(req.Inputs.DistributionGroups)

	if req.DryRun {
		return infer.CreateResponse[DNSZoneState]{
			ID: "preview",
			Output: DNSZoneState{
				Name:               req.Inputs.Name,
				Domain:             req.Inputs.Domain,
				Enabled:            req.Inputs.Enabled,
				EnableSearchDomain: req.Inputs.EnableSearchDomain,
				DistributionGroups: req.Inputs.DistributionGroups,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[DNSZoneState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	enabled := req.Inputs.Enabled

	zone, err := client.DNSZones.CreateZone(ctx, nbapi.ZoneRequest{
		Name:               req.Inputs.Name,
		Domain:             req.Inputs.Domain,
		Enabled:            &enabled,
		EnableSearchDomain: req.Inputs.EnableSearchDomain,
		DistributionGroups: req.Inputs.DistributionGroups,
	})
	if err != nil {
		return infer.CreateResponse[DNSZoneState]{}, fmt.Errorf("creating DNS zone failed: %w", err)
	}

	return infer.CreateResponse[DNSZoneState]{
		ID: zone.Id,
		Output: DNSZoneState{
			Name:               zone.Name,
			Domain:             zone.Domain,
			Enabled:            zone.Enabled,
			EnableSearchDomain: zone.EnableSearchDomain,
			DistributionGroups: zone.DistributionGroups,
		},
	}, nil
}

// Read reads a DNS zone from NetBird.
func (*DNSZone) Read(ctx context.Context, req infer.ReadRequest[DNSZoneArgs, DNSZoneState]) (infer.ReadResponse[DNSZoneArgs, DNSZoneState], error) {
	p.GetLogger(ctx).Debugf("Read:DNSZone[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[DNSZoneArgs, DNSZoneState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	zone, err := client.DNSZones.GetZone(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[DNSZoneArgs, DNSZoneState]{}, fmt.Errorf("reading DNS zone failed: %w", err)
	}

	return infer.ReadResponse[DNSZoneArgs, DNSZoneState]{
		ID: req.ID,
		Inputs: DNSZoneArgs{
			Name:               zone.Name,
			Domain:             zone.Domain,
			Enabled:            zone.Enabled,
			EnableSearchDomain: zone.EnableSearchDomain,
			DistributionGroups: zone.DistributionGroups,
		},
		State: DNSZoneState{
			Name:               zone.Name,
			Domain:             zone.Domain,
			Enabled:            zone.Enabled,
			EnableSearchDomain: zone.EnableSearchDomain,
			DistributionGroups: zone.DistributionGroups,
		},
	}, nil
}

// Update updates a DNS zone in NetBird.
func (*DNSZone) Update(ctx context.Context, req infer.UpdateRequest[DNSZoneArgs, DNSZoneState]) (infer.UpdateResponse[DNSZoneState], error) {
	p.GetLogger(ctx).Debugf("Update:DNSZone[%s]", req.ID)

	slices.Sort(req.Inputs.DistributionGroups)

	if req.DryRun {
		return infer.UpdateResponse[DNSZoneState]{
			Output: DNSZoneState{
				Name:               req.Inputs.Name,
				Domain:             req.Inputs.Domain,
				Enabled:            req.Inputs.Enabled,
				EnableSearchDomain: req.Inputs.EnableSearchDomain,
				DistributionGroups: req.Inputs.DistributionGroups,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[DNSZoneState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	enabled := req.Inputs.Enabled

	zone, err := client.DNSZones.UpdateZone(ctx, req.ID, nbapi.ZoneRequest{
		Name:               req.Inputs.Name,
		Domain:             req.Inputs.Domain,
		Enabled:            &enabled,
		EnableSearchDomain: req.Inputs.EnableSearchDomain,
		DistributionGroups: req.Inputs.DistributionGroups,
	})
	if err != nil {
		return infer.UpdateResponse[DNSZoneState]{}, fmt.Errorf("updating DNS zone failed: %w", err)
	}

	return infer.UpdateResponse[DNSZoneState]{
		Output: DNSZoneState{
			Name:               zone.Name,
			Domain:             zone.Domain,
			Enabled:            zone.Enabled,
			EnableSearchDomain: zone.EnableSearchDomain,
			DistributionGroups: zone.DistributionGroups,
		},
	}, nil
}

// Delete removes a DNS zone from NetBird.
func (*DNSZone) Delete(ctx context.Context, req infer.DeleteRequest[DNSZoneState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:DNSZone[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.DNSZones.DeleteZone(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting DNS zone failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between DNSZoneArgs and DNSZoneState.
func (*DNSZone) Diff(ctx context.Context, req infer.DiffRequest[DNSZoneArgs, DNSZoneState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:DNSZone[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Domain != req.State.Domain {
		diff["domain"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.EnableSearchDomain != req.State.EnableSearchDomain {
		diff["enableSearchDomain"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlice(req.Inputs.DistributionGroups, req.State.DistributionGroups) {
		diff["distributionGroups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:DNSZone[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a DNS zone.
func (*DNSZone) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[DNSZoneArgs], error) {
	p.GetLogger(ctx).Debugf("Check:DNSZone old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[DNSZoneArgs](ctx, req.NewInputs)

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{
			Property: "name",
			Reason:   "name must not be empty",
		})
	}

	if isBlank(args.Domain) {
		failures = append(failures, p.CheckFailure{
			Property: "domain",
			Reason:   "domain must not be empty",
		})
	}

	return infer.CheckResponse[DNSZoneArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*DNSZone) WireDependencies(field infer.FieldSelector, args *DNSZoneArgs, state *DNSZoneState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Domain).DependsOn(field.InputField(&args.Domain))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.EnableSearchDomain).DependsOn(field.InputField(&args.EnableSearchDomain))
	field.OutputField(&state.DistributionGroups).DependsOn(field.InputField(&args.DistributionGroups))
}
