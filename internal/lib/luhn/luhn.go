package luhn

import (
	"strconv"
	"strings"
)

// Valid check number is valid or not based on Luhn algorithm
func Valid(order string) bool {
	orderID := strings.ReplaceAll(order, " ", "")
	sum := 0
	shouldDouble := len(orderID)%2 == 0

	for i := 0; i <= len(orderID)-1; i++ {
		digit, err := strconv.Atoi(string(orderID[i]))
		if err != nil {
			return false
		}

		if shouldDouble {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}
