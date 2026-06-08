package function

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// LookupPeer looks up an existing NetBird peer by hostname/name.
type LookupPeer struct{}

// Annotate describes the function.
func (f *LookupPeer) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up an existing NetBird peer by name and return its ID, IP address, and group memberships.")
}

// LookupPeerArgs are the inputs for LookupPeer.
type LookupPeerArgs struct {
	Name string `pulumi:"name"`
}

// Annotate provides field descriptions for LookupPeerArgs.
func (a *LookupPeerArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "The name (hostname) of the peer to look up.")
}

// LookupPeerResult is the output of LookupPeer.
type LookupPeerResult struct {
	ID        string   `pulumi:"peerId"`
	Name      string   `pulumi:"name"`
	IP        string   `pulumi:"ip"`
	DNSLabel  string   `pulumi:"dnsLabel"`
	Connected bool     `pulumi:"connected"`
	Hostname  string   `pulumi:"hostname"`
	OS        string   `pulumi:"os"`
	Groups    []string `pulumi:"groups"`
}

// Annotate provides field descriptions for LookupPeerResult.
func (r *LookupPeerResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.ID, "The NetBird peer ID.")
	ann.Describe(&r.Name, "The peer name.")
	ann.Describe(&r.IP, "The WireGuard IP address assigned to the peer.")
	ann.Describe(&r.DNSLabel, "The DNS label used to form the peer's FQDN.")
	ann.Describe(&r.Connected, "Whether the peer is currently connected to the management server.")
	ann.Describe(&r.Hostname, "The OS hostname of the machine.")
	ann.Describe(&r.OS, "Operating system string reported by the peer.")
	ann.Describe(&r.Groups, "IDs of groups the peer belongs to.")
}

// Invoke looks up a peer by name.
func (f *LookupPeer) Invoke(ctx context.Context, req infer.FunctionRequest[LookupPeerArgs]) (infer.FunctionResponse[LookupPeerResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupPeerResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	peers, err := client.Peers.List(ctx)
	if err != nil {
		return infer.FunctionResponse[LookupPeerResult]{}, fmt.Errorf("listing peers failed: %w", err)
	}

	for _, peer := range peers {
		if peer.Name != req.Input.Name {
			continue
		}

		groups := make([]string, len(peer.Groups))
		for i, g := range peer.Groups {
			groups[i] = g.Id
		}

		return infer.FunctionResponse[LookupPeerResult]{
			Output: LookupPeerResult{
				ID:        peer.Id,
				Name:      peer.Name,
				IP:        peer.Ip,
				DNSLabel:  peer.DnsLabel,
				Connected: peer.Connected,
				Hostname:  peer.Hostname,
				OS:        peer.Os,
				Groups:    groups,
			},
		}, nil
	}

	return infer.FunctionResponse[LookupPeerResult]{}, fmt.Errorf("peer %q not found", req.Input.Name)
}
