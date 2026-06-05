package resource

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Route represents a Pulumi resource for NetBird routes.
type Route struct{}

// Annotate provides documentation for Route.
func (r *Route) Annotate(a infer.Annotator) {
	a.Describe(r, "A NetBird route resource for directing traffic through exit nodes or routing peers.")
}

// RouteArgs defines input fields for creating or updating a route.
type RouteArgs struct {
	NetworkID           string    `pulumi:"networkId"`
	Description         string    `pulumi:"description"`
	Enabled             bool      `pulumi:"enabled"`
	Masquerade          bool      `pulumi:"masquerade"`
	Metric              int       `pulumi:"metric"`
	KeepRoute           bool      `pulumi:"keepRoute"`
	Network             *string   `pulumi:"network,optional"`
	Domains             *[]string `pulumi:"domains,optional"`
	Groups              []string  `pulumi:"groups"`
	Peer                *string   `pulumi:"peer,optional"`
	PeerGroups          *[]string `pulumi:"peerGroups,optional"`
	AccessControlGroups *[]string `pulumi:"accessControlGroups,optional"`
	SkipAutoApply       *bool     `pulumi:"skipAutoApply,optional"`
}

// Annotate provides documentation for RouteArgs fields.
func (a *RouteArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&a.NetworkID, "Route network identifier, used to group HA routes.")
	annotator.Describe(&a.Description, "Route description.")
	annotator.Describe(&a.Enabled, "Whether the route is enabled.")
	annotator.Describe(&a.Masquerade, "Whether masquerading is enabled for this route.")
	annotator.Describe(&a.Metric, "Route metric; lower value means higher priority.")
	annotator.Describe(&a.KeepRoute, "Keep the route after a domain no longer resolves to the IP.")
	annotator.Describe(&a.Network, "Network CIDR range (conflicts with domains).")
	annotator.Describe(&a.Domains, "Domain list for dynamic resolution (conflicts with network).")
	annotator.Describe(&a.Groups, "Group IDs whose peers will use this route.")
	annotator.Describe(&a.Peer, "Peer ID acting as the routing peer (conflicts with peerGroups).")
	annotator.Describe(&a.PeerGroups, "Peer group IDs acting as routing peers (conflicts with peer).")
	annotator.Describe(&a.AccessControlGroups, "Access control group IDs associated with this route.")
	annotator.Describe(&a.SkipAutoApply, "Skip auto-application for exit-node (0.0.0.0/0) routes.")
}

// RouteState represents the output state of a route resource.
type RouteState struct {
	NetworkID           string    `pulumi:"networkId"`
	Description         string    `pulumi:"description"`
	Enabled             bool      `pulumi:"enabled"`
	Masquerade          bool      `pulumi:"masquerade"`
	Metric              int       `pulumi:"metric"`
	KeepRoute           bool      `pulumi:"keepRoute"`
	Network             *string   `pulumi:"network,optional"`
	Domains             *[]string `pulumi:"domains,optional"`
	Groups              []string  `pulumi:"groups"`
	Peer                *string   `pulumi:"peer,optional"`
	PeerGroups          *[]string `pulumi:"peerGroups,optional"`
	AccessControlGroups *[]string `pulumi:"accessControlGroups,optional"`
	SkipAutoApply       *bool     `pulumi:"skipAutoApply,optional"`
	NetworkType         *string   `pulumi:"networkType,optional"`
}

// Annotate provides documentation for RouteState fields.
func (s *RouteState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&s.NetworkID, "Route network identifier, used to group HA routes.")
	annotator.Describe(&s.Description, "Route description.")
	annotator.Describe(&s.Enabled, "Whether the route is enabled.")
	annotator.Describe(&s.Masquerade, "Whether masquerading is enabled for this route.")
	annotator.Describe(&s.Metric, "Route metric; lower value means higher priority.")
	annotator.Describe(&s.KeepRoute, "Keep the route after a domain no longer resolves to the IP.")
	annotator.Describe(&s.Network, "Network CIDR range.")
	annotator.Describe(&s.Domains, "Domain list for dynamic resolution.")
	annotator.Describe(&s.Groups, "Group IDs whose peers will use this route.")
	annotator.Describe(&s.Peer, "Peer ID acting as the routing peer.")
	annotator.Describe(&s.PeerGroups, "Peer group IDs acting as routing peers.")
	annotator.Describe(&s.AccessControlGroups, "Access control group IDs associated with this route.")
	annotator.Describe(&s.SkipAutoApply, "Skip auto-application for exit-node (0.0.0.0/0) routes.")
	annotator.Describe(&s.NetworkType, "Network type (IPv4, IPv6, or domain) — computed by the API.")
}

