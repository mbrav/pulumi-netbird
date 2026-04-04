package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// DNSRecord represents a NetBird DNS record resource.
type DNSRecord struct{}

// Annotate adds a description to the DNSRecord resource type.
func (d *DNSRecord) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d, "A DNS record within a NetBird DNS zone.")
}

// DNSRecordArgs defines input fields for creating or updating a DNS record.
type DNSRecordArgs struct {
	ZoneID  string        `pulumi:"zoneID"`
	Name    string        `pulumi:"name"`
	Content string        `pulumi:"content"`
	TTL     int           `pulumi:"ttl"`
	Type    DNSRecordType `pulumi:"type"`
}

// Annotate provides documentation for DNSRecordArgs fields.
func (d *DNSRecordArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d.ZoneID, "ID of the DNS zone this record belongs to.")
	annotator.Describe(&d.Name, "FQDN for the DNS record. Must be a subdomain within or match the zone's domain.")
	annotator.Describe(&d.Content, "DNS record content (IP address for A/AAAA, domain for CNAME).")
	annotator.Describe(&d.TTL, "Time to live in seconds.")
	annotator.Describe(&d.Type, "DNS record type.")
}

// DNSRecordState represents the output state of a DNS record resource.
type DNSRecordState struct {
	ZoneID  string        `pulumi:"zoneID"`
	Name    string        `pulumi:"name"`
	Content string        `pulumi:"content"`
	TTL     int           `pulumi:"ttl"`
	Type    DNSRecordType `pulumi:"type"`
}

// Annotate provides documentation for DNSRecordState fields.
func (d *DNSRecordState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d.ZoneID, "ID of the DNS zone this record belongs to.")
	annotator.Describe(&d.Name, "FQDN for the DNS record. Must be a subdomain within or match the zone's domain.")
	annotator.Describe(&d.Content, "DNS record content (IP address for A/AAAA, domain for CNAME).")
	annotator.Describe(&d.TTL, "Time to live in seconds.")
	annotator.Describe(&d.Type, "DNS record type.")
}

// DNSRecordType defines the allowed DNS record types.
// This wraps the nbapi type to allow method definitions (like Values()).
type DNSRecordType string

const (
	// DNSRecordTypeA represents an A record.
	DNSRecordTypeA DNSRecordType = DNSRecordType(nbapi.DNSRecordTypeA)
	// DNSRecordTypeAAAA represents an AAAA record.
	DNSRecordTypeAAAA DNSRecordType = DNSRecordType(nbapi.DNSRecordTypeAAAA)
	// DNSRecordTypeCNAME represents a CNAME record.
	DNSRecordTypeCNAME DNSRecordType = DNSRecordType(nbapi.DNSRecordTypeCNAME)
)

// Values returns the valid enum values for DNSRecordType, used by Pulumi for schema generation and validation.
func (DNSRecordType) Values() []infer.EnumValue[DNSRecordType] {
	return []infer.EnumValue[DNSRecordType]{
		{Name: "A", Value: DNSRecordTypeA, Description: "IPv4 address record."},
		{Name: "AAAA", Value: DNSRecordTypeAAAA, Description: "IPv6 address record."},
		{Name: "CNAME", Value: DNSRecordTypeCNAME, Description: "Canonical name record."},
	}
}

