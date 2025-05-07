package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseGenerator(t *testing.T) {
	testCases := []struct {
		name        string
		typeName    string
		outputDir   string
		templateDir string
		fileStyle   string
		varStyle    string
		expected    *BaseGenerator
	}{
		{
			name:        "valid generator",
			typeName:    "User",
			outputDir:   "./output",
			templateDir: "./templates",
			fileStyle:   "snake",
			varStyle:    "camel",
			expected: &BaseGenerator{
				Type:        "User",
				OutputDir:   "./output",
				TemplateDir: "./templates",
				FileStyle:   "snake",
				VarStyle:    "camel",
			},
		},
		{
			name:        "empty values",
			typeName:    "",
			outputDir:   "",
			templateDir: "",
			fileStyle:   "",
			varStyle:    "",
			expected: &BaseGenerator{
				Type:        "",
				OutputDir:   "",
				TemplateDir: "",
				FileStyle:   "",
				VarStyle:    "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generator := NewBaseGenerator(tc.typeName, tc.outputDir, tc.templateDir, tc.fileStyle, tc.varStyle)
			assert.Equal(t, tc.expected, generator)
		})
	}
}

func TestBaseGeneratorValidate(t *testing.T) {
	testCases := []struct {
		name        string
		typeName    string
		outputDir   string
		templateDir string
		fileStyle   string
		varStyle    string
		expectError bool
	}{
		{
			name:        "valid generator",
			typeName:    "User",
			outputDir:   "./output",
			templateDir: "./templates",
			fileStyle:   "snake",
			varStyle:    "camel",
			expectError: false,
		},
		{
			name:        "missing type",
			typeName:    "",
			outputDir:   "./output",
			templateDir: "./templates",
			fileStyle:   "snake",
			varStyle:    "camel",
			expectError: true,
		},
		{
			name:        "missing output dir",
			typeName:    "User",
			outputDir:   "",
			templateDir: "./templates",
			fileStyle:   "snake",
			varStyle:    "camel",
			expectError: true,
		},
		{
			name:        "missing template dir",
			typeName:    "User",
			outputDir:   "./output",
			templateDir: "",
			fileStyle:   "snake",
			varStyle:    "camel",
			expectError: true,
		},
		{
			name:        "all empty",
			typeName:    "",
			outputDir:   "",
			templateDir: "",
			fileStyle:   "",
			varStyle:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generator := NewBaseGenerator(tc.typeName, tc.outputDir, tc.templateDir, tc.fileStyle, tc.varStyle)
			err := generator.Validate()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
