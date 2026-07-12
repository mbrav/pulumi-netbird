package function

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// GetReverseProxyClusters lists all available reverse proxy clusters, optionally
// filtered by type ("account" or "shared").
type GetReverseProxyClusters struct{}

// Annotate describes the function.
func (f *GetReverseProxyClusters) Annotate(a infer.Annotator) {
	a.Describe(f, "List all NetBird reverse proxy clusters. Clusters are auto-provisioned "+
		"server-side (there is no create endpoint); use this to discover a cluster address "+
		"to pin a ReverseProxyDomain against via its targetCluster input.")
}

// GetReverseProxyClustersArgs are the inputs for GetReverseProxyClusters.
type GetReverseProxyClustersArgs struct {
	Type *string `pulumi:"type,optional"`
}

// Annotate provides field descriptions for GetReverseProxyClustersArgs.
func (a *GetReverseProxyClustersArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Type, "Optional cluster type filter: 'account' (owned/operated by the account, BYOP) "+
		"or 'shared' (operated by NetBird and shared across accounts). When set, only clusters of this type are returned.")
}

// ProxyClusterSummary is a brief summary of a NetBird reverse proxy cluster.
type ProxyClusterSummary struct {
	ID                  string `pulumi:"id"`
	Address             string `pulumi:"address"`
	Type                string `pulumi:"type"`
	Online              bool   `pulumi:"online"`
	ConnectedProxies    int    `pulumi:"connectedProxies"`
	Private             bool   `pulumi:"private"`
	RequireSubdomain    bool   `pulumi:"requireSubdomain"`
	SupportsCrowdsec    bool   `pulumi:"supportsCrowdsec"`
	SupportsCustomPorts bool   `pulumi:"supportsCustomPorts"`
}

// Annotate provides field descriptions for ProxyClusterSummary.
func (p *ProxyClusterSummary) Annotate(ann infer.Annotator) {
	ann.Describe(&p.ID, "Unique identifier of the proxy cluster.")
	ann.Describe(&p.Address, "Cluster address used for CNAME targets; the value to pass as a ReverseProxyDomain targetCluster.")
	ann.Describe(&p.Type, "Source of the cluster: 'account' (BYOP) or 'shared' (operated by NetBird).")
	ann.Describe(&p.Online, "Whether at least one proxy in the cluster has heartbeated within the active window.")
	ann.Describe(&p.ConnectedProxies, "Number of proxy nodes currently connected.")
	ann.Describe(&p.Private, "True when at least one connected proxy is embedded in a netbird client and serving over a WireGuard tunnel.")
	ann.Describe(&p.RequireSubdomain, "Whether services on this cluster must include a subdomain label.")
	ann.Describe(&p.SupportsCrowdsec, "Whether all active proxies in the cluster have CrowdSec configured.")
	ann.Describe(&p.SupportsCustomPorts, "Whether the cluster supports binding arbitrary TCP/UDP ports.")
}

// GetReverseProxyClustersResult is the output of GetReverseProxyClusters.
type GetReverseProxyClustersResult struct {
	Clusters []ProxyClusterSummary `pulumi:"clusters"`
}

// Annotate provides field descriptions for GetReverseProxyClustersResult.
func (r *GetReverseProxyClustersResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Clusters, "The list of reverse proxy clusters matching the filter criteria.")
}

// Invoke lists reverse proxy clusters, applying an optional type filter.
func (f *GetReverseProxyClusters) Invoke(
	ctx context.Context,
	req infer.FunctionRequest[GetReverseProxyClustersArgs],
) (infer.FunctionResponse[GetReverseProxyClustersResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[GetReverseProxyClustersResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiClusters, err := client.ReverseProxyClusters.List(ctx)
	if err != nil {
		return infer.FunctionResponse[GetReverseProxyClustersResult]{}, fmt.Errorf("listing reverse proxy clusters failed: %w", err)
	}

	clusters := make([]ProxyClusterSummary, 0, len(apiClusters))

	for _, cluster := range apiClusters {
		if req.Input.Type != nil && string(cluster.Type) != *req.Input.Type {
			continue
		}

		clusters = append(clusters, ProxyClusterSummary{
			ID:                  cluster.Id,
			Address:             cluster.Address,
			Type:                string(cluster.Type),
			Online:              cluster.Online,
			ConnectedProxies:    cluster.ConnectedProxies,
			Private:             derefBool(cluster.Private),
			RequireSubdomain:    derefBool(cluster.RequireSubdomain),
			SupportsCrowdsec:    derefBool(cluster.SupportsCrowdsec),
			SupportsCustomPorts: derefBool(cluster.SupportsCustomPorts),
		})
	}

	return infer.FunctionResponse[GetReverseProxyClustersResult]{
		Output: GetReverseProxyClustersResult{
			Clusters: clusters,
		},
	}, nil
}

// derefBool returns the value of an optional bool pointer, defaulting to false when nil.
func derefBool(b *bool) bool {
	if b == nil {
		return false
	}

	return *b
}
