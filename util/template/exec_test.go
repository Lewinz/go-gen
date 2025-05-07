package template

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealCommander(t *testing.T) {
	commander := &RealCommander{}
	cmd := commander.Command("echo", "test")
	assert.NotNil(t, cmd)
	assert.Equal(t, filepath.Base(cmd.Path), "echo")
	assert.Equal(t, []string{"test"}, cmd.Args[1:])
}

func TestDefaultCommander(t *testing.T) {
	assert.NotNil(t, defaultCommander)
	cmd := defaultCommander.Command("echo", "test")
	assert.NotNil(t, cmd)
	assert.Equal(t, filepath.Base(cmd.Path), "echo")
	assert.Equal(t, []string{"test"}, cmd.Args[1:])
}
