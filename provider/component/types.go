package component

import (
	"fmt"
	"strings"
)

type ACLFile struct {
	Groups GroupMap `json:"groups"`
	ACLs   []ACL    `json:"acls"`
}

type GroupMap map[string][]string

type ACL struct {
	Action string   `json:"action"`
	Src    []string `json:"src"`
	Dst    []string `json:"dst"`
	Proto  string   `json:"proto,omitempty"`
}

type Network struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NetworkResource struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	NetworkID   string   `json:"network_id"`
	Address     string   `json:"address"`
	Enabled     bool     `json:"enabled"`
	GroupIDs    []string `json:"group_ids"`
}

type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Policy struct {
	Type       string       `yaml:"type"`
	Properties PolicyFields `yaml:"properties"`
}

type PolicyFields struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Enabled     bool         `yaml:"enabled"`
	Rules       []PolicyRule `yaml:"rules"`
}

type PolicyRule struct {
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description,omitempty"`
	Action        string   `yaml:"action"`
	Enabled       bool     `yaml:"enabled"`
	Bidirectional bool     `yaml:"bidirectional"`
	Protocol      string   `yaml:"protocol"`
	Ports         []string `yaml:"ports,omitempty"`
	Sources       []string `yaml:"sources,omitempty"`
	Destinations  []string `yaml:"destinations,omitempty"`
}

// ACLRule represents a decomposed ACL destination
type ACLRule struct {
	Name    string              // Used as resource name
	Address *string             // IP or CIDR
	Ports   *[]string           // If present
	Group   *string             // If present
	Type    string              // Type of resource
	Dest    *map[string]ACLRule // Only Src Rules must have Dest not nil
}

func (v ACLRule) String() string {
	var addr, group, ports string

	if v.Address != nil {
		addr = *v.Address
	} else {
		addr = "<nil>"
	}

	if v.Group != nil {
		group = *v.Group
	} else {
		group = "<nil>"
	}

	if v.Ports != nil {
		ports = strings.Join(*v.Ports, ",")
	} else {
		ports = "<nil>"
	}

	return fmt.Sprintf("Type=%s Name=%s Address=%s, Group=%s, Ports=%s",
		v.Type,
		v.Name,
		addr,
		group,
		ports,
	)
}
