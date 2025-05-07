package model

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetModelCmd(t *testing.T) {
	cmd := GetModelCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "model", cmd.Use)
	assert.Equal(t, "Generate model code", cmd.Short)
	assert.Equal(t, "Generate model code for different databases like MongoDB, MySQL, etc.", cmd.Long)
}

func TestModelCmdFlags(t *testing.T) {
	cmd := GetModelCmd()

	// Check if all flags exist
	assert.NotNil(t, cmd.Flag("type"))
	assert.NotNil(t, cmd.Flag("dir"))
	assert.NotNil(t, cmd.Flag("template"))
	assert.NotNil(t, cmd.Flag("file-style"))

	// Check if flags are required by trying to execute mongo subcommand without required flags
	cmd.SetArgs([]string{"mongo"})
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestMongoSubcommand(t *testing.T) {
	cmd := GetModelCmd()

	// Find mongo subcommand
	var mongoCmd *cobra.Command
	for _, subCmd := range cmd.Commands() {
		if subCmd.Use == "mongo" {
			mongoCmd = subCmd
			break
		}
	}

	assert.NotNil(t, mongoCmd)
	assert.Equal(t, "mongo", mongoCmd.Use)
	assert.Equal(t, "Generate MongoDB model code", mongoCmd.Short)
	assert.Equal(t, "Generate MongoDB model code with specified type and naming style.", mongoCmd.Long)
}

func TestFileStyleFlag(t *testing.T) {
	cmd := GetModelCmd()

	// Test default value
	assert.Equal(t, "snake", cmd.Flag("file-style").DefValue)

	// Test valid values
	validStyles := []string{"snake", "camel", "pascal", "kebab"}
	for _, style := range validStyles {
		err := cmd.Flag("file-style").Value.Set(style)
		assert.NoError(t, err)
		assert.Equal(t, style, cmd.Flag("file-style").Value.String())
	}
}

func TestCommandExecution(t *testing.T) {
	cmd := GetModelCmd()

	// Test without required flags
	cmd.SetArgs([]string{"mongo"})
	err := cmd.Execute()
	assert.Error(t, err)

	// Test with required flags
	cmd.SetArgs([]string{
		"mongo",
		"--type", "User",
		"--dir", "./output",
		"--template", "./templates",
	})

	// Note: We can't fully test the execution here because it would require
	// actual file system operations and template processing.
	// This would be better tested with integration tests.
}
