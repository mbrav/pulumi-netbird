package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

const tokenTimeFormat = "2006-01-02T15:04:05Z07:00"

// Token represents a NetBird personal access token (PAT) issued for a user.
type Token struct{}

// Annotate adds a description to the Token resource type.
func (t *Token) Annotate(a infer.Annotator) {
	a.Describe(&t, "A NetBird personal access token (PAT) for a user. The plaintext token is "+
		"returned only once, on creation, and is exposed as a secret output. The token cannot be "+
		"modified after creation; any input change forces a replacement.")
}

// TokenArgs defines input fields for creating a personal access token.
type TokenArgs struct {
	UserID    string `pulumi:"userId"`
	Name      string `pulumi:"name"`
	ExpiresIn int    `pulumi:"expiresIn"`
}

// Annotate provides documentation for TokenArgs fields.
func (t *TokenArgs) Annotate(a infer.Annotator) {
	a.Describe(&t.UserID, "ID of the user the token is issued for.")
	a.Describe(&t.Name, "Display name of the token.")
	a.Describe(&t.ExpiresIn, "Token lifetime in days.")
}

// TokenState represents the output state of a personal access token resource.
type TokenState struct {
	UserID         string  `pulumi:"userId"`
	Name           string  `pulumi:"name"`
	ExpiresIn      int     `pulumi:"expiresIn"`
	Token          *string `provider:"secret"                pulumi:"token,optional"`
	CreatedAt      *string `pulumi:"createdAt,optional"`
	CreatedBy      *string `pulumi:"createdBy,optional"`
	ExpirationDate *string `pulumi:"expirationDate,optional"`
	LastUsed       *string `pulumi:"lastUsed,optional"`
}

// Annotate provides documentation for TokenState fields.
func (t *TokenState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&t.UserID, "ID of the user the token is issued for.")
	annotator.Describe(&t.Name, "Display name of the token.")
	annotator.Describe(&t.ExpiresIn, "Token lifetime in days.")
	annotator.Describe(&t.Token, "Plaintext token value. Only populated on creation; never returned by the API afterwards.")
	annotator.Describe(&t.CreatedAt, "Timestamp the token was created.")
	annotator.Describe(&t.CreatedBy, "User ID of the principal that created the token.")
	annotator.Describe(&t.ExpirationDate, "Timestamp the token expires.")
	annotator.Describe(&t.LastUsed, "Timestamp the token was last used, if ever.")
}

