package component

import (
	"fmt"
	"strings"
)

// parseACLRules builds a nested structure where each source ACLRule contains its destination rules.
func parseACLRules(acls []ACL, srcRules map[string]*ACLRule) {
	for _, acl := range acls {
		for _, src := range acl.Src {
			// Parse the source ACL entry to create a source rule
			srcRule := parseACLEntry(src)
			srcKey := fmt.Sprintf("%s-%s", srcRule.Type, srcRule.Name)

			// Use existing srcRule entry if it exists, otherwise initialize a new one
			if existing, exists := srcRules[srcKey]; exists {
				srcRule = existing
			} else {
				// Initialize the Dest map for the new srcRule
				srcRule.Dest = &map[string]ACLRule{}
				srcRules[srcKey] = srcRule
			}

			// Add destination rules to the source rule
			for _, dst := range acl.Dst {
				dstRule := parseACLEntry(dst)
				dstKey := fmt.Sprintf("%s-%s", dstRule.Type, dstRule.Name)

				// Initialize Dest map if it is nil (should not happen but for safety)
				if srcRule.Dest == nil {
					srcRule.Dest = &map[string]ACLRule{}
				}

				// Add the destination rule if it does not already exist
				if _, exists := (*srcRule.Dest)[dstKey]; !exists {
					(*srcRule.Dest)[dstKey] = *dstRule
				}
			}
		}
	}
}

// parseGroupRules adds destination group entries to each user listed in the GroupMap.
// Each key in the map is a group name (used as a destination), and the value is a list of users (sources).
func parseGroupRules(groups GroupMap, srcRules map[string]*ACLRule) {
	for dstEntry, users := range groups {
		dstRule := parseACLEntry(dstEntry)
		if dstRule == nil {
			continue
		}
		dstKey := fmt.Sprintf("%s-%s", dstRule.Type, dstRule.Name)

		for _, user := range users {
			srcEntry := fmt.Sprintf("group:%s", user)
			srcRule := parseACLEntry(srcEntry)
			if srcRule == nil {
				continue
			}
			srcKey := fmt.Sprintf("%s-%s", srcRule.Type, srcRule.Name)

			// Get or create the source rule
			if existing, exists := srcRules[srcKey]; exists {
				srcRule = existing
			} else {
				srcRule.Dest = &map[string]ACLRule{}
				srcRules[srcKey] = srcRule
			}

			// Initialize destination map if needed
			if srcRule.Dest == nil {
				srcRule.Dest = &map[string]ACLRule{}
			}

			// Add the destination rule
			if _, exists := (*srcRule.Dest)[dstKey]; !exists {
				(*srcRule.Dest)[dstKey] = *dstRule
			}
		}
	}
}

// parseACLEntry parses a single ACL entry and returns an ACLRule.
func parseACLEntry(entry string) *ACLRule {
	var addr *string
	var ports *[]string
	var name string

	// Check if the entry is a tag or group
	if strings.HasPrefix(entry, "tag:") || strings.HasPrefix(entry, "group:") {
		tagParts := strings.SplitN(entry, ":", 3)
		ports = parsePorts(tagParts)

		// Replace '*' with 'star' in the name
		name = strings.ReplaceAll(tagParts[1], "*", "star")

		return &ACLRule{
			Name:    name,
			Group:   &name,
			Address: nil,
			Ports:   ports,
			Type:    "group",
		}
	}

	// Parse address and ports from the entry
	addrStr, portsSlice := parseAddressAndPorts(entry)
	addr = stringPtrIfNotEmpty(addrStr)
	if len(portsSlice) > 0 {
		ports = &portsSlice
	}

	// Generate name and type based on the address
	name, resType := generateNameFromAddress(addrStr)

	return &ACLRule{
		Name:    name,
		Group:   nil,
		Address: addr,
		Ports:   ports,
		Type:    resType,
	}
}

// parsePorts extracts ports from the given parts of an ACL entry.
func parsePorts(parts []string) *[]string {
	if len(parts) == 3 && parts[2] != "*" {
		split := strings.Split(parts[2], ",")
		return &split
	}
	return nil
}

// parseAddressAndPorts splits the entry into address and ports.
func parseAddressAndPorts(entry string) (string, []string) {
	parts := strings.SplitN(entry, ":", 2)
	addr := parts[0]
	if len(parts) == 2 && parts[1] != "*" {
		return addr, strings.Split(parts[1], ",")
	}
	return addr, nil
}

// generateNameFromAddress generates a name and type based on the address format.
func generateNameFromAddress(addr string) (string, string) {
	base := strings.ReplaceAll(addr, ".", "-")
	base = strings.ReplaceAll(base, "/", "-")
	if strings.HasSuffix(addr, "/32") {
		return base[:len(base)-3], "host" // Remove the /32 suffix for host type
	}
	return base, "cidr" // Default to CIDR type
}

// stringPtrIfNotEmpty returns a pointer to the string if it is not empty, otherwise returns nil.
func stringPtrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
