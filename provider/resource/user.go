package resource

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/management/server/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// User represents a NetBird user resource (invitation).
type User struct{}

// Annotate describes the resource and its fields.
func (u *User) Annotate(a infer.Annotator) {
	a.Describe(u, "A NetBird user that receives an invite and is optionally assigned groups and roles.")
}

// UserArgs defines input arguments for creating a NetBird user.
type UserArgs struct {
	// Email is the user's email address to which the invite is sent.
	Email *string `pulumi:"email"`
	// Name is the full name of the user.
	Name *string `pulumi:"name"`
	// Role is the NetBird account role to assign (e.g. 'admin', 'user').
	Role string `pulumi:"role"`
	// IsServiceUser is true if the user is a service identity.
	IsServiceUser bool `pulumi:"is_service_user"`
	// AutoGroups is the list of group IDs to auto-assign peers to.
	AutoGroups []string `pulumi:"auto_groups"`
	// IsBlocked indicates whether the user is blocked from accessing the system.
	// Used only on update, not create.
	IsBlocked *bool `pulumi:"blocked"`
}

// Annotate adds descriptions for SDK schema generation.
func (user *UserArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&user.Email, "Email address to send user invite to.")
	annotator.Describe(&user.Name, "Full name of the user.")
	annotator.Describe(&user.Role, "NetBird account role (e.g., 'admin', 'user').")
	annotator.Describe(&user.IsServiceUser, "Whether this user is a service identity.")
	annotator.Describe(&user.AutoGroups, "List of group IDs to auto-assign this user’s peers to.")
	annotator.Describe(&user.IsBlocked, "Indicates whether the user is blocked from accessing the system. Used only on update, not create.")
}

// UserState represents the stored state of a NetBird user in Pulumi.
type UserState struct {
	Email         *string  `pulumi:"email"`
	Name          *string  `pulumi:"name"`
	Role          string   `pulumi:"role"`
	IsServiceUser bool     `pulumi:"isServiceUser"`
	AutoGroups    []string `pulumi:"autoGroups"`
	IsBlocked     *bool    `pulumi:"blocked"`
}

// Annotate documents the stored state for the Pulumi schema.
func (user *UserState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&user.Email, "Email address of the user.")
	annotator.Describe(&user.Name, "Full name of the user.")
	annotator.Describe(&user.Role, "NetBird account role assigned to the user.")
	annotator.Describe(&user.IsServiceUser, "Whether this user is a service identity.")
	annotator.Describe(&user.AutoGroups, "Groups this user’s peers are automatically assigned to.")
	annotator.Describe(&user.IsBlocked, "Indicates whether the user is blocked from accessing the system")
}

