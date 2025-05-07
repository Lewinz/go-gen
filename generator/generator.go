package generator

import "fmt"

// Generator defines the interface for code generators
type Generator interface {
	// Generate executes the code generation
	Generate() error
	// Validate validates the generation parameters
	Validate() error
}

// BaseGenerator provides the basic implementation of a generator
type BaseGenerator struct {
	Type        string // Model type
	OutputDir   string // Output directory
	TemplateDir string // Template directory
	FileStyle   string // File naming style
}

// NewBaseGenerator creates a new base generator
func NewBaseGenerator(typeName, outputDir, templateDir, fileStyle string) *BaseGenerator {
	return &BaseGenerator{
		Type:        typeName,
		OutputDir:   outputDir,
		TemplateDir: templateDir,
		FileStyle:   fileStyle,
	}
}

// Validate implements basic parameter validation
func (g *BaseGenerator) Validate() error {
	if g.Type == "" {
		return fmt.Errorf("type is required")
	}
	if g.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}
	if g.TemplateDir == "" {
		return fmt.Errorf("template directory is required")
	}
	return nil
}
