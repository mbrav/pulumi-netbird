package resource

import (
	"slices"
)

// Helper to stringify a pointer safely.
func strPtr(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}

// Helper to compare string pointers safely.
func equalPtr(atrA, strB *string) bool {
	if atrA == nil && strB == nil {
		return true
	}

	if atrA == nil || strB == nil {
		return false
	}

	return *atrA == *strB
}

// helper function to compare optional []string values.
func equalSlicePtr(sliceA, sliceB *[]string) bool {
	if sliceA == nil && sliceB == nil {
		return true
	}

	if sliceA == nil || sliceB == nil {
		return false
	}

	// Copy and sort both before comparing
	aSorted := slices.Clone(*sliceA)
	bSorted := slices.Clone(*sliceB)
	slices.Sort(aSorted)
	slices.Sort(bSorted)

	if len(aSorted) != len(bSorted) {
		return false
	}

	for i := range aSorted {
		if aSorted[i] != bSorted[i] {
			return false
		}
	}

	return true
}
