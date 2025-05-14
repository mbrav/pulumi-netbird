package main

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

// Output-friendly struct (matches Pulumi YAML format)
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
