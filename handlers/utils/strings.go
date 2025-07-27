package utils

// Pluralise returns the plural form of a word based on the count.
func Pluralise(count int64, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
