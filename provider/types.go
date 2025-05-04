package provider

import nbapi "github.com/netbirdio/netbird/management/server/http/api"

// Group represents a resource for managing NetBird groups.
type Group struct{}

// GroupArgs represents the input arguments for creating or updating a group.
type GroupArgs struct {
	Name  string    `pulumi:"name"`
	Peers *[]string `pulumi:"peers,optional"`
}

// GroupState represents the state of the group resource.
type GroupState struct {
	Name  string    `pulumi:"name"`
	Peers *[]string `pulumi:"peers,optional"`
	NbID  string    `pulumi:"nbId"`
}

// Network represents a resource for managing NetBird networks.
type Network struct{}

// NetworkArgs represents the input arguments for creating or updating a network.
type NetworkArgs struct {
	Name        string  `pulumi:"name"`
	Description *string `pulumi:"description,optional"`
}

// NetworkState represents the state of the network resource.
type NetworkState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	NetworkArgs
	Name        string  `pulumi:"name"`
	Description *string `pulumi:"description,optional"`
	NbID        string  `pulumi:"nbId"`
}

// NetworkResource represents a Pulumi resource for NetBird network resources.
type NetworkResource struct{}

// NetworkResourceArgs represents the input arguments for creating or updating a network resource.
type NetworkResourceArgs struct {
	Name        string    `pulumi:"name"`
	Description *string   `pulumi:"description,optional"`
	NetworkID   string    `pulumi:"network_id"`
	Address     string    `pulumi:"address"`
	Enabled     bool      `pulumi:"enabled"`
	GroupIDs    *[]string `pulumi:"group_ids,optional"`
}

// NetworkResourceState represents the state of a network resource.
type NetworkResourceState struct {
	Name        string    `pulumi:"name"`
	Description *string   `pulumi:"description"`
	NbID        string    `pulumi:"nbId"`
	NetworkID   string    `pulumi:"network_id"`
	Address     string    `pulumi:"address"`
	Enabled     bool      `pulumi:"enabled"`
	GroupIDs    *[]string `pulumi:"group_ids,optional"`
}

// NetworkRouter represents a Pulumi resource for NetBird network resources.
type NetworkRouter struct{}

// NetworkRouterArgs represents the input arguments for creating or updating a network router.
type NetworkRouterArgs struct {
	NetworkID  string    `pulumi:"network_id"`
	Enabled    bool      `pulumi:"enabled"`
	Masquerade bool      `pulumi:"masquerade"`
	Metric     int       `pulumi:"metric"`
	Peer       *string   `pulumi:"peer,optional"`
	PeerGroups *[]string `pulumi:"peer_groups,optional"`
}

// NetworkResourceArgs represents the state of a network router.
type NetworkRouterState struct {
	NbID       string    `pulumi:"nbId"`
	NetworkID  string    `pulumi:"network_id"`
	Enabled    bool      `pulumi:"enabled"`
	Masquerade bool      `pulumi:"masquerade"`
	Metric     int       `pulumi:"metric"`
	Peer       *string   `pulumi:"peer"`
	PeerGroups *[]string `pulumi:"peer_groups"`
}

// Peer represents a resource for managing NetBird peers.
type Peer struct{}

// PeerArgs represents the input arguments for a peer resource.
type PeerArgs struct {
	PeerID string `pulumi:"peerId"`
}

// PeerState represents the state of the peer resource.
type PeerState struct {
	PeerID     string `pulumi:"peerId"`
	Name       string `pulumi:"name"`
	SshEnabled bool   `pulumi:"sshEnabled"`
}

// Policy represents a resource for managing NetBird policies.
type Policy struct{}

// PolicyArgs are the input arguments for a policy resource.
type PolicyArgs struct {
	Name                string                   `pulumi:"name"`
	Description         *string                  `pulumi:"description"`
	Enabled             bool                     `pulumi:"enabled"`
	Rules               []nbapi.PolicyRuleUpdate `pulumi:"rules"`
	SourcePostureChecks *[]string                `pulumi:"sourcePostureChecks"`
}

// PolicyState is the persisted state of the resource.
type PolicyState struct {
	NbID                string                   `pulumi:"nbId"`
	Name                string                   `pulumi:"name"`
	Description         *string                  `pulumi:"description"`
	Enabled             bool                     `pulumi:"enabled"`
	Rules               []nbapi.PolicyRuleUpdate `pulumi:"rules"`
	SourcePostureChecks *[]string                `pulumi:"sourcePostureChecks"`
}
