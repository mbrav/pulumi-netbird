package resource

import (
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
func equalSlicePtr(sliceA, sliceB *[]string) bool {
	if sliceA == nil && sliceB == nil {
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

func equalResourcesPtr(resourcesA, resourcesB *[]Resource) bool {
	if resourcesA == nil && resourcesB == nil {
		return true
	}

	if resourcesA == nil || resourcesB == nil {
		return false
	}

	if len(*resourcesA) != len(*resourcesB) {
		return false
	}

	aSorted := slices.Clone(*resourcesA)
	bSorted := slices.Clone(*resourcesB)

	slices.SortFunc(aSorted, func(resA, resB Resource) int {
		if resA.Type != resB.Type {
			if resA.Type < resB.Type {
				return -1
			}

			return 1
		}

		return strings.Compare(resA.ID, resB.ID)
	})

	slices.SortFunc(bSorted, func(resA, resB Resource) int {
		if resA.Type != resB.Type {
			if resA.Type < resB.Type {
				return -1
			}

			return 1
		}

		return strings.Compare(resA.ID, resB.ID)
	})

	for i := range aSorted {
		if !equalResourcePtr(&aSorted[i], &bSorted[i]) {
			return false
		}
	}

	return true
}

// equalReverseProxyTargets compares two slices of ReverseProxyTarget by their key fields.
func equalReverseProxyTargets(targetsA, targetsB []ReverseProxyTarget) bool {
	if len(targetsA) != len(targetsB) {
		return false
	}

	for idx := range targetsA {
		if targetsA[idx].Enabled != targetsB[idx].Enabled ||
			targetsA[idx].Port != targetsB[idx].Port ||
			targetsA[idx].Protocol != targetsB[idx].Protocol ||
			targetsA[idx].TargetType != targetsB[idx].TargetType ||
			!equalPtr(targetsA[idx].Host, targetsB[idx].Host) ||
			!equalPtr(targetsA[idx].Path, targetsB[idx].Path) {
			return false
		}
	}

	return true
}
