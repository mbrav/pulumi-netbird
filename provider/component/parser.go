package component

//
// import (
// 	"fmt"
// 	"strings"
//
// 	"github.com/mbrav/pulumi-netbird/sdk/go/netbird/resource"
// )
//
// // parseACLs processes both ACL and group rules efficiently.
// // It avoids creating a Cartesian product for resources and groups group rules together.
// // The function takes an ACLJSON pointer as input and returns a map of ACLEntry keyed by their unique string keys.
// // It initializes group and subgroup mappings, processes ACL entries to populate group and network resource entries,
// // and generates policy rules for group-to-group and group-to-host relationships.
// func parseACLs(aclJson *ACLJSON) map[string]ACLEntry {
// 	// Initialize the ruleMap
// 	aclMap := make(map[string]ACLEntry)
//
// 	// Create groups
// 	// And generate a map for storing group to group maps
// 	subGroupMap := make(map[string][]ACLEntry)
// 	for group, sG := range aclJson.Groups {
// 		cleanGroupName := strings.ReplaceAll(group, "group:", "")
// 		parentGroup := ACLEntry{
// 			Name: cleanGroupName,
// 			Type: ACLTypeGroup,
// 		}
// 		// Check if group exists, otherwise create
// 		if _, exists := aclMap[parentGroup.Key()]; !exists {
// 			// Save the unique key to map
// 			aclMap[parentGroup.Key()] = parentGroup
// 		}
//
// 		// Generate child groups
// 		var childGroups []ACLEntry
// 		for _, cG := range sG {
// 			cleanCildGroupName := strings.ReplaceAll(cG, "group:", "")
// 			childGroup := ACLEntry{
// 				Name: cleanCildGroupName,
// 				Type: ACLTypeGroup,
// 			}
// 			// Check if group exists
// 			if _, exists := aclMap[childGroup.Key()]; !exists {
// 				// Save the unique key to map
// 				aclMap[childGroup.Key()] = childGroup
// 			}
// 			childGroups = append(childGroups, childGroup)
// 		}
//
// 		// Save the unique key to map
// 		subGroupMap[parentGroup.Key()] = childGroups
// 	}
//
// 	for index, acl := range aclJson.ACLs {
// 		// Initialize rule
//
// 		// Get rules
// 		groupSources, netSources := parseACLEntries(acl.Dst)
// 		groupDestinations, netDestinations := parseACLEntries(acl.Dst)
//
// 		// Create source Groups
// 		for g, group := range groupSources {
// 			// Check if group exists in subGroupMap
// 			if childGroups, exists := subGroupMap[g]; exists {
// 				for _, cG := range childGroups {
// 					groupSources[cG.Key()] = cG
// 				}
// 			}
// 			// Check if key exists
// 			if _, exists := aclMap[g]; exists {
// 				continue
// 			}
// 			// Save the unique key to map
// 			aclMap[g] = group
// 		}
//
// 		// Create destination Groups
// 		for g, group := range groupDestinations {
// 			// Check if group exists in subGroupMap
// 			if childGroups, exists := subGroupMap[g]; exists {
// 				for _, cG := range childGroups {
// 					groupDestinations[cG.Key()] = cG
// 				}
// 			}
// 			// Check if key exists
// 			if _, exists := aclMap[g]; exists {
// 				continue
// 			}
// 			// Save the unique key to map
// 			aclMap[g] = group
// 		}
//
// 		// Create source NetworkResources
// 		for e, entry := range netSources {
// 			// Delete entry for netsources
// 			delete(netSources, e)
// 			// Convert NetworkResource type to a NetworkResource to Group conversion type
// 			entry.Type = ACLTypeNetResToGroup
// 			// Set Protocol and Ports to nil since these irrelevant for Groups
// 			// entry.Ports = nil
// 			// entry.Protocol = nil
// 			entryKey := entry.Key()
// 			groupSources[entryKey] = entry
// 			// Check if key exists
// 			if _, exists := aclMap[entryKey]; exists {
// 				continue
// 			}
// 			// Save the unique key to map
// 			aclMap[entryKey] = entry
// 		}
//
// 		// Create destination NetworkResources
// 		for n, net := range netDestinations {
// 			// Check if key exists
// 			if _, exists := aclMap[n]; exists {
// 				continue
// 			}
// 			// Save the unique key to map
// 			aclMap[n] = net
// 		}
//
// 		// Set protocol
// 		var protocol *resource.Protocol
// 		if acl.Proto != nil {
// 			switch *acl.Proto {
// 			case "tcp":
// 				p := resource.ProtocolTcp
// 				protocol = &p
// 			case "udp":
// 				p := resource.ProtocolUdp
// 				protocol = &p
// 			case "icmp":
// 				p := resource.ProtocolIcmp
// 				protocol = &p
// 			default:
// 				protocol = nil
// 			}
// 		}
//
// 		// Generate acl with group to group (many-to-many)
// 		if len(groupDestinations) > 0 {
// 			aclRule := ACLEntry{
// 				Type:         ACLTypePolicyGroupToGroup,
// 				Sources:      &groupSources,
// 				Protocol:     protocol,
// 				Destinations: &groupDestinations,
// 			}
// 			keyName := aclRule.Key()
// 			aclRule.Name = keyName
// 			// Check if key exists
// 			if _, exists := aclMap[keyName]; exists {
// 				keyName = fmt.Sprintf("%s-%d", keyName, index)
// 			}
//
// 			// Save the unique key to map
// 			aclMap[keyName] = aclRule
// 		}
//
// 		// Generate acl with group to host (many-to-one)
// 		for _, net := range netDestinations {
// 			aclRule := ACLEntry{
// 				Type:        ACLTypePolicyGroupToNetRes,
// 				Sources:     &groupSources,
// 				Protocol:    protocol,
// 				Destination: &net,
// 				Ports:       net.Ports,
// 			}
// 			keyName := aclRule.Key()
// 			aclRule.Name = keyName
// 			// Check if key exists
// 			if _, exists := aclMap[keyName]; exists {
// 				keyName = fmt.Sprintf("%s-%d", keyName, index)
// 			}
//
// 			// Save the unique key to map
// 			aclMap[keyName] = aclRule
// 		}
// 	}
// 	return aclMap
// }
//
// // parseACLEntries parses a list of ACL entry strings and returns two maps of ACLEntry objects.
// // The first map contains group-based ACL entries, and the second contains network-based ACL entries.
// // Each entry string is analyzed to determine if it represents a group/tag or a network address with optional ports and protocol.
// // The function normalizes group names, infers protocols based on port numbers, and determines the type of each ACL entry.
// func parseACLEntries(entries []string) (map[string]ACLEntry, map[string]ACLEntry) {
// 	groupACLs := make(map[string]ACLEntry)
// 	netACLs := make(map[string]ACLEntry)
// 	for _, entry := range entries {
//
// 		var ports *[]string
// 		var name string
// 		var protocol *resource.Protocol
//
// 		// Check if the entry is a tag or group
// 		if strings.HasPrefix(entry, "tag:") || strings.HasPrefix(entry, "group:") {
// 			tagParts := strings.SplitN(entry, ":", 3)
//
// 			// Replace '*' with 'star' in the name
// 			nam = strings.ReplaceAll(tagParts[1], ":", "")
// 			groupACL := ACLEntry{
// 				Name: name,
// 				Type: ACLTypeGroup,
// 			}
//
// 			// Save the unique key to map
// 			groupACLs[groupACL.Key()] = groupACL
// 			continue
// 		}
//
// 		// Parse address and ports
// 		addrStr, portsSlice := parseAddressAndPorts(entry)
//
// 		// Determine ports
// 		if len(portsSlice) > 0 {
// 			ports = &portsSlice
// 		}
//
// 		// Generate name and type based on the address
// 		name, resType := getAddressType(addrStr)
//
// 		netACL := ACLEntry{
// 			Name:     name,
// 			Address:  &addrStr,
// 			Ports:    ports,
// 			Protocol: protocol,
// 			Type:     resType,
// 		}
// 		// Save the unique key to map
// 		netACLs[netACL.Key()] = netACL
// 	}
// 	return groupACLs, netACLs
// }
//
// // parseAddressAndPorts splits an entry string into an address and a slice of ports.
// // The entry should be in the format "address:port1,port2" or just "address".
// // If no CIDR is present in the address, "/32" is appended to treat it as a host address.
// // If ports are specified (and not "*"), they are split into a slice; otherwise, nil is returned for ports.
// func parseAddressAndPorts(entry string) (string, []string) {
// 	parts := strings.SplitN(entry, ":", 2)
// 	addr := parts[0]
//
// 	// Append /32 only if no CIDR is present
// 	if !strings.Contains(addr, "/") {
// 		addr += "/32"
// 	}
//
// 	if len(parts) == 2 && parts[1] != "*" {
// 		return addr, strings.Split(parts[1], ",")
// 	}
// 	return addr, nil
// }
//
// // getAddressType generates a normalized name and determines the ACLType based on the address format.
// // It replaces dots and slashes in the address with dashes to create a base name.
// // If the address does not contain a CIDR notation, or if it is explicitly a /32, it is treated as a host.
// // Otherwise, the address is treated as a network.
// func getAddressType(addr string) (string, ACLType) {
// 	base := strings.ReplaceAll(addr, ".", "-")
// 	base = strings.ReplaceAll(base, "/", "-")
//
// 	// If there's no CIDR notation at all, treat as host
// 	if !strings.Contains(addr, "/") {
// 		return base, ACLTypeHost
// 	}
//
// 	// If CIDR is explicitly /32, treat as host
// 	if strings.Contains(addr, "/32") {
// 		return base, ACLTypeHost
// 	}
//
// 	// Default: treat as network
// 	return base, ACLTypeNetwork
// }
