package bot

// RemoveIndex is just a general function which removes
// an item from a slice given an index.
func RemoveIndex(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
