package collate

var (
	// binCollatorInstance is a singleton used for all collations when newCollationEnabled is false.
	binCollatorInstance = &binCollator{}
)

// Collator provides functionality for comparing strings for a given
// collation order.
type Collator interface {
	// Compare returns an integer comparing the two strings. The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
	Compare(a, b string) int
	// Key returns the collate key for str. If the collation is padding, make sure the PadLen >= len(rune[]str) in opt.
	Key(str string) []byte
	// Pattern get a collation-aware WildcardPattern.
	Pattern() WildcardPattern
}

// WildcardPattern is the interface used for wildcard pattern match.
type WildcardPattern interface {
	// Compile compiles the patternStr with specified escape character.
	Compile(patternStr string, escape byte)
	// DoMatch tries to match the str with compiled pattern, `Compile()` must be called before calling it.
	DoMatch(str string) bool
}

// GetCollator get the collator according to collate, it will return the binary collator if the corresponding collator doesn't exist.
func GetCollator(collate string) Collator {
	//if atomic.LoadInt32(&newCollationEnabled) == 1 {
	//	ctor, ok := newCollatorMap[collate]
	//	if !ok {
	//		logutil.BgLogger().Warn(
	//			"Unable to get collator by name, use binCollator instead.",
	//			zap.String("name", collate),
	//			zap.Stack("stack"))
	//		return newCollatorMap["utf8mb4_bin"]
	//	}
	//	return ctor
	//}
	return binCollatorInstance
}

// CompatibleCollate checks whether the two collate are the same.
func CompatibleCollate(collate1, collate2 string) bool {
	if (collate1 == "utf8mb4_general_ci" || collate1 == "utf8_general_ci") && (collate2 == "utf8mb4_general_ci" || collate2 == "utf8_general_ci") {
		return true
	} else if (collate1 == "utf8mb4_bin" || collate1 == "utf8_bin") && (collate2 == "utf8mb4_bin" || collate2 == "utf8_bin") {
		return true
	} else if (collate1 == "utf8mb4_unicode_ci" || collate1 == "utf8_unicode_ci") && (collate2 == "utf8mb4_unicode_ci" || collate2 == "utf8_unicode_ci") {
		return true
	} else {
		return collate1 == collate2
	}
}

// IsCICollation returns if the collation is case-sensitive
func IsCICollation(collate string) bool {
	return collate == "utf8_general_ci" || collate == "utf8mb4_general_ci" ||
		collate == "utf8_unicode_ci" || collate == "utf8mb4_unicode_ci"
}
