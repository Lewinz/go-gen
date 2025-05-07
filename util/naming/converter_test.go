package naming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConverter(t *testing.T) {
	testCases := []struct {
		name     string
		style    Style
		expected Style
	}{
		{"snake style", StyleSnake, StyleSnake},
		{"camel style", StyleCamel, StyleCamel},
		{"pascal style", StylePascal, StylePascal},
		{"kebab style", StyleKebab, StyleKebab},
		{"invalid style", "invalid", "invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			converter := NewConverter(tc.style)
			assert.NotNil(t, converter)
			assert.Equal(t, tc.expected, converter.style)
		})
	}
}

func TestConvert(t *testing.T) {
	testCases := []struct {
		name     string
		style    Style
		input    string
		expected string
	}{
		// Snake case tests
		{"snake from pascal", StyleSnake, "UserModel", "user_model"},
		{"snake from camel", StyleSnake, "userModel", "user_model"},
		{"snake from kebab", StyleSnake, "user-model", "user_model"},
		{"snake from mixed", StyleSnake, "User_Model", "user_model"},
		{"snake with number", StyleSnake, "User2Model", "user2_model"},
		{"snake with space", StyleSnake, "User Model", "user_model"},
		{"snake already snake", StyleSnake, "user_model", "user_model"},

		// Camel case tests
		{"camel from pascal", StyleCamel, "UserModel", "userModel"},
		{"camel from snake", StyleCamel, "user_model", "userModel"},
		{"camel from kebab", StyleCamel, "user-model", "userModel"},
		{"camel from mixed", StyleCamel, "User_Model", "userModel"},
		{"camel with number", StyleCamel, "user2_model", "user2Model"},
		{"camel with space", StyleCamel, "user model", "userModel"},
		{"camel already camel", StyleCamel, "userModel", "userModel"},

		// Pascal case tests
		{"pascal from snake", StylePascal, "user_model", "UserModel"},
		{"pascal from camel", StylePascal, "userModel", "UserModel"},
		{"pascal from kebab", StylePascal, "user-model", "UserModel"},
		{"pascal from mixed", StylePascal, "User_Model", "UserModel"},
		{"pascal with number", StylePascal, "user2_model", "User2Model"},
		{"pascal with space", StylePascal, "user model", "UserModel"},
		{"pascal already pascal", StylePascal, "UserModel", "UserModel"},

		// Kebab case tests
		{"kebab from pascal", StyleKebab, "UserModel", "user-model"},
		{"kebab from snake", StyleKebab, "user_model", "user-model"},
		{"kebab from camel", StyleKebab, "userModel", "user-model"},
		{"kebab from mixed", StyleKebab, "User_Model", "user-model"},
		{"kebab with number", StyleKebab, "User2Model", "user2-model"},
		{"kebab with space", StyleKebab, "User Model", "user-model"},
		{"kebab already kebab", StyleKebab, "user-model", "user-model"},

		// Edge cases
		{"empty string", StyleSnake, "", ""},
		{"single char", StyleSnake, "A", "a"},
		{"single word", StyleSnake, "user", "user"},
		{"invalid style", "invalid", "UserModel", "UserModel"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			converter := NewConverter(tc.style)
			result := converter.Convert(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSplitIntoWords(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{"pascal case", "UserModel", []string{"user", "model"}},
		{"camel case", "userModel", []string{"user", "model"}},
		{"snake case", "user_model", []string{"user", "model"}},
		{"kebab case", "user-model", []string{"user", "model"}},
		{"mixed case", "User_Model", []string{"user", "model"}},
		{"with number", "User2Model", []string{"user2", "model"}},
		{"with space", "User Model", []string{"user", "model"}},
		{"empty string", "", []string{}},
		{"single char", "A", []string{"a"}},
		{"single word", "user", []string{"user"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := splitIntoWords(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected string
	}{
		{"simple", []string{"user", "model"}, "userModel"},
		{"with number", []string{"user2", "model"}, "user2Model"},
		{"single word", []string{"user"}, "user"},
		{"empty", []string{}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toCamelCase(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToPascalCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected string
	}{
		{"simple", []string{"user", "model"}, "UserModel"},
		{"with number", []string{"user2", "model"}, "User2Model"},
		{"single word", []string{"user"}, "User"},
		{"empty", []string{}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toPascalCase(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
