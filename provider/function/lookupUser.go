package function

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// LookupUser looks up an existing NetBird user by email address.
type LookupUser struct{}

// Annotate describes the function.
func (f *LookupUser) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up an existing NetBird user by email and return their ID, role, and auto-group assignments.")
}

// LookupUserArgs are the inputs for LookupUser.
type LookupUserArgs struct {
	Email string `pulumi:"email"`
}

// Annotate provides field descriptions for LookupUserArgs.
func (a *LookupUserArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Email, "The email address of the user to look up.")
}

// LookupUserResult is the output of LookupUser.
type LookupUserResult struct {
	ID         string   `pulumi:"userId"`
	Name       string   `pulumi:"name"`
	Email      string   `pulumi:"email"`
	Role       string   `pulumi:"role"`
	IsBlocked  bool     `pulumi:"isBlocked"`
	AutoGroups []string `pulumi:"autoGroups"`
}

// Annotate provides field descriptions for LookupUserResult.
func (r *LookupUserResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.ID, "The NetBird user ID.")
	ann.Describe(&r.Name, "The user's display name.")
	ann.Describe(&r.Email, "The user's email address.")
	ann.Describe(&r.Role, "The user's role in the NetBird account (e.g. admin, user).")
	ann.Describe(&r.IsBlocked, "Whether the user is blocked from accessing the system.")
	ann.Describe(&r.AutoGroups, "Group IDs automatically assigned to peers registered by this user.")
}

// Invoke looks up a user by email.
func (f *LookupUser) Invoke(ctx context.Context, req infer.FunctionRequest[LookupUserArgs]) (infer.FunctionResponse[LookupUserResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupUserResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	users, err := client.Users.List(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupUserResult]{}, fmt.Errorf("listing users failed: %w", err)
	}

	for _, user := range users {
		if user.Email != req.Input.Email {
			continue
		}

		autoGroups := slices.Clone(user.AutoGroups)
		slices.Sort(autoGroups)

		return infer.FunctionResponse[LookupUserResult]{
			Output: LookupUserResult{
				ID:         user.Id,
				Name:       user.Name,
				Email:      user.Email,
				Role:       user.Role,
				IsBlocked:  user.IsBlocked,
				AutoGroups: autoGroups,
			},
		}, nil
	}

	return infer.FunctionResponse[LookupUserResult]{}, fmt.Errorf("user with email %q not found", req.Input.Email)
}
