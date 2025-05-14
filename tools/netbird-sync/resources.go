package main

import (
	"fmt"
	"strings"
)

func generateGroupResources(groups GroupMap) map[string]any {
	groupResources := make(map[string]any)

	for key, peers := range groups {
		name := strings.TrimPrefix(key, "group:")
		groupResources[name] = map[string]any{
			"type": "netbird:resource:Group",
			"properties": map[string]any{
				"name":  name,
				"peers": peers,
			},
		}
	}

	return groupResources
}

func generatePolicyResources(acls []ACL) map[string]Policy {
	resources := make(map[string]Policy)

	for i, acl := range acls {
		name := fmt.Sprintf("policy-%03d", i+1)

		rule := PolicyRule{
			Name:          fmt.Sprintf("ACL Rule %03d", i+1),
			Description:   fmt.Sprintf("From %s to %s", strings.Join(acl.Src, ","), strings.Join(acl.Dst, ",")),
			Action:        acl.Action,
			Enabled:       true,
			Bidirectional: false,
			Protocol:      defaultProto(acl.Proto),
			Sources:       acl.Src,
			Destinations:  acl.Dst,
		}

		res := Policy{
			Type: "netbird:resource:Policy",
			Properties: PolicyFields{
				Name:        fmt.Sprintf("ACL %03d", i+1),
				Description: "Auto-generated from ACL JSON",
				Enabled:     true,
				Rules:       []PolicyRule{rule},
			},
		}

		resources[name] = res
	}

	return resources
}

func defaultProto(proto string) string {
	if proto == "" {
		return "tcp"
	}
	return proto
}
