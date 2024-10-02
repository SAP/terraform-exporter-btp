package tfimportprovider

import "slices"

func isSubset(superSet []string, subset []string) (string, bool) {
	for _, value := range subset {
		if !slices.Contains(superSet, value) {
			return value, false
		}
	}
	return "", true
}
