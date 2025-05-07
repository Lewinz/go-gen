package generator

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
		expected    *BaseGenerator
	}{
		{
			name:        "valid generator",
			typeName:    "User",
			outputDir:   "./output",
			templateDir: "./templates",
			fileStyle:   "snake",
			expected: &BaseGenerator{
				Type:        "User",
				OutputDir:   "./output",
				TemplateDir: "./templates",
				FileStyle:   "snake",
			},
		},
		{
			name:        "empty values",
			typeName:    "",
			outputDir:   "",
			templateDir: "",
			fileStyle:   "",
			expected: &BaseGenerator{
				Type:        "",
				OutputDir:   "",
				TemplateDir: "",
				FileStyle:   "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generator := NewBaseGenerator(tc.typeName, tc.outputDir, tc.templateDir, tc.fileStyle)
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
		expectError bool
	}{
		{
			name:        "valid generator",
			typeName:    "User",
			outputDir:   "./output",
			templateDir: "./templates",
			fileStyle:   "snake",
			expectError: false,
		},
		{
			name:        "missing type",
			typeName:    "",
			outputDir:   "./output",
			templateDir: "./templates",
			fileStyle:   "snake",
			expectError: true,
		},
		{
			name:        "missing output dir",
			typeName:    "User",
			outputDir:   "",
			templateDir: "./templates",
			fileStyle:   "snake",
			expectError: true,
		},
		{
			name:        "missing template dir",
			typeName:    "User",
			outputDir:   "./output",
			templateDir: "",
			fileStyle:   "snake",
			expectError: true,
		},
		{
			name:        "all empty",
			typeName:    "",
			outputDir:   "",
			templateDir: "",
			fileStyle:   "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generator := NewBaseGenerator(tc.typeName, tc.outputDir, tc.templateDir, tc.fileStyle)
			err := generator.Validate()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
