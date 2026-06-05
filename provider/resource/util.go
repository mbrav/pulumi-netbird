package resource

import (
	"fmt"
	"slices"
	"strings"
)

// strPtr helper function to stringify a pointer safely.
func strPtr(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}

// equalPtr compares two pointers of any comparable type safely.
func equalPtr[T comparable](ptrA, ptrB *T) bool {
	if ptrA == nil && ptrB == nil {
		return true
	}

	if ptrA == nil || ptrB == nil {
		return false
	}

	return *ptrA == *ptrB
}

// equalSlice compares two []string slices, ignoring order.
func equalSlice(sliceA, sliceB []string) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}

	aSorted := slices.Clone(sliceA)
	bSorted := slices.Clone(sliceB)

	slices.Sort(aSorted)
	slices.Sort(bSorted)

	for i := range aSorted {
		if aSorted[i] != bSorted[i] {
			return false
		}
	}

	return true
}

// equalSlicePtr compares two *[]string values by delegating to equalSlice.
// Treats nil and empty slice as equal.
func equalSlicePtr(sliceA, sliceB *[]string) bool {
	aLen := 0
	if sliceA != nil {
		aLen = len(*sliceA)
	}

	bLen := 0
	if sliceB != nil {
		bLen = len(*sliceB)
	}

	if aLen == 0 && bLen == 0 {
		return true
	}

	if sliceA == nil || sliceB == nil {
		return false
	}

	return equalSlice(*sliceA, *sliceB)
}

// boolVal safely converts a pointer to bool to a bool value.
func boolVal(p *bool) bool {
	if p == nil {
		return false
	}

	return *p
}

func isBlank(v string) bool {
	return strings.TrimSpace(v) == ""
}

// isNotFoundErr returns true when err represents a 404 / "not found" response from the NetBird API.
func isNotFoundErr(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "not found")
}

// parseNestedID splits a compound "<parentID>/<childID>" import ID.
// Both parts must be non-empty; otherwise an error is returned to the caller.
func parseNestedID(kind, id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("%s import ID must be in the format <parentID>/<childID>, got %q", kind, id)
	}

	return parts[0], parts[1], nil
}

// sortedStrings returns a sorted clone of s, leaving the original unmodified.
func sortedStrings(s []string) []string {
	c := slices.Clone(s)
	slices.Sort(c)

	return c
}
