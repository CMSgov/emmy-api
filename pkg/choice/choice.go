package choice

// Ternary operator
//
//nolint:revive // Condition flag is intentional for a ternary-style helper.
func Ternary[T any](condition bool, isTrue, isFalse T) T {
	if condition {
		return isTrue
	}

	return isFalse
}

// Function ternary operator
//
//nolint:revive // Condition flag is intentional for a ternary-style helper.
func FuncTernary[T any](condition bool, isTrue, isFalse func() T) T {
	if condition {
		return isTrue()
	}

	return isFalse()
}
