package naming

import (
	"strings"
	"unicode"
)

// Style defines naming style
type Style string

const (
	StyleSnake  Style = "snake"  // user_model
	StyleCamel  Style = "camel"  // userModel
	StylePascal Style = "pascal" // UserModel
	StyleKebab  Style = "kebab"  // user-model
)

// Converter is a naming style converter
type Converter struct {
	style Style
}

// NewConverter creates a new naming style converter
func NewConverter(style Style) *Converter {
	return &Converter{style: style}
}

// Convert converts input string to specified style
func (c *Converter) Convert(input string) string {
	if input == "" {
		return ""
	}

	// Split input string into words
	words := splitIntoWords(input)

	switch c.style {
	case StyleSnake:
		return strings.Join(words, "_")
	case StyleCamel:
		return toCamelCase(words)
	case StylePascal:
		return toPascalCase(words)
	case StyleKebab:
		return strings.Join(words, "-")
	default:
		return input
	}
}

// splitIntoWords splits input string into words
func splitIntoWords(s string) []string {
	if s == "" {
		return []string{}
	}

	// Handle existing separators
	s = strings.ReplaceAll(s, "-", "_")

	var words []string
	var current strings.Builder

	for i, r := range s {
		// Save current word when encountering separator
		if r == '_' || unicode.IsSpace(r) {
			if current.Len() > 0 {
				words = append(words, strings.ToLower(current.String()))
				current.Reset()
			}
			continue
		}

		// Save current word when encountering uppercase letter (except first character)
		if i > 0 && unicode.IsUpper(r) && !unicode.IsUpper(rune(s[i-1])) {
			if current.Len() > 0 {
				words = append(words, strings.ToLower(current.String()))
				current.Reset()
			}
		}

		current.WriteRune(r)
	}

	// Save last word
	if current.Len() > 0 {
		words = append(words, strings.ToLower(current.String()))
	}

	return words
}

// toCamelCase converts word list to camel case
func toCamelCase(words []string) string {
	if len(words) == 0 {
		return ""
	}

	result := words[0]
	for i := 1; i < len(words); i++ {
		result += strings.Title(words[i])
	}
	return result
}

// toPascalCase converts word list to pascal case
func toPascalCase(words []string) string {
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for _, word := range words {
		result.WriteString(strings.Title(word))
	}
	return result.String()
}