// Create creates a new DNS record in a NetBird DNS zone.
func (*DNSRecord) Create(ctx context.Context, req infer.CreateRequest[DNSRecordArgs]) (infer.CreateResponse[DNSRecordState], error) {
	p.GetLogger(ctx).Debugf("Create:DNSRecord name=%s, zone_id=%s", req.Inputs.Name, req.Inputs.ZoneID)

	if req.DryRun {
		return infer.CreateResponse[DNSRecordState]{
			ID: "preview",
			Output: DNSRecordState{
				ZoneID:  req.Inputs.ZoneID,
				Name:    req.Inputs.Name,
				Content: req.Inputs.Content,
				TTL:     req.Inputs.TTL,
				Type:    req.Inputs.Type,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[DNSRecordState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	record, err := client.DNSZones.CreateRecord(ctx, req.Inputs.ZoneID, nbapi.DNSRecordRequest{
		Name:    req.Inputs.Name,
		Content: req.Inputs.Content,
		Ttl:     req.Inputs.TTL,
		Type:    nbapi.DNSRecordType(req.Inputs.Type),
	})
	if err != nil {
		return infer.CreateResponse[DNSRecordState]{}, fmt.Errorf("creating DNS record failed: %w", err)
	}

	return infer.CreateResponse[DNSRecordState]{
		ID: record.Id,
		Output: DNSRecordState{
			ZoneID:  req.Inputs.ZoneID,
			Name:    record.Name,
			Content: record.Content,
			TTL:     record.Ttl,
			Type:    DNSRecordType(record.Type),
		},
	}, nil
}

// Read reads a DNS record from NetBird.
func (*DNSRecord) Read(ctx context.Context, req infer.ReadRequest[DNSRecordArgs, DNSRecordState]) (infer.ReadResponse[DNSRecordArgs, DNSRecordState], error) {
	p.GetLogger(ctx).Debugf("Read:DNSRecord[%s] zone_id=%s", req.ID, req.State.ZoneID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[DNSRecordArgs, DNSRecordState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	record, err := client.DNSZones.GetRecord(ctx, req.State.ZoneID, req.ID)
	if err != nil {
		return infer.ReadResponse[DNSRecordArgs, DNSRecordState]{}, fmt.Errorf("reading DNS record failed: %w", err)
	}

	return infer.ReadResponse[DNSRecordArgs, DNSRecordState]{
		ID: req.ID,
		Inputs: DNSRecordArgs{
			ZoneID:  req.State.ZoneID,
			Name:    record.Name,
			Content: record.Content,
			TTL:     record.Ttl,
			Type:    DNSRecordType(record.Type),
		},
		State: DNSRecordState{
			ZoneID:  req.State.ZoneID,
			Name:    record.Name,
			Content: record.Content,
			TTL:     record.Ttl,
			Type:    DNSRecordType(record.Type),
		},
	}, nil
}

// Update updates a DNS record in NetBird.
func (*DNSRecord) Update(ctx context.Context, req infer.UpdateRequest[DNSRecordArgs, DNSRecordState]) (infer.UpdateResponse[DNSRecordState], error) {
	p.GetLogger(ctx).Debugf("Update:DNSRecord[%s] zone_id=%s", req.ID, req.State.ZoneID)

	if req.DryRun {
		return infer.UpdateResponse[DNSRecordState]{
			Output: DNSRecordState{
				ZoneID:  req.Inputs.ZoneID,
				Name:    req.Inputs.Name,
				Content: req.Inputs.Content,
				TTL:     req.Inputs.TTL,
				Type:    req.Inputs.Type,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[DNSRecordState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	record, err := client.DNSZones.UpdateRecord(ctx, req.State.ZoneID, req.ID, nbapi.DNSRecordRequest{
		Name:    req.Inputs.Name,
		Content: req.Inputs.Content,
		Ttl:     req.Inputs.TTL,
		Type:    nbapi.DNSRecordType(req.Inputs.Type),
	})
	if err != nil {
		return infer.UpdateResponse[DNSRecordState]{}, fmt.Errorf("updating DNS record failed: %w", err)
	}

	return infer.UpdateResponse[DNSRecordState]{
		Output: DNSRecordState{
			ZoneID:  req.State.ZoneID,
			Name:    record.Name,
			Content: record.Content,
			TTL:     record.Ttl,
			Type:    DNSRecordType(record.Type),
		},
	}, nil
}

// Delete removes a DNS record from NetBird.
func (*DNSRecord) Delete(ctx context.Context, req infer.DeleteRequest[DNSRecordState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:DNSRecord[%s] zone_id=%s", req.ID, req.State.ZoneID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.DNSZones.DeleteRecord(ctx, req.State.ZoneID, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting DNS record failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between DNSRecordArgs and DNSRecordState.
func (*DNSRecord) Diff(ctx context.Context, req infer.DiffRequest[DNSRecordArgs, DNSRecordState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:DNSRecord[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.ZoneID != req.State.ZoneID {
		diff["zoneID"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Content != req.State.Content {
		diff["content"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.TTL != req.State.TTL {
		diff["ttl"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Type != req.State.Type {
		diff["type"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:DNSRecord[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a DNS record.
func (*DNSRecord) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[DNSRecordArgs], error) {
	p.GetLogger(ctx).Debugf("Check:DNSRecord old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[DNSRecordArgs](ctx, req.NewInputs)

	if isBlank(args.ZoneID) {
		failures = append(failures, p.CheckFailure{
			Property: "zoneID",
			Reason:   "zoneID must not be empty",
		})
	}

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{
			Property: "name",
			Reason:   "name must not be empty",
		})
	}

	if isBlank(args.Content) {
		failures = append(failures, p.CheckFailure{
			Property: "content",
			Reason:   "content must not be empty",
		})
	}

	if args.TTL < 1 {
		failures = append(failures, p.CheckFailure{
			Property: "ttl",
			Reason:   "ttl must be greater than 0",
		})
	}

	return infer.CheckResponse[DNSRecordArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*DNSRecord) WireDependencies(field infer.FieldSelector, args *DNSRecordArgs, state *DNSRecordState) {
	field.OutputField(&state.ZoneID).DependsOn(field.InputField(&args.ZoneID))
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Content).DependsOn(field.InputField(&args.Content))
	field.OutputField(&state.TTL).DependsOn(field.InputField(&args.TTL))
	field.OutputField(&state.Type).DependsOn(field.InputField(&args.Type))
}