// Create creates a new NetBird user.
func (*User) Create(ctx context.Context, req infer.CreateRequest[UserArgs]) (infer.CreateResponse[UserState], error) {
	p.GetLogger(ctx).Debugf("Create:User name=%s, email=%s", strPtr(req.Inputs.Name), strPtr(req.Inputs.Email))

	if req.DryRun {
		return infer.CreateResponse[UserState]{
			ID: "preview",
			Output: UserState{
				Name:          req.Inputs.Name,
				Email:         req.Inputs.Email,
				Role:          req.Inputs.Role,
				IsServiceUser: req.Inputs.IsServiceUser,
				AutoGroups:    req.Inputs.AutoGroups,
				IsBlocked:     nil,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[UserState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	user, err := client.Users.Create(ctx, nbapi.UserCreateRequest{
		Name:          req.Inputs.Name,
		Email:         req.Inputs.Email,
		Role:          req.Inputs.Role,
		IsServiceUser: req.Inputs.IsServiceUser,
		AutoGroups:    req.Inputs.AutoGroups,
	})
	if err != nil {
		return infer.CreateResponse[UserState]{}, fmt.Errorf("creating network failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:UserAPI name=%s, id=%s", user.Name, user.Id)

	return infer.CreateResponse[UserState]{
		ID: user.Id,
		Output: UserState{
			Name:          req.Inputs.Name,
			Email:         req.Inputs.Email,
			Role:          req.Inputs.Role,
			IsServiceUser: req.Inputs.IsServiceUser,
			AutoGroups:    req.Inputs.AutoGroups,
			IsBlocked:     nil, // IsBlocked is not set on create, only on update
		},
	}, nil
}

// Read fetches the current state of a user from NetBird.
func (*User) Read(ctx context.Context, req infer.ReadRequest[UserArgs, UserState]) (infer.ReadResponse[UserArgs, UserState], error) {
	p.GetLogger(ctx).Debugf("Read:User[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[UserArgs, UserState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	// Netbird Api does not implement user get by Id
	users, err := client.Users.List(ctx)
	if err != nil {
		return infer.ReadResponse[UserArgs, UserState]{}, fmt.Errorf("reading user failed: %w", err)
	}

	var foundUser *nbapi.User

	for _, u := range users {
		if u.Id == req.ID {
			foundUser = &u

			break
		}
	}

	if foundUser == nil {
		return infer.ReadResponse[UserArgs, UserState]{}, fmt.Errorf("user with ID %s not found", req.ID)
	}

	p.GetLogger(ctx).Debugf("Read:UserAPI[%s] name=%s, email=%s", foundUser.Id, foundUser.Name, foundUser.Email)

	return infer.ReadResponse[UserArgs, UserState]{
		ID: req.ID,
		Inputs: UserArgs{
			Name:          &foundUser.Name,
			Email:         &foundUser.Email,
			Role:          foundUser.Role,
			IsServiceUser: *foundUser.IsServiceUser,
			AutoGroups:    foundUser.AutoGroups,
			IsBlocked:     &foundUser.IsBlocked,
		},
		State: UserState{
			Name:          &foundUser.Name,
			Email:         &foundUser.Email,
			Role:          foundUser.Role,
			IsServiceUser: *foundUser.IsServiceUser,
			AutoGroups:    foundUser.AutoGroups,
			IsBlocked:     &foundUser.IsBlocked,
		},
	}, nil
}

// Update updates the state of the NetBird User resource if needed.
func (*User) Update(ctx context.Context, req infer.UpdateRequest[UserArgs, UserState]) (infer.UpdateResponse[UserState], error) {
	p.GetLogger(ctx).Debugf("Update:User[%s]", req.ID)

	// Check if isBlocked is nil, if so, we set it to false
	var isBlocked *bool
	if req.Inputs.IsBlocked != nil {
		isBlocked = req.Inputs.IsBlocked
	} else {
		f := false
		isBlocked = &f
	}

	if req.DryRun {
		return infer.UpdateResponse[UserState]{
			Output: UserState{
				Name:          req.Inputs.Name,
				Email:         req.Inputs.Email,
				Role:          req.Inputs.Role,
				IsServiceUser: req.Inputs.IsServiceUser,
				AutoGroups:    req.Inputs.AutoGroups,
				IsBlocked:     isBlocked,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[UserState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	_, err = client.Users.Update(ctx, req.ID, nbapi.UserRequest{
		Role:       req.Inputs.Role,
		AutoGroups: req.Inputs.AutoGroups,
		IsBlocked:  *isBlocked,
	})
	if err != nil {
		return infer.UpdateResponse[UserState]{}, fmt.Errorf("updating peer failed: %w", err)
	}

	return infer.UpdateResponse[UserState]{
		Output: UserState{
			Name:          req.Inputs.Name,
			Email:         req.Inputs.Email,
			Role:          req.Inputs.Role,
			IsServiceUser: req.Inputs.IsServiceUser,
			AutoGroups:    req.Inputs.AutoGroups,
			IsBlocked:     isBlocked,
		},
	}, nil
}

// Delete removes a user from NetBird.
func (*User) Delete(ctx context.Context, req infer.DeleteRequest[UserState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:User[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Users.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting user failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*User) Diff(ctx context.Context, req infer.DiffRequest[UserArgs, UserState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:User[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Email != req.State.Email {
		diff["email"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.Role != req.State.Role {
		diff["role"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.IsServiceUser != req.State.IsServiceUser {
		diff["is_service_user"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if equalSlice(req.Inputs.AutoGroups, req.State.AutoGroups) {
		diff["auto_groups"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	if req.Inputs.IsBlocked != req.State.IsBlocked {
		diff["blocked"] = p.PropertyDiff{
			InputDiff: false,
			Kind:      p.Update,
		}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation and default setting.
func (*User) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[UserArgs], error) {
	p.GetLogger(ctx).Debugf("Check:User old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())
	args, failures, err := infer.DefaultCheck[UserArgs](ctx, req.NewInputs)

	return infer.CheckResponse[UserArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*User) WireDependencies(f infer.FieldSelector, args *UserArgs, state *UserState) {
	f.OutputField(&state.Email).DependsOn(f.InputField(&args.Email))
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.Role).DependsOn(f.InputField(&args.Role))
	f.OutputField(&state.IsServiceUser).DependsOn(f.InputField(&args.IsServiceUser))
	f.OutputField(&state.AutoGroups).DependsOn(f.InputField(&args.AutoGroups))
	f.OutputField(&state.IsBlocked).DependsOn(f.InputField(&args.IsBlocked))
}
