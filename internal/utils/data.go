package utils

// IsCorrectDataType checks whether the data type is correct.
func IsCorrectDataType(t string) bool {
	if t == "CREDENTIALS" || t == "TEXT" || t == "BINARY" || t == "CARD" {
		return true
	}
	return false
}
