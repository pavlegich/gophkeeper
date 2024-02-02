package utils

import "strconv"

// IsValidDataType checks whether the data type is correct.
func IsValidDataType(t string) bool {
	if t == "credentials" || t == "text" || t == "binary" || t == "card" {
		return true
	}
	return false
}

// IsValidCardNumber checks whether the bank card number
// is valid using Luhn algorithm.
func IsValidCardNumber(number int) bool {
	if len(strconv.Itoa(number)) != 16 {
		return false
	}
	return (number%10+checksum(number/10))%10 == 0
}

// checkSum checks one part of bank card number
// for validity using Luhn algorithm.
func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
