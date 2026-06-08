package function

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// LookupSetupKey looks up an existing NetBird setup key by name.
type LookupSetupKey struct{}

// Annotate describes the function.
func (f *LookupSetupKey) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up an existing NetBird setup key by name and return its metadata (not the secret key value).")
}

// LookupSetupKeyArgs are the inputs for LookupSetupKey.
type LookupSetupKeyArgs struct {
	Name string `pulumi:"name"`
}

// Annotate provides field descriptions for LookupSetupKeyArgs.
func (a *LookupSetupKeyArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "The name of the setup key to look up.")
}

// LookupSetupKeyResult is the output of LookupSetupKey.
type LookupSetupKeyResult struct {
	ID           string   `pulumi:"setupKeyId"`
	Name         string   `pulumi:"name"`
	Type         string   `pulumi:"type"`
	State        string   `pulumi:"state"`
	Revoked      bool     `pulumi:"revoked"`
	Ephemeral    bool     `pulumi:"ephemeral"`
	UsageLimit   int      `pulumi:"usageLimit"`
	AutoGroups   []string `pulumi:"autoGroups"`
	Expires      string   `pulumi:"expires"`
	LastUsed     string   `pulumi:"lastUsed"`
}

// Annotate provides field descriptions for LookupSetupKeyResult.
func (r *LookupSetupKeyResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.ID, "The NetBird setup key ID.")
	ann.Describe(&r.Name, "The setup key name.")
	ann.Describe(&r.Type, "The setup key type: 'one-off' or 'reusable'.")
	ann.Describe(&r.State, "The setup key state: 'valid', 'overused', 'expired', or 'revoked'.")
	ann.Describe(&r.Revoked, "Whether the setup key has been revoked.")
	ann.Describe(&r.Ephemeral, "Whether peers registered with this key are ephemeral.")
	ann.Describe(&r.UsageLimit, "Maximum number of times the key can be used (0 = unlimited for reusable).")
	ann.Describe(&r.AutoGroups, "Group IDs automatically assigned to peers that register with this key.")
	ann.Describe(&r.Expires, "Key expiration timestamp in RFC3339 format.")
	ann.Describe(&r.LastUsed, "Timestamp of last key use in RFC3339 format.")
}

// Invoke looks up a setup key by name.
func (f *LookupSetupKey) Invoke(ctx context.Context, req infer.FunctionRequest[LookupSetupKeyArgs]) (infer.FunctionResponse[LookupSetupKeyResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupSetupKeyResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	keys, err := client.SetupKeys.List(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupSetupKeyResult]{}, fmt.Errorf("listing setup keys failed: %w", err)
	}

	for _, key := range keys {
		if key.Name != req.Input.Name {
			continue
		}

		autoGroups := slices.Clone(key.AutoGroups)
		slices.Sort(autoGroups)

		return infer.FunctionResponse[LookupSetupKeyResult]{
			Output: LookupSetupKeyResult{
				ID:         key.Id,
				Name:       key.Name,
				Type:       key.Type,
				State:      key.State,
				Revoked:    key.Revoked,
				Ephemeral:  key.Ephemeral,
				UsageLimit: key.UsageLimit,
				AutoGroups: autoGroups,
				Expires:    key.Expires.Format("2006-01-02T15:04:05Z07:00"),
				LastUsed:   key.LastUsed.Format("2006-01-02T15:04:05Z07:00"),
			},
		}, nil
	}

	return infer.FunctionResponse[LookupSetupKeyResult]{}, fmt.Errorf("setup key %q not found", req.Input.Name)
}