func routeStateFromAPI(route *nbapi.Route) RouteState {
	groups := slices.Clone(route.Groups)
	slices.Sort(groups)

	peerGroups := route.PeerGroups
	if peerGroups != nil {
		sorted := slices.Clone(*peerGroups)
		slices.Sort(sorted)
		peerGroups = &sorted
	}

	acGroups := route.AccessControlGroups
	if acGroups != nil {
		sorted := slices.Clone(*acGroups)
		slices.Sort(sorted)
		acGroups = &sorted
	}

	// Normalize empty string to nil so that routes using peerGroups (not peer)
	// don't produce a spurious diff against inputs that omit the peer field.
	peer := route.Peer
	if peer != nil && *peer == "" {
		peer = nil
	}

	var networkType *string

	if route.NetworkType != "" {
		nt := route.NetworkType
		networkType = &nt
	}

	return RouteState{
		NetworkID:           route.NetworkId,
		Description:         route.Description,
		Enabled:             route.Enabled,
		Masquerade:          route.Masquerade,
		Metric:              route.Metric,
		KeepRoute:           route.KeepRoute,
		Network:             route.Network,
		Domains:             route.Domains,
		Groups:              groups,
		Peer:                peer,
		PeerGroups:          peerGroups,
		AccessControlGroups: acGroups,
		SkipAutoApply:       route.SkipAutoApply,
		NetworkType:         networkType,
	}
}

func routeArgsFromState(state RouteState) RouteArgs {
	return RouteArgs{
		NetworkID:           state.NetworkID,
		Description:         state.Description,
		Enabled:             state.Enabled,
		Masquerade:          state.Masquerade,
		Metric:              state.Metric,
		KeepRoute:           state.KeepRoute,
		Network:             state.Network,
		Domains:             state.Domains,
		Groups:              state.Groups,
		Peer:                state.Peer,
		PeerGroups:          state.PeerGroups,
		AccessControlGroups: state.AccessControlGroups,
		SkipAutoApply:       state.SkipAutoApply,
	}
}

func routeRequest(args RouteArgs) nbapi.RouteRequest {
	return nbapi.RouteRequest{
		NetworkId:           args.NetworkID,
		Description:         args.Description,
		Enabled:             args.Enabled,
		Masquerade:          args.Masquerade,
		Metric:              args.Metric,
		KeepRoute:           args.KeepRoute,
		Network:             args.Network,
		Domains:             args.Domains,
		Groups:              args.Groups,
		Peer:                args.Peer,
		PeerGroups:          args.PeerGroups,
		AccessControlGroups: args.AccessControlGroups,
		SkipAutoApply:       args.SkipAutoApply,
	}
}

