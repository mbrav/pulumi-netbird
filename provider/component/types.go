package component

//
// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"sort"
// 	"strings"
//
// 	"github.com/mbrav/pulumi-netbird/sdk/go/netbird/resource"
// )
//
// type ACLJSON struct {
// 	Groups    GroupMap            `json:"groups"`
// 	TagOwners map[string][]string `json:"tagOwners"`
// 	ACLs      []ACL               `json:"acls"`
// }
//
// type GroupMap map[string][]string
//
// type ACL struct {
// 	Action string   `json:"action"`
// 	Src    []string `json:"src"`
// 	Dst    []string `json:"dst"`
// 	Proto  *string  `json:"proto,omitempty"`
// }
//
// // ACLEntry represents a decomposed ACL destination
// type ACLEntry struct {
// 	Type ACLType `json:"type"` // Type of ACL
//
// 	// For Groups and Network resources
// 	Name     string             `json:"name"`               // Used as resource name
// 	Address  *string            `json:"address,omitempty"`  // IP or CIDR
// 	Ports    *[]string          `json:"ports,omitempty"`    // If present
// 	Protocol *resource.Protocol `json:"protocol,omitempty"` // Optional protocol Type
//
// 	// For Policy Rules
// 	Sources      *map[string]ACLEntry `json:"sources,omitempty"`      // Source Resources
// 	Destinations *map[string]ACLEntry `json:"destinations,omitempty"` // Destination Resources
// 	Destination  *ACLEntry            `json:"destination,omitempty"`  // Single destination Resource
//
// 	// For Pulumi resources
// 	ResourceGroup          *resource.Group           `json:"resourceGroup,omitempty"`   // Optional resource Group
// 	ResourceNeworkResource *resource.NetworkResource `json:"networkResource,omitempty"` // Optional resource NetworkResource
// 	ResourcePolicy         *resource.Policy          `json:"policy,omitempty"`          // Optional resource Policy
// }
//
// type ACLType string
//
// const (
// 	ACLTypeHost                = ACLType("host")
// 	ACLTypeNetwork             = ACLType("network")
// 	ACLTypeGroup               = ACLType("group")
// 	ACLTypeNetResToGroup       = ACLType("net-res-to-group")
// 	ACLTypePolicyGroupToGroup  = ACLType("policy-gtg")
// 	ACLTypePolicyGroupToNetRes = ACLType("policy-gtn")
// )
//
// // String method for ACLEntry
// func (v ACLEntry) String() string {
// 	var addr, groups, ports string
//
// 	if v.Address != nil {
// 		addr = *v.Address
// 	} else {
// 		addr = "<nil>"
// 	}
//
// 	if v.Sources != nil {
// 		groups = fmt.Sprintf("%d", len(*v.Sources))
// 	} else {
// 		groups = "<nil>"
// 	}
//
// 	if v.Ports != nil {
// 		ports = strings.Join(*v.Ports, ",")
// 	} else {
// 		ports = "<nil>"
// 	}
//
// 	return fmt.Sprintf("Type=%s Name=%s Address=%s, Group=%s, Ports=%s",
// 		v.Type,
// 		v.Name,
// 		addr,
// 		groups,
// 		ports,
// 	)
// }
//
// // Key returns or generates a hash for the ACL rule.
// func (v ACLEntry) Key() string {
// 	switch v.Type {
// 	case ACLTypePolicyGroupToNetRes:
// 		var srcString, dstString string
//
// 		// Sort map keys for consistent source order
// 		srcKeys := make([]string, 0, len(*v.Sources))
// 		for key := range *v.Sources {
// 			srcKeys = append(srcKeys, key)
// 		}
// 		sort.Strings(srcKeys)
// 		for _, key := range srcKeys {
// 			srcString += (*v.Sources)[key].Key()
// 		}
//
// 		// Hash sources but keep destination key as-is
// 		srcHash := hash(srcString)[:10]
// 		dstString = v.Destination.Key()
//
// 		return fmt.Sprintf("%s-%s-%s", v.Type, srcHash, dstString)
//
// 	case ACLTypePolicyGroupToGroup:
// 		var srcString, dstString string
//
// 		// Sort and process sources
// 		srcKeys := make([]string, 0, len(*v.Sources))
// 		for key := range *v.Sources {
// 			srcKeys = append(srcKeys, key)
// 		}
// 		sort.Strings(srcKeys)
// 		for _, key := range srcKeys {
// 			srcString += (*v.Sources)[key].Key()
// 		}
//
// 		// Sort and process destinations
// 		dstKeys := make([]string, 0, len(*v.Destinations))
// 		for key := range *v.Destinations {
// 			dstKeys = append(dstKeys, key)
// 		}
// 		sort.Strings(dstKeys)
// 		for _, key := range dstKeys {
// 			dstString += (*v.Destinations)[key].Key()
// 		}
//
// 		srcHash := hash(srcString)[:10]
// 		dstHash := hash(dstString)[:10]
//
// 		return fmt.Sprintf("%s-%s-%s", v.Type, srcHash, dstHash)
//
// 	default:
// 		return fmt.Sprintf("%s-%s", v.Type, v.Name)
// 	}
// }
//
// // Hash function to generate a hash from a string.
// func hash(s string) string {
// 	h := sha256.New()
// 	h.Write([]byte(s))
// 	return hex.EncodeToString(h.Sum(nil))
// }
