package provider

import (
	"context"
	"fmt"

	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Policy represents a resource for managing NetBird policies.
type Policy struct{}

// PolicyArgs are the input arguments for a policy resource.
type PolicyArgs struct {
	Name                string                   `pulumi:"name"`
	Description         *string                  `pulumi:"description"`
	Enabled             bool                     `pulumi:"enabled"`
	Rules               []nbapi.PolicyRuleUpdate `pulumi:"rules"`
	SourcePostureChecks *[]string                `pulumi:"sourcePostureChecks"`
}

// PolicyState is the persisted state of the resource.
type PolicyState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	// PolicyArgs
	Name                string                   `pulumi:"name"`
	Description         *string                  `pulumi:"description"`
	Enabled             bool                     `pulumi:"enabled"`
	Rules               []nbapi.PolicyRuleUpdate `pulumi:"rules"`
	SourcePostureChecks *[]string                `pulumi:"sourcePostureChecks"`
	NbID                string                   `pulumi:"nbId"`
}

// Policy annotation
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
	a.Describe(&p.NbID, "The internal NetBird ID of the policy.")
}

func (Policy) Create(ctx context.Context, name string, input PolicyArgs, preview bool) (string, PolicyState, error) {
	state := PolicyState{
		Name:                input.Name,
		Description:         input.Description,
		Enabled:             input.Enabled,
		Rules:               input.Rules,
		SourcePostureChecks: input.SourcePostureChecks,
	}
	if preview {
		return name, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return name, state, err
	}

	created, err := client.Policies.Create(ctx, nbapi.PolicyUpdate{
		Name:                input.Name,
		Description:         input.Description,
		Enabled:             input.Enabled,
		Rules:               input.Rules,
		SourcePostureChecks: input.SourcePostureChecks,
	})
	if err != nil {
		return "", state, fmt.Errorf("creating policy failed: %w", err)
	}

	state.NbID = *created.Id
	return name, state, nil
}

func (Policy) Read(ctx context.Context, id string, inputs PolicyArgs, state PolicyState) (PolicyArgs, PolicyState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return inputs, state, err
	}

	pol, err := client.Policies.Get(ctx, state.NbID)
	if err != nil {
		return inputs, state, fmt.Errorf("reading policy failed: %w", err)
	}

	return PolicyArgs{
			Name:                pol.Name,
			Description:         pol.Description,
			Enabled:             pol.Enabled,
			Rules:               []nbapi.PolicyRuleUpdate{},
			SourcePostureChecks: &pol.SourcePostureChecks,
		}, PolicyState{
			NbID:                *pol.Id,
			Name:                pol.Name,
			Description:         pol.Description,
			Enabled:             pol.Enabled,
			Rules:               []nbapi.PolicyRuleUpdate{},
			SourcePostureChecks: &pol.SourcePostureChecks,
		}, nil
}

func (Policy) Update(ctx context.Context, id string, old PolicyArgs, new PolicyArgs, state PolicyState) (PolicyState, error) {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return state, err
	}

	updated, err := client.Policies.Update(ctx, state.NbID, nbapi.PolicyCreate{
		Name:                new.Name,
		Description:         new.Description,
		Enabled:             new.Enabled,
		Rules:               new.Rules,
		SourcePostureChecks: new.SourcePostureChecks,
	})
	if err != nil {
		return state, fmt.Errorf("updating policy failed: %w", err)
	}

	return PolicyState{
		NbID:                *updated.Id,
		Name:                updated.Name,
		Description:         updated.Description,
		Enabled:             updated.Enabled,
		Rules:               new.Rules,
		SourcePostureChecks: &updated.SourcePostureChecks,
	}, nil
}

func (Policy) Delete(ctx context.Context, id string, props PolicyState) error {
	client, err := getNetBirdClient(ctx)
	if err != nil {
		return err
	}

	if err := client.Policies.Delete(ctx, props.NbID); err != nil {
		return fmt.Errorf("deleting policy failed: %w", err)
	}
	return nil
}

func (Policy) Import(ctx context.Context, id string, input PolicyArgs, preview bool) (string, PolicyState, error) {
	state := PolicyState{}

	if preview {
		// We can't fetch data in preview, but we still set the ID
		state.NbID = id
		return id, state, nil
	}

	client, err := getNetBirdClient(ctx)
	if err != nil {
		return "", state, err
	}

	pol, err := client.Policies.Get(ctx, id)
	if err != nil {
		return "", state, fmt.Errorf("importing policy failed: %w", err)
	}

	state = PolicyState{
		NbID:                *pol.Id,
		Name:                pol.Name,
		Description:         pol.Description,
		Enabled:             pol.Enabled,
		Rules:               []nbapi.PolicyRuleUpdate{}, // not returned by GET
		SourcePostureChecks: &pol.SourcePostureChecks,
	}

	return id, state, nil
}
