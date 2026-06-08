package function

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// LookupRoute looks up an existing NetBird route by network CIDR.
type LookupRoute struct{}

// Annotate describes the function.
func (f *LookupRoute) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up an existing NetBird route by network CIDR and return its ID, routing peers, and configuration.")
}

// LookupRouteArgs are the inputs for LookupRoute.
type LookupRouteArgs struct {
	Network string `pulumi:"network"`
}

// Annotate provides field descriptions for LookupRouteArgs.
func (a *LookupRouteArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Network, "The network CIDR (e.g. '10.0.0.0/8') to look up. Returns the first matching route.")
}

// LookupRouteResult is the output of LookupRoute.
type LookupRouteResult struct {
	ID          string    `pulumi:"routeId"`
	Description string    `pulumi:"description"`
	Network     string    `pulumi:"network"`
	Domains     []string  `pulumi:"domains"`
	Enabled     bool      `pulumi:"enabled"`
	Masquerade  bool      `pulumi:"masquerade"`
	Metric      int       `pulumi:"metric"`
	Peer        *string   `pulumi:"peer,optional"`
	PeerGroups  []string  `pulumi:"peerGroups"`
	Groups      []string  `pulumi:"groups"`
}

// Annotate provides field descriptions for LookupRouteResult.
func (r *LookupRouteResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.ID, "The NetBird route ID.")
	ann.Describe(&r.Description, "The route description.")
	ann.Describe(&r.Network, "The network CIDR range.")
	ann.Describe(&r.Domains, "Domain list for dynamic resolution (if this is a domain route).")
	ann.Describe(&r.Enabled, "Whether the route is enabled.")
	ann.Describe(&r.Masquerade, "Whether the routing peer masquerades traffic to this prefix.")
	ann.Describe(&r.Metric, "Route metric; lower number means higher priority.")
	ann.Describe(&r.Peer, "The peer ID acting as the routing peer (mutually exclusive with peerGroups).")
	ann.Describe(&r.PeerGroups, "Group IDs whose peers act as routing peers (mutually exclusive with peer).")
	ann.Describe(&r.Groups, "Group IDs that have access to this route.")
}

// Invoke looks up a route by network CIDR.
func (f *LookupRoute) Invoke(ctx context.Context, req infer.FunctionRequest[LookupRouteArgs]) (infer.FunctionResponse[LookupRouteResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupRouteResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	routes, err := client.Routes.List(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupRouteResult]{}, fmt.Errorf("listing routes failed: %w", err)
	}

	for _, route := range routes {
		if route.Network == nil || *route.Network != req.Input.Network {
			continue
		}

		domains := []string{}
		if route.Domains != nil {
			domains = *route.Domains
		}

		peerGroups := []string{}
		if route.PeerGroups != nil {
			peerGroups = *route.PeerGroups
		}

		return infer.FunctionResponse[LookupRouteResult]{
			Output: LookupRouteResult{
				ID:          route.Id,
				Description: route.Description,
				Network:     *route.Network,
				Domains:     domains,
				Enabled:     route.Enabled,
				Masquerade:  route.Masquerade,
				Metric:      route.Metric,
				Peer:        route.Peer,
				PeerGroups:  peerGroups,
				Groups:      route.Groups,
			},
		}, nil
	}

	return infer.FunctionResponse[LookupRouteResult]{}, fmt.Errorf("route with network %q not found", req.Input.Network)
}
