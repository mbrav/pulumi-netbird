package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// DNSSettings represents the global DNS settings resource.
type DNSSettings struct{}

// Annotate adds a description to the DNSSettings resource type.
func (d *DNSSettings) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d, "NetBird global DNS settings. This is a singleton resource — only one instance exists per account.")
}

// DNSSettingsArgs defines input fields for DNS settings.
type DNSSettingsArgs struct {
	DisabledManagementGroups []string `pulumi:"disabledManagementGroups"`
}

// Annotate provides documentation for DNSSettingsArgs fields.
func (d *DNSSettingsArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d.DisabledManagementGroups, "Group IDs whose DNS management is disabled.")
}

// DNSSettingsState represents the output state of the DNS settings resource.
type DNSSettingsState struct {
	DisabledManagementGroups []string `pulumi:"disabledManagementGroups"`
}

// Annotate provides documentation for DNSSettingsState fields.
func (d *DNSSettingsState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&d.DisabledManagementGroups, "Group IDs whose DNS management is disabled.")
}

// Create initialises the DNS settings resource by applying the desired settings.
// Since this is a singleton, Create calls UpdateSettings and uses a fixed ID.
func (*DNSSettings) Create(ctx context.Context, req infer.CreateRequest[DNSSettingsArgs]) (infer.CreateResponse[DNSSettingsState], error) {
	p.GetLogger(ctx).Debugf("Create:DNSSettings disabledGroups=%v", req.Inputs.DisabledManagementGroups)

	if req.DryRun {
		return infer.CreateResponse[DNSSettingsState]{
			ID: "dns-settings",
			Output: DNSSettingsState{
				DisabledManagementGroups: req.Inputs.DisabledManagementGroups,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[DNSSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.DNS.UpdateSettings(ctx, nbapi.DNSSettings{
		DisabledManagementGroups: req.Inputs.DisabledManagementGroups,
	})
	if err != nil {
		return infer.CreateResponse[DNSSettingsState]{}, fmt.Errorf("creating DNS settings failed: %w", err)
	}

	return infer.CreateResponse[DNSSettingsState]{
		ID: "dns-settings",
		Output: DNSSettingsState{
			DisabledManagementGroups: updated.DisabledManagementGroups,
		},
	}, nil
}

// Read reads the current DNS settings from NetBird.
func (*DNSSettings) Read(ctx context.Context, req infer.ReadRequest[DNSSettingsArgs, DNSSettingsState]) (infer.ReadResponse[DNSSettingsArgs, DNSSettingsState], error) {
	p.GetLogger(ctx).Debugf("Read:DNSSettings[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[DNSSettingsArgs, DNSSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	settings, err := client.DNS.GetSettings(ctx)
	if err != nil {
		return infer.ReadResponse[DNSSettingsArgs, DNSSettingsState]{}, fmt.Errorf("reading DNS settings failed: %w", err)
	}

	return infer.ReadResponse[DNSSettingsArgs, DNSSettingsState]{
		ID: req.ID,
		Inputs: DNSSettingsArgs{
			DisabledManagementGroups: settings.DisabledManagementGroups,
		},
		State: DNSSettingsState{
			DisabledManagementGroups: settings.DisabledManagementGroups,
		},
	}, nil
}

// Update updates the DNS settings.
func (*DNSSettings) Update(ctx context.Context, req infer.UpdateRequest[DNSSettingsArgs, DNSSettingsState]) (infer.UpdateResponse[DNSSettingsState], error) {
	p.GetLogger(ctx).Debugf("Update:DNSSettings[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[DNSSettingsState]{
			Output: DNSSettingsState{
				DisabledManagementGroups: req.Inputs.DisabledManagementGroups,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[DNSSettingsState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.DNS.UpdateSettings(ctx, nbapi.DNSSettings{
		DisabledManagementGroups: req.Inputs.DisabledManagementGroups,
	})
	if err != nil {
		return infer.UpdateResponse[DNSSettingsState]{}, fmt.Errorf("updating DNS settings failed: %w", err)
	}

	return infer.UpdateResponse[DNSSettingsState]{
		Output: DNSSettingsState{
			DisabledManagementGroups: updated.DisabledManagementGroups,
		},
	}, nil
}

// Delete is a no-op because DNS settings are a singleton and cannot be deleted.
func (*DNSSettings) Delete(ctx context.Context, req infer.DeleteRequest[DNSSettingsState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:DNSSettings[%s] (no-op, singleton resource)", req.ID)

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between DNSSettingsArgs and DNSSettingsState.
func (*DNSSettings) Diff(ctx context.Context, req infer.DiffRequest[DNSSettingsArgs, DNSSettingsState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:DNSSettings[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if !equalSlice(req.Inputs.DisabledManagementGroups, req.State.DisabledManagementGroups) {
		diff["disabledManagementGroups"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	p.GetLogger(ctx).Debugf("Diff:DNSSettings[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// WireDependencies explicitly defines input/output relationships.
func (*DNSSettings) WireDependencies(field infer.FieldSelector, args *DNSSettingsArgs, state *DNSSettingsState) {
	field.OutputField(&state.DisabledManagementGroups).DependsOn(field.InputField(&args.DisabledManagementGroups))
}