// Create issues a new personal access token for a user.
func (*Token) Create(ctx context.Context, req infer.CreateRequest[TokenArgs]) (infer.CreateResponse[TokenState], error) {
	p.GetLogger(ctx).Debugf("Create:Token userId=%s name=%s", req.Inputs.UserID, req.Inputs.Name)

	if req.DryRun {
		return infer.CreateResponse[TokenState]{
			ID: "preview",
			Output: TokenState{
				UserID:         req.Inputs.UserID,
				Name:           req.Inputs.Name,
				ExpiresIn:      req.Inputs.ExpiresIn,
				Token:          nil,
				CreatedAt:      nil,
				CreatedBy:      nil,
				ExpirationDate: nil,
				LastUsed:       nil,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[TokenState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	generated, err := client.Tokens.Create(ctx, req.Inputs.UserID, nbapi.PersonalAccessTokenRequest{
		Name:      req.Inputs.Name,
		ExpiresIn: req.Inputs.ExpiresIn,
	})
	if err != nil {
		return infer.CreateResponse[TokenState]{}, fmt.Errorf("creating token failed: %w", err)
	}

	pat := generated.PersonalAccessToken
	plain := generated.PlainToken

	p.GetLogger(ctx).Debugf("Create:TokenAPI id=%s name=%s", pat.Id, pat.Name)

	state := tokenStateFromAPI(req.Inputs.UserID, req.Inputs.ExpiresIn, pat)
	state.Token = &plain

	return infer.CreateResponse[TokenState]{
		ID:     pat.Id,
		Output: state,
	}, nil
}

// Read fetches the current state of a personal access token from NetBird.
// Supports the compound import ID "userID/tokenID".
func (*Token) Read(ctx context.Context, req infer.ReadRequest[TokenArgs, TokenState]) (infer.ReadResponse[TokenArgs, TokenState], error) {
	p.GetLogger(ctx).Debugf("Read:Token[%s]", req.ID)

	userID := req.State.UserID
	tokenID := req.ID

	if userID == "" {
		var parseErr error

		userID, tokenID, parseErr = parseNestedID("Token", req.ID)
		if parseErr != nil {
			return infer.ReadResponse[TokenArgs, TokenState]{}, parseErr
		}
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[TokenArgs, TokenState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	pat, err := client.Tokens.Get(ctx, userID, tokenID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[TokenArgs, TokenState]{
				ID:     "",
				Inputs: TokenArgs{},  //nolint:exhaustruct
				State:  TokenState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[TokenArgs, TokenState]{}, fmt.Errorf("reading token failed: %w", err)
	}

	state := tokenStateFromAPI(userID, req.State.ExpiresIn, *pat)
	// The plaintext token is only ever returned on creation; preserve any prior value.
	state.Token = req.State.Token

	return infer.ReadResponse[TokenArgs, TokenState]{
		ID: tokenID,
		Inputs: TokenArgs{
			UserID:    userID,
			Name:      pat.Name,
			ExpiresIn: req.Inputs.ExpiresIn,
		},
		State: state,
	}, nil
}

// Delete removes a personal access token from NetBird.
func (*Token) Delete(ctx context.Context, req infer.DeleteRequest[TokenState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:Token[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Tokens.Delete(ctx, req.State.UserID, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting token failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state. The API has no update
// endpoint for tokens, so any input change forces a replacement.
func (*Token) Diff(ctx context.Context, req infer.DiffRequest[TokenArgs, TokenState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Token[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.UserID != req.State.UserID {
		diff["userId"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.ExpiresIn != req.State.ExpiresIn {
		diff["expiresIn"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a personal access token.
func (*Token) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[TokenArgs], error) {
	p.GetLogger(ctx).Debugf("Check:Token old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[TokenArgs](ctx, req.NewInputs)

	if isBlank(args.UserID) {
		failures = append(failures, p.CheckFailure{
			Property: "userId",
			Reason:   "userId must not be empty",
		})
	}

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{
			Property: "name",
			Reason:   "name must not be empty",
		})
	}

	if args.ExpiresIn <= 0 {
		failures = append(failures, p.CheckFailure{
			Property: "expiresIn",
			Reason:   "expiresIn must be greater than 0",
		})
	}

	return infer.CheckResponse[TokenArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*Token) WireDependencies(field infer.FieldSelector, args *TokenArgs, state *TokenState) {
	field.OutputField(&state.UserID).DependsOn(field.InputField(&args.UserID))
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.ExpiresIn).DependsOn(field.InputField(&args.ExpiresIn))
}

// tokenStateFromAPI maps an API token into resource state (excluding the plaintext token).
func tokenStateFromAPI(userID string, expiresIn int, pat nbapi.PersonalAccessToken) TokenState {
	createdAt := pat.CreatedAt.Format(tokenTimeFormat)
	createdBy := pat.CreatedBy
	expirationDate := pat.ExpirationDate.Format(tokenTimeFormat)

	var lastUsed *string

	if pat.LastUsed != nil {
		formatted := pat.LastUsed.Format(tokenTimeFormat)
		lastUsed = &formatted
	}

	return TokenState{
		UserID:         userID,
		Name:           pat.Name,
		ExpiresIn:      expiresIn,
		Token:          nil,
		CreatedAt:      &createdAt,
		CreatedBy:      &createdBy,
		ExpirationDate: &expirationDate,
		LastUsed:       lastUsed,
	}
}
