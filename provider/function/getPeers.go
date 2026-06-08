package function

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// GetPeers lists all peers, optionally filtered to those belonging to a specific group.
type GetPeers struct{}

// Annotate describes the function.
func (f *GetPeers) Annotate(a infer.Annotator) {
	a.Describe(f, "List all NetBird peers, optionally filtered to those belonging to a specific group ID.")
}

// GetPeersArgs are the inputs for GetPeers.
type GetPeersArgs struct {
	GroupID *string `pulumi:"groupId,optional"`
}

// Annotate provides field descriptions for GetPeersArgs.
func (a *GetPeersArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.GroupID, "Optional group ID to filter peers. When set, only peers that belong to this group are returned.")
}

// PeerSummary is a brief summary of a NetBird peer.
type PeerSummary struct {
	ID        string   `pulumi:"id"`
	Name      string   `pulumi:"name"`
	IP        string   `pulumi:"ip"`
	DNSLabel  string   `pulumi:"dnsLabel"`
	Connected bool     `pulumi:"connected"`
	Hostname  string   `pulumi:"hostname"`
	Groups    []string `pulumi:"groups"`
}

// Annotate provides field descriptions for PeerSummary.
func (p *PeerSummary) Annotate(ann infer.Annotator) {
	ann.Describe(&p.ID, "The peer ID.")
	ann.Describe(&p.Name, "The peer name.")
	ann.Describe(&p.IP, "The WireGuard IP address assigned to the peer.")
	ann.Describe(&p.DNSLabel, "The DNS label used to form the peer's FQDN.")
	ann.Describe(&p.Connected, "Whether the peer is currently connected to the management server.")
	ann.Describe(&p.Hostname, "The OS hostname of the machine.")
	ann.Describe(&p.Groups, "IDs of groups the peer belongs to.")
}

// GetPeersResult is the output of GetPeers.
type GetPeersResult struct {
	Peers []PeerSummary `pulumi:"peers"`
}

// Annotate provides field descriptions for GetPeersResult.
func (r *GetPeersResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Peers, "The list of peers matching the filter criteria.")
}

// Invoke lists peers, applying an optional group filter.
func (f *GetPeers) Invoke(ctx context.Context, req infer.FunctionRequest[GetPeersArgs]) (infer.FunctionResponse[GetPeersResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[GetPeersResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiPeers, err := client.Peers.List(ctx)
	if err != nil {
		return infer.FunctionResponse[GetPeersResult]{}, fmt.Errorf("listing peers failed: %w", err)
	}

	peers := make([]PeerSummary, 0, len(apiPeers))

	for _, peer := range apiPeers {
		groups := make([]string, len(peer.Groups))
		for i, g := range peer.Groups {
			groups[i] = g.Id
		}

		if req.Input.GroupID != nil {
			matched := slices.Contains(groups, *req.Input.GroupID)

			if !matched {
				continue
			}
		}

		peers = append(peers, PeerSummary{
			ID:        peer.Id,
			Name:      peer.Name,
			IP:        peer.Ip,
			DNSLabel:  peer.DnsLabel,
			Connected: peer.Connected,
			Hostname:  peer.Hostname,
			Groups:    groups,
		})
	}

	return infer.FunctionResponse[GetPeersResult]{
		Output: GetPeersResult{
			Peers: peers,
		},
	}, nil
}
