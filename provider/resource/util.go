package resource

import (
	"slices"
)

// Helper to stringify a pointer safely.
func strPtr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// Helper to compare string pointers safely.
func equalPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

// helper function to compare optional []string values.
func equalSlicePtr(a, b *[]string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	// Copy and sort both before comparing
	aSorted := slices.Clone(*a)
	bSorted := slices.Clone(*b)
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
