package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// AgentNetworkGuardrail represents a reusable Agent Network guardrail: a set of checks
// (model allowlist, prompt capture) that can be attached to one or more policies.
type AgentNetworkGuardrail struct{}

// Annotate adds a description to the AgentNetworkGuardrail resource type.
func (a *AgentNetworkGuardrail) Annotate(ann infer.Annotator) {
	ann.Describe(&a, "A NetBird Agent Network guardrail: a reusable set of checks (model "+
		"allowlist, prompt capture) attachable to one or more AgentNetworkPolicy resources.")
}

// AgentNetworkGuardrailModelAllowlist restricts requests to an explicit set of catalog model IDs.
type AgentNetworkGuardrailModelAllowlist struct {
	Enabled bool     `pulumi:"enabled"`
	Models  []string `pulumi:"models"`
}

// Annotate provides documentation for AgentNetworkGuardrailModelAllowlist fields.
func (m *AgentNetworkGuardrailModelAllowlist) Annotate(a infer.Annotator) {
	a.Describe(&m.Enabled, "Whether the model allowlist check is enforced.")
	a.Describe(&m.Models, "Allowed catalog model IDs. Requests for any other model are denied.")
}

// AgentNetworkGuardrailPromptCapture controls request/response prompt capture for the guardrail.
type AgentNetworkGuardrailPromptCapture struct {
	Enabled   bool `pulumi:"enabled"`
	RedactPii bool `pulumi:"redactPii"`
}

// Annotate provides documentation for AgentNetworkGuardrailPromptCapture fields.
func (c *AgentNetworkGuardrailPromptCapture) Annotate(a infer.Annotator) {
	a.Describe(&c.Enabled, "Whether prompt/response capture is enforced for requests passing through this guardrail.")
	a.Describe(&c.RedactPii, "Whether captured prompts have PII redacted.")
}

// AgentNetworkGuardrailChecks bundles the guardrail's check parameters.
type AgentNetworkGuardrailChecks struct {
	ModelAllowlist AgentNetworkGuardrailModelAllowlist `pulumi:"modelAllowlist"`
	PromptCapture  AgentNetworkGuardrailPromptCapture  `pulumi:"promptCapture"`
}

// Annotate provides documentation for AgentNetworkGuardrailChecks fields.
func (c *AgentNetworkGuardrailChecks) Annotate(a infer.Annotator) {
	a.Describe(&c.ModelAllowlist, "Restricts requests to an explicit set of catalog model IDs.")
	a.Describe(&c.PromptCapture, "Controls request/response prompt capture.")
}

// AgentNetworkGuardrailArgs defines input fields for an Agent Network guardrail.
type AgentNetworkGuardrailArgs struct {
	Name        string                      `pulumi:"name"`
	Description *string                     `pulumi:"description,optional"`
	Checks      AgentNetworkGuardrailChecks `pulumi:"checks"`
}

// Annotate provides documentation for AgentNetworkGuardrailArgs fields.
func (a *AgentNetworkGuardrailArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name for the guardrail.")
	ann.Describe(&a.Description, "Optional human-readable description.")
	ann.Describe(&a.Checks, "Guardrail check parameters. Each entry has an enabled flag plus per-check configuration; disabled entries are inert.")
}

// AgentNetworkGuardrailState represents the output state of an Agent Network guardrail.
type AgentNetworkGuardrailState struct {
	Name        string                      `pulumi:"name"`
	Description *string                     `pulumi:"description,optional"`
	Checks      AgentNetworkGuardrailChecks `pulumi:"checks"`
	CreatedAt   *string                     `pulumi:"createdAt,optional"`
	UpdatedAt   *string                     `pulumi:"updatedAt,optional"`
}

// Annotate provides documentation for AgentNetworkGuardrailState fields.
func (a *AgentNetworkGuardrailState) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Display name for the guardrail.")
	ann.Describe(&a.Description, "Optional human-readable description.")
	ann.Describe(&a.Checks, "Guardrail check parameters.")
	ann.Describe(&a.CreatedAt, "Timestamp when the guardrail was created.")
	ann.Describe(&a.UpdatedAt, "Timestamp when the guardrail was last updated.")
}

