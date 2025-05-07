package mongo

import (
	"fmt"
	"os"
	"strings"

	"github.com/lewinz/go-gen/generator"
	"github.com/lewinz/go-gen/util/naming"
	"github.com/lewinz/go-gen/util/template"
)

// MongoGenerator is a MongoDB model generator
type MongoGenerator struct {
	*generator.BaseGenerator
	engine *template.Engine
}

// NewMongoGenerator creates a new MongoDB generator
func NewMongoGenerator(base *generator.BaseGenerator) *MongoGenerator {
	// If file style is not specified, use snake case by default
	if base.FileStyle == "" {
		base.FileStyle = "snake"
	}
	return &MongoGenerator{
		BaseGenerator: base,
		engine:        template.NewEngine(naming.Style(base.FileStyle)),
	}
}

// Generate implements MongoDB model generation
func (g *MongoGenerator) Generate() error {
	if err := g.Validate(); err != nil {
		return err
	}

	// Ensure output directory exists
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Generate code using template engine
	return g.engine.Generate(g.TemplateDir, g.OutputDir, g.Type)
}

// Validate implements MongoDB-specific parameter validation
func (g *MongoGenerator) Validate() error {
	if err := g.BaseGenerator.Validate(); err != nil {
		return err
	}

	// Validate naming styles
	if !isValidStyle(g.FileStyle) {
		return fmt.Errorf("invalid file style: %s", g.FileStyle)
	}

	// Validate template directory
	if !isValidTemplatePath(g.TemplateDir) {
		return fmt.Errorf("invalid template path: %s", g.TemplateDir)
	}

	return nil
}

// isValidStyle checks if the naming style is valid
func isValidStyle(style string) bool {
	switch style {
	case "snake", "camel", "pascal", "kebab":
		return true
	default:
		return false
	}
}

// isValidTemplatePath checks if the template path is valid
func isValidTemplatePath(path string) bool {
	// Check if it's a git repository URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "git@") {
		return true
	}

	// Check if it's a local path
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
