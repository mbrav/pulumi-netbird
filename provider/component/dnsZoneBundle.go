package component

import (
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DNSRecordSpec holds configuration for a single DNS record within a DNSZoneBundle.
type DNSRecordSpec struct {
	Name    string `pulumi:"name"`
	Type    string `pulumi:"type"`
	Content string `pulumi:"content"`
	TTL     int    `pulumi:"ttl"`
}

// Annotate adds schema descriptions to DNSRecordSpec fields.
func (s *DNSRecordSpec) Annotate(a infer.Annotator) {
	a.Describe(&s.Name, "Fully-qualified DNS record name (e.g. api.corp.example.com).")
	a.Describe(&s.Type, "Record type: A, AAAA, or CNAME.")
	a.Describe(&s.Content, "Record value: an IP address for A/AAAA or a hostname for CNAME.")
	a.Describe(&s.TTL, "Time-to-live in seconds.")
}

// DNSZoneBundleArgs are the inputs for a DNSZoneBundle component.
type DNSZoneBundleArgs struct {
	Name               string          `pulumi:"name"`
	Domain             string          `pulumi:"domain"`
	Enabled            bool            `pulumi:"enabled"`
	EnableSearchDomain bool            `pulumi:"enableSearchDomain"`
	DistributionGroups []string        `pulumi:"distributionGroups"`
	Records            []DNSRecordSpec `pulumi:"records"`
}

// Annotate adds schema descriptions to DNSZoneBundleArgs fields.
func (d *DNSZoneBundleArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&d.Name, "Logical name for the DNS zone resource.")
	ann.Describe(&d.Domain, "DNS domain (e.g. corp.example.com).")
	ann.Describe(&d.Enabled, "Whether the DNS zone is active.")
	ann.Describe(&d.EnableSearchDomain, "Whether to enable the zone as a search domain for peers.")
	ann.Describe(&d.DistributionGroups, "Group IDs whose peers receive this DNS zone.")
	ann.Describe(&d.Records, "DNS records to create within the zone.")
}

// DNSZoneBundleState holds the outputs of a DNSZoneBundle component.
type DNSZoneBundleState struct {
	pulumi.ResourceState

	ZoneID    pulumi.StringOutput      `pulumi:"zoneId"`
	RecordIDs pulumi.StringArrayOutput `pulumi:"recordIds"`
}

// Annotate adds schema descriptions to DNSZoneBundleState fields.
func (s *DNSZoneBundleState) Annotate(a infer.Annotator) {
	a.Describe(&s.ZoneID, "ID of the created DNSZone resource.")
	a.Describe(&s.RecordIDs, "IDs of the created DNSRecord resources, in declaration order.")
}

// DNSZoneBundle is the ComponentResource anchor for the DNSZoneBundle component.
type DNSZoneBundle struct{}

// Construct implements infer.ComponentResource and creates the child resources.
func (*DNSZoneBundle) Construct(
	ctx *pulumi.Context, name, typ string,
	args DNSZoneBundleArgs, opts pulumi.ResourceOption,
) (*DNSZoneBundleState, error) {
	return newDNSZoneBundle(ctx, name, typ, args, opts)
}

func newDNSZoneBundle(
	ctx *pulumi.Context,
	name, typ string,
	args DNSZoneBundleArgs,
	opts ...pulumi.ResourceOption,
) (*DNSZoneBundleState, error) {
	comp := &DNSZoneBundleState{} //nolint:exhaustruct

	err := ctx.RegisterComponentResource(typ, name, comp, opts...)
	if err != nil {
		return nil, fmt.Errorf("registering DNSZoneBundle component: %w", err)
	}

	distributionGroups := make(pulumi.StringArray, len(args.DistributionGroups))
	for j, gid := range args.DistributionGroups {
		distributionGroups[j] = pulumi.String(gid)
	}

	zoneInputs := pulumi.Map{
		"name":               pulumi.String(args.Name),
		"domain":             pulumi.String(args.Domain),
		"enabled":            pulumi.Bool(args.Enabled),
		"enableSearchDomain": pulumi.Bool(args.EnableSearchDomain),
		"distributionGroups": distributionGroups,
	}

	var zone pulumi.CustomResourceState

	err = ctx.RegisterResource(tokenDNSZone, name+"-zone", zoneInputs, &zone, pulumi.Parent(comp))
	if err != nil {
		return nil, fmt.Errorf("creating DNSZone: %w", err)
	}

	recordIDs := make(pulumi.StringArray, len(args.Records))

	for recIdx, rec := range args.Records {
		recordInputs := pulumi.Map{
			"zoneID":  zone.ID().ToStringOutput(),
			"name":    pulumi.String(rec.Name),
			"type":    pulumi.String(rec.Type),
			"content": pulumi.String(rec.Content),
			"ttl":     pulumi.Int(rec.TTL),
		}

		var recChild pulumi.CustomResourceState

		err = ctx.RegisterResource(
			tokenDNSRecord,
			name+"-record-"+rec.Name,
			recordInputs,
			&recChild,
			pulumi.Parent(comp),
		)
		if err != nil {
			return nil, fmt.Errorf("creating DNSRecord %q: %w", rec.Name, err)
		}

		recordIDs[recIdx] = recChild.ID().ToStringOutput()
	}

	comp.ZoneID = zone.ID().ToStringOutput()
	comp.RecordIDs = recordIDs.ToStringArrayOutput()

	return comp, nil
}
