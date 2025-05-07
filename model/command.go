package model

import (
	"github.com/lewinz/go-gen/generator"
	"github.com/lewinz/go-gen/model/mongo"
	"github.com/spf13/cobra"
)

var (
	// Command line arguments
	typeName    string
	outputDir   string
	templateDir string
	fileStyle   string

	// modelCmd is the model generation command
	modelCmd = &cobra.Command{
		Use:   "model",
		Short: "Generate model code",
		Long:  `Generate model code for different databases like MongoDB, MySQL, etc.`,
	}

	// mongoCmd is the MongoDB model generation command
	mongoCmd = &cobra.Command{
		Use:   "mongo",
		Short: "Generate MongoDB model code",
		Long:  `Generate MongoDB model code with specified type and naming style.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create base generator
			base := generator.NewBaseGenerator(typeName, outputDir, templateDir, fileStyle)

			// Create MongoDB generator
			generator := mongo.NewMongoGenerator(base)

			// Execute generation
			return generator.Generate()
		},
	}
)

func init() {
	// Add MongoDB subcommand
	modelCmd.AddCommand(mongoCmd)

	// Add common parameters
	modelCmd.PersistentFlags().StringVar(&typeName, "type", "", "Model type name (required)")
	modelCmd.PersistentFlags().StringVar(&outputDir, "dir", "", "Output directory (required)")
	modelCmd.PersistentFlags().StringVar(&templateDir, "template", "", "Template directory (required)")
	modelCmd.PersistentFlags().StringVar(&fileStyle, "file-style", "snake", "File naming style (snake|camel|pascal|kebab)")

	// Set required parameters
	if err := modelCmd.MarkPersistentFlagRequired("type"); err != nil {
		panic(err)
	}
	if err := modelCmd.MarkPersistentFlagRequired("dir"); err != nil {
		panic(err)
	}
	if err := modelCmd.MarkPersistentFlagRequired("template"); err != nil {
		panic(err)
	}
}

// GetModelCmd returns the model generation command
func GetModelCmd() *cobra.Command {
	return modelCmd
}