// Create creates a new Agent Network guardrail.
func (*AgentNetworkGuardrail) Create(
	ctx context.Context, req infer.CreateRequest[AgentNetworkGuardrailArgs],
) (infer.CreateResponse[AgentNetworkGuardrailState], error) {
	p.GetLogger(ctx).Debugf("Create:AgentNetworkGuardrail name=%s", req.Inputs.Name)

	if req.DryRun {
		return infer.CreateResponse[AgentNetworkGuardrailState]{
			ID:     "preview",
			Output: agentNetworkGuardrailStateFromArgs(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[AgentNetworkGuardrailState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// The API always echoes description on the wire and treats "" as unset, so an
	// omitted input travels as "" rather than being left out: that keeps inputs and
	// state in agreement and makes removing the field clear it server-side.
	description := strPtr(req.Inputs.Description)

	created, err := client.AgentNetwork.CreateGuardrail(ctx, nbapi.AgentNetworkGuardrailRequest{
		Name:        req.Inputs.Name,
		Description: &description,
		Checks:      toAPIAgentNetworkGuardrailChecks(req.Inputs.Checks),
	})
	if err != nil {
		return infer.CreateResponse[AgentNetworkGuardrailState]{}, fmt.Errorf("creating agent network guardrail failed: %w", err)
	}

	return infer.CreateResponse[AgentNetworkGuardrailState]{
		ID:     created.Id,
		Output: agentNetworkGuardrailStateFromAPI(*created),
	}, nil
}

// Read fetches the current state of an Agent Network guardrail from NetBird.
func (*AgentNetworkGuardrail) Read(
	ctx context.Context, req infer.ReadRequest[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState],
) (infer.ReadResponse[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState], error) {
	p.GetLogger(ctx).Debugf("Read:AgentNetworkGuardrail[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	guardrail, err := client.AgentNetwork.GetGuardrail(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState]{
				ID:     "",
				Inputs: AgentNetworkGuardrailArgs{},  //nolint:exhaustruct
				State:  AgentNetworkGuardrailState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState]{}, fmt.Errorf("reading agent network guardrail failed: %w", err)
	}

	state := agentNetworkGuardrailStateFromAPI(*guardrail)

	return infer.ReadResponse[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState]{
		ID: req.ID,
		Inputs: AgentNetworkGuardrailArgs{
			Name:        guardrail.Name,
			Description: &guardrail.Description,
			Checks:      state.Checks,
		},
		State: state,
	}, nil
}

// Update updates an Agent Network guardrail.
func (*AgentNetworkGuardrail) Update(
	ctx context.Context, req infer.UpdateRequest[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState],
) (infer.UpdateResponse[AgentNetworkGuardrailState], error) {
	p.GetLogger(ctx).Debugf("Update:AgentNetworkGuardrail[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[AgentNetworkGuardrailState]{
			Output: agentNetworkGuardrailStateFromArgs(req.Inputs, req.State.CreatedAt, req.State.UpdatedAt),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AgentNetworkGuardrailState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// The API always echoes description on the wire and treats "" as unset, so an
	// omitted input travels as "" rather than being left out: that keeps inputs and
	// state in agreement and makes removing the field clear it server-side.
	description := strPtr(req.Inputs.Description)

	updated, err := client.AgentNetwork.UpdateGuardrail(ctx, req.ID, nbapi.AgentNetworkGuardrailRequest{
		Name:        req.Inputs.Name,
		Description: &description,
		Checks:      toAPIAgentNetworkGuardrailChecks(req.Inputs.Checks),
	})
	if err != nil {
		return infer.UpdateResponse[AgentNetworkGuardrailState]{}, fmt.Errorf("updating agent network guardrail failed: %w", err)
	}

	return infer.UpdateResponse[AgentNetworkGuardrailState]{
		Output: agentNetworkGuardrailStateFromAPI(*updated),
	}, nil
}

// Delete removes an Agent Network guardrail from NetBird.
func (*AgentNetworkGuardrail) Delete(ctx context.Context, req infer.DeleteRequest[AgentNetworkGuardrailState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:AgentNetworkGuardrail[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.AgentNetwork.DeleteGuardrail(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting agent network guardrail failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*AgentNetworkGuardrail) Diff(
	ctx context.Context, req infer.DiffRequest[AgentNetworkGuardrailArgs, AgentNetworkGuardrailState],
) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:AgentNetworkGuardrail[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalOptionalStr(req.Inputs.Description, req.State.Description) {
		diff["description"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlice(req.Inputs.Checks.ModelAllowlist.Models, req.State.Checks.ModelAllowlist.Models) ||
		req.Inputs.Checks.ModelAllowlist.Enabled != req.State.Checks.ModelAllowlist.Enabled ||
		req.Inputs.Checks.PromptCapture != req.State.Checks.PromptCapture {
		diff["checks"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:AgentNetworkGuardrail[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation.
func (*AgentNetworkGuardrail) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[AgentNetworkGuardrailArgs], error) {
	p.GetLogger(ctx).Debugf("Check:AgentNetworkGuardrail old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[AgentNetworkGuardrailArgs](ctx, req.NewInputs)

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{Property: "name", Reason: "name must not be empty"})
	}

	for i, modelID := range args.Checks.ModelAllowlist.Models {
		if isBlank(modelID) {
			failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("checks.modelAllowlist.models[%d]", i), Reason: "model id must not be empty"})
		}
	}

	return infer.CheckResponse[AgentNetworkGuardrailArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*AgentNetworkGuardrail) WireDependencies(field infer.FieldSelector, args *AgentNetworkGuardrailArgs, state *AgentNetworkGuardrailState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Description).DependsOn(field.InputField(&args.Description))
	field.OutputField(&state.Checks).DependsOn(field.InputField(&args.Checks))
}

func agentNetworkGuardrailStateFromArgs(args AgentNetworkGuardrailArgs, createdAt, updatedAt *string) AgentNetworkGuardrailState {
	return AgentNetworkGuardrailState{
		Name:        args.Name,
		Description: args.Description,
		Checks:      args.Checks,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func agentNetworkGuardrailStateFromAPI(guardrail nbapi.AgentNetworkGuardrail) AgentNetworkGuardrailState {
	var createdAt, updatedAt *string

	if guardrail.CreatedAt != nil {
		formatted := guardrail.CreatedAt.Format(idpTimeFormat)
		createdAt = &formatted
	}

	if guardrail.UpdatedAt != nil {
		formatted := guardrail.UpdatedAt.Format(idpTimeFormat)
		updatedAt = &formatted
	}

	return AgentNetworkGuardrailState{
		Name:        guardrail.Name,
		Description: &guardrail.Description,
		Checks:      fromAPIAgentNetworkGuardrailChecks(guardrail.Checks),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func toAPIAgentNetworkGuardrailChecks(checks AgentNetworkGuardrailChecks) nbapi.AgentNetworkGuardrailChecks {
	var out nbapi.AgentNetworkGuardrailChecks

	out.ModelAllowlist.Enabled = checks.ModelAllowlist.Enabled
	out.ModelAllowlist.Models = checks.ModelAllowlist.Models
	out.PromptCapture.Enabled = checks.PromptCapture.Enabled
	out.PromptCapture.RedactPii = checks.PromptCapture.RedactPii

	return out
}

func fromAPIAgentNetworkGuardrailChecks(checks nbapi.AgentNetworkGuardrailChecks) AgentNetworkGuardrailChecks {
	return AgentNetworkGuardrailChecks{
		ModelAllowlist: AgentNetworkGuardrailModelAllowlist{
			Enabled: checks.ModelAllowlist.Enabled,
			Models:  checks.ModelAllowlist.Models,
		},
		PromptCapture: AgentNetworkGuardrailPromptCapture{
			Enabled:   checks.PromptCapture.Enabled,
			RedactPii: checks.PromptCapture.RedactPii,
		},
	}
}
