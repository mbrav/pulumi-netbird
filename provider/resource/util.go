package resource

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/netbirdio/netbird/shared/management/client/rest"
)

// strPtr helper function to stringify a pointer safely.
func strPtr(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}

// equalOptionalStr compares an optional string input against the value the API
// reports, treating nil and "" as the same thing.
//
// Use it for optional string fields the API always returns on the wire (a
// non-pointer `string` in the response type) and normalises so that empty
// means unset. Comparing those with equalPtr reports a diff on every preview
// once the field is omitted from the configuration — nil input never equals a
// pointer to "" — and the update never converges. Resources using this must
// also send an explicit "" for a nil input (see strPtr) so that removing the
// field from the configuration actually clears it server-side.
func equalOptionalStr(input, state *string) bool {
	return strPtr(input) == strPtr(state)
}

// equalServerAssignedPtr compares an optional input against server state,
// accepting whatever the server assigned when the input is nil.
//
// Use it for fields the server fills in or defaults on its own and that a
// request cannot clear (an IdP sync interval, a user's name and email). For
// those, a nil input can never equal the returned value, so an unconditional
// comparison proposes the same update — or, worse, the same replacement — on
// every preview forever. Once the field is set in the configuration it is
// compared normally, so real changes are still detected.
func equalServerAssignedPtr[T comparable](input, state *T) bool {
	if input == nil {
		return true
	}

	return equalPtr(input, state)
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

// isConflictErr returns true when err represents a 409 Conflict response from the NetBird API.
func isConflictErr(err error) bool {
	apiErr, ok := errors.AsType[*rest.APIError](err)
	if !ok {
		return false
	}

	return apiErr.StatusCode == http.StatusConflict
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