// Create creates a new NetBird route.
func (*Route) Create(ctx context.Context, req infer.CreateRequest[RouteArgs]) (infer.CreateResponse[RouteState], error) {
	p.GetLogger(ctx).Debugf("Create:Route networkId=%s", req.Inputs.NetworkID)

	if req.DryRun {
		return infer.CreateResponse[RouteState]{
			ID: "preview",
			Output: routeStateFromAPI(&nbapi.Route{
				NetworkId:           req.Inputs.NetworkID,
				Id:                  "preview",
				Description:         req.Inputs.Description,
				Enabled:             req.Inputs.Enabled,
				Masquerade:          req.Inputs.Masquerade,
				Metric:              req.Inputs.Metric,
				KeepRoute:           req.Inputs.KeepRoute,
				NetworkType:         "",
				Network:             req.Inputs.Network,
				Domains:             req.Inputs.Domains,
				Groups:              req.Inputs.Groups,
				Peer:                req.Inputs.Peer,
				PeerGroups:          req.Inputs.PeerGroups,
				AccessControlGroups: req.Inputs.AccessControlGroups,
				SkipAutoApply:       req.Inputs.SkipAutoApply,
			}),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[RouteState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	route, err := client.Routes.Create(ctx, routeRequest(req.Inputs))
	if err != nil {
		return infer.CreateResponse[RouteState]{}, fmt.Errorf("creating route failed: %w", err)
	}

	p.GetLogger(ctx).Debugf("Create:RouteAPI id=%s", route.Id)

	return infer.CreateResponse[RouteState]{
		ID:     route.Id,
		Output: routeStateFromAPI(route),
	}, nil
}

// Read fetches the current state of a route from NetBird.
func (*Route) Read(ctx context.Context, req infer.ReadRequest[RouteArgs, RouteState]) (infer.ReadResponse[RouteArgs, RouteState], error) {
	p.GetLogger(ctx).Debugf("Read:Route[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[RouteArgs, RouteState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	route, err := client.Routes.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[RouteArgs, RouteState]{
				ID:     "",
				Inputs: RouteArgs{},  //nolint:exhaustruct
				State:  RouteState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[RouteArgs, RouteState]{}, fmt.Errorf("reading route failed: %w", err)
	}

	state := routeStateFromAPI(route)

	return infer.ReadResponse[RouteArgs, RouteState]{
		ID:     route.Id,
		Inputs: routeArgsFromState(state),
		State:  state,
	}, nil
}

// Update updates an existing NetBird route.
func (*Route) Update(ctx context.Context, req infer.UpdateRequest[RouteArgs, RouteState]) (infer.UpdateResponse[RouteState], error) {
	p.GetLogger(ctx).Debugf("Update:Route[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[RouteState]{
			Output: routeStateFromAPI(&nbapi.Route{
				NetworkId:           req.Inputs.NetworkID,
				Id:                  req.ID,
				Description:         req.Inputs.Description,
				Enabled:             req.Inputs.Enabled,
				Masquerade:          req.Inputs.Masquerade,
				Metric:              req.Inputs.Metric,
				KeepRoute:           req.Inputs.KeepRoute,
				NetworkType:         "",
				Network:             req.Inputs.Network,
				Domains:             req.Inputs.Domains,
				Groups:              req.Inputs.Groups,
				Peer:                req.Inputs.Peer,
				PeerGroups:          req.Inputs.PeerGroups,
				AccessControlGroups: req.Inputs.AccessControlGroups,
				SkipAutoApply:       req.Inputs.SkipAutoApply,
			}),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[RouteState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	route, err := client.Routes.Update(ctx, req.ID, routeRequest(req.Inputs))
	if err != nil {
		return infer.UpdateResponse[RouteState]{}, fmt.Errorf("updating route failed: %w", err)
	}

	return infer.UpdateResponse[RouteState]{
		Output: routeStateFromAPI(route),
	}, nil
}

// Delete removes a route from NetBird.
func (*Route) Delete(ctx context.Context, req infer.DeleteRequest[RouteState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:Route[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.Routes.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting route failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between inputs and prior state.
func (*Route) Diff(ctx context.Context, req infer.DiffRequest[RouteArgs, RouteState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:Route[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.NetworkID != req.State.NetworkID {
		diff["networkId"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.Description != req.State.Description {
		diff["description"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Masquerade != req.State.Masquerade {
		diff["masquerade"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Metric != req.State.Metric {
		diff["metric"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.KeepRoute != req.State.KeepRoute {
		diff["keepRoute"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.Network, req.State.Network) {
		diff["network"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.Domains, req.State.Domains) {
		diff["domains"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlice(req.Inputs.Groups, req.State.Groups) {
		diff["groups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.Peer, req.State.Peer) {
		diff["peer"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.PeerGroups, req.State.PeerGroups) {
		diff["peerGroups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.AccessControlGroups, req.State.AccessControlGroups) {
		diff["accessControlGroups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.SkipAutoApply, req.State.SkipAutoApply) {
		diff["skipAutoApply"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:Route[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check provides input validation.
func (*Route) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[RouteArgs], error) {
	p.GetLogger(ctx).Debugf("Check:Route old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[RouteArgs](ctx, req.NewInputs)
	failures = routeCheckArgs(args, failures)

	return infer.CheckResponse[RouteArgs]{Inputs: args, Failures: failures}, err
}

func routeCheckArgs(args RouteArgs, failures []p.CheckFailure) []p.CheckFailure { //nolint:cyclop,gocognit
	if isBlank(args.NetworkID) {
		failures = append(failures, p.CheckFailure{Property: "networkId", Reason: "networkId must not be empty"})
	}

	if args.Network != nil && isBlank(*args.Network) {
		failures = append(failures, p.CheckFailure{Property: "network", Reason: "network must not be blank when provided"})
	}

	if args.Network == nil && (args.Domains == nil || len(*args.Domains) == 0) {
		failures = append(failures, p.CheckFailure{Property: "network", Reason: "either network or domains must be provided"})
	}

	if args.Network != nil && args.Domains != nil && len(*args.Domains) > 0 {
		failures = append(failures, p.CheckFailure{Property: "network", Reason: "network and domains are mutually exclusive"})
	}

	if args.Domains != nil {
		for i, d := range *args.Domains {
			if isBlank(d) {
				failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("domains[%d]", i), Reason: "domain must not be empty"})
			}
		}
	}

	for i, g := range args.Groups {
		if isBlank(g) {
			failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("groups[%d]", i), Reason: "group id must not be empty"})
		}
	}

	hasPeer := args.Peer != nil && !isBlank(*args.Peer)
	hasPeerGroups := args.PeerGroups != nil && len(*args.PeerGroups) > 0

	if !hasPeer && !hasPeerGroups {
		failures = append(failures, p.CheckFailure{Property: "peer", Reason: "either peer or peerGroups must be provided"})
	}

	if args.PeerGroups != nil {
		for i, pg := range *args.PeerGroups {
			if isBlank(pg) {
				failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("peerGroups[%d]", i), Reason: "peer group id must not be empty"})
			}
		}
	}

	if args.AccessControlGroups != nil {
		for i, acg := range *args.AccessControlGroups {
			if isBlank(acg) {
				failures = append(failures, p.CheckFailure{Property: fmt.Sprintf("accessControlGroups[%d]", i), Reason: "access control group id must not be empty"})
			}
		}
	}

	return failures
}

// WireDependencies defines input/output field relationships.
func (*Route) WireDependencies(selector infer.FieldSelector, args *RouteArgs, state *RouteState) {
	selector.OutputField(&state.NetworkID).DependsOn(selector.InputField(&args.NetworkID))
	selector.OutputField(&state.Description).DependsOn(selector.InputField(&args.Description))
	selector.OutputField(&state.Enabled).DependsOn(selector.InputField(&args.Enabled))
	selector.OutputField(&state.Masquerade).DependsOn(selector.InputField(&args.Masquerade))
	selector.OutputField(&state.Metric).DependsOn(selector.InputField(&args.Metric))
	selector.OutputField(&state.KeepRoute).DependsOn(selector.InputField(&args.KeepRoute))
	selector.OutputField(&state.Network).DependsOn(selector.InputField(&args.Network))
	selector.OutputField(&state.Domains).DependsOn(selector.InputField(&args.Domains))
	selector.OutputField(&state.Groups).DependsOn(selector.InputField(&args.Groups))
	selector.OutputField(&state.Peer).DependsOn(selector.InputField(&args.Peer))
	selector.OutputField(&state.PeerGroups).DependsOn(selector.InputField(&args.PeerGroups))
	selector.OutputField(&state.AccessControlGroups).DependsOn(selector.InputField(&args.AccessControlGroups))
	selector.OutputField(&state.SkipAutoApply).DependsOn(selector.InputField(&args.SkipAutoApply))
}
