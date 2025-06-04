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
}

// Annotate adds descriptions for SDK schema generation.
func (p *UserArgs) Annotate(a infer.Annotator) {
	a.Describe(&p.Email, "Email address to send user invite to.")
	a.Describe(&p.Name, "Full name of the user.")
	a.Describe(&p.Role, "NetBird account role (e.g., 'admin', 'user').")
	a.Describe(&p.IsServiceUser, "Whether this user is a service identity.")
	a.Describe(&p.AutoGroups, "List of group IDs to auto-assign this user’s peers to.")
}

// UserState represents the stored state of a NetBird user in Pulumi.
type UserState struct {
	Email         *string  `pulumi:"email"`
	Name          *string  `pulumi:"name"`
	Role          string   `pulumi:"role"`
	IsServiceUser bool     `pulumi:"isServiceUser"`
	AutoGroups    []string `pulumi:"autoGroups"`
}

// Annotate documents the stored state for the Pulumi schema.
func (p *UserState) Annotate(a infer.Annotator) {
	a.Describe(&p.Email, "Email address of the user.")
	a.Describe(&p.Name, "Full name of the user.")
	a.Describe(&p.Role, "NetBird account role assigned to the user.")
	a.Describe(&p.IsServiceUser, "Whether this user is a service identity.")
	a.Describe(&p.AutoGroups, "Groups this user’s peers are automatically assigned to.")
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
		},
	}, nil
}

// // Read fetches the current state of a user from NetBird.
// func (*User) Read(ctx context.Context, req infer.ReadRequest[UserArgs, UserState]) (infer.ReadResponse[UserArgs, UserState], error) {
// 	p.GetLogger(ctx).Debugf("Read:User[%s]", req.ID)
//
// 	client, err := config.GetNetBirdClient(ctx)
// 	if err != nil {
// 		return infer.ReadResponse[UserArgs, UserState]{}, fmt.Errorf("error getting NetBird client: %w", err)
// 	}
//
// 	// Netbird Api does not implement user get by Id
// 	user, err := client.Users.List(ctx)
// 	if err != nil {
// 		return infer.ReadResponse[UserArgs, UserState]{}, fmt.Errorf("reading user failed: %w", err)
// 	}
//
// 	p.GetLogger(ctx).Debugf("Read:UserAPI[%s] name=%s, email=%s", user.Id, user.Name, deref(user.Email))
//
// 	return infer.ReadResponse[UserArgs, UserState]{
// 		ID: req.ID,
// 		Inputs: UserArgs{
// 			Name:          user.Name,
// 			Email:         user.Email,
// 			Role:          user.Role,
// 			IsServiceUser: user.IsServiceUser,
// 			AutoGroups:    user.AutoGroups,
// 		},
// 		State: UserState{
// 			Name:          user.Name,
// 			Email:         user.Email,
// 			Role:          user.Role,
// 			IsServiceUser: user.IsServiceUser,
// 			AutoGroups:    user.AutoGroups,
// 		},
// 	}, nil
// }

//
// // Update updates the state of the NetBird User if needed.
// func (*User) Update(ctx context.Context, req infer.UpdateRequest[UserArgs, UserState]) (infer.UpdateResponse[UserState], error) {
// 	p.GetLogger(ctx).Debugf("Update:User[%s]", req.ID)
//
// 	if req.DryRun {
// 		return infer.UpdateResponse[UserState]{
// 			Output: UserState{
// 				Name:                        req.Inputs.Name,
// 				InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
// 				LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
// 				SSHEnabled:                  req.Inputs.SSHEnabled,
// 				ApprovalRequired:            nil,
// 			},
// 		}, nil
// 	}
//
// 	client, err := config.GetNetBirdClient(ctx)
// 	if err != nil {
// 		return infer.UpdateResponse[UserState]{}, fmt.Errorf("error getting NetBird client: %w", err)
// 	}
//
// 	_, err = client.Users.Update(ctx, req.ID, nbapi.UserRequest{
// 		Name:                        req.Inputs.Name,
// 		InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
// 		LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
// 		SshEnabled:                  req.Inputs.SSHEnabled,
// 		ApprovalRequired:            nil, // ApprovalRequired is not supported in for Cloud version only
// 	})
// 	if err != nil {
// 		return infer.UpdateResponse[UserState]{}, fmt.Errorf("updating peer failed: %w", err)
// 	}
//
// 	return infer.UpdateResponse[UserState]{
// 		Output: UserState{
// 			Name:                        req.Inputs.Name,
// 			InactivityExpirationEnabled: req.Inputs.InactivityExpirationEnabled,
// 			LoginExpirationEnabled:      req.Inputs.LoginExpirationEnabled,
// 			SSHEnabled:                  req.Inputs.SSHEnabled,
// 			ApprovalRequired:            nil,
// 		},
// 	}, nil
// }
//
// // Diff detects changes between inputs and prior state.
// func (*User) Diff(ctx context.Context, req infer.DiffRequest[UserArgs, UserState]) (infer.DiffResponse, error) {
// 	p.GetLogger(ctx).Debugf("Diff:User[%s]", req.ID)
//
// 	diff := map[string]p.PropertyDiff{}
//
// 	if req.Inputs.Name != req.State.Name {
// 		diff["name"] = p.PropertyDiff{
// 			InputDiff: false,
// 			Kind:      p.Update,
// 		}
// 	}
//
// 	// if *req.Inputs.ApprovalRequired != *req.State.ApprovalRequired {
// 	// 	diff["approvalRequired"] = p.PropertyDiff{
// 	// 		InputDiff: false,
// 	// 		Kind:      p.Update,
// 	// 	}
// 	// }
//
// 	if req.Inputs.InactivityExpirationEnabled != req.State.InactivityExpirationEnabled {
// 		diff["inactivityExpirationEnabled"] = p.PropertyDiff{
// 			InputDiff: false,
// 			Kind:      p.Update,
// 		}
// 	}
//
// 	if req.Inputs.LoginExpirationEnabled != req.State.LoginExpirationEnabled {
// 		diff["loginExpirationEnabled"] = p.PropertyDiff{
// 			InputDiff: false,
// 			Kind:      p.Update,
// 		}
// 	}
//
// 	if req.Inputs.SSHEnabled != req.State.SSHEnabled {
// 		diff["sshEnabled"] = p.PropertyDiff{
// 			InputDiff: false,
// 			Kind:      p.Update,
// 		}
// 	}
//
// 	return infer.DiffResponse{
// 		DeleteBeforeReplace: false,
// 		HasChanges:          len(diff) > 0,
// 		DetailedDiff:        diff,
// 	}, nil
// }
