package template

import (
	"os/exec"
)

// Commander 定义命令执行器接口
type Commander interface {
	Command(name string, args ...string) *exec.Cmd
}

// RealCommander 实现真实的命令执行
type RealCommander struct{}

// Command 执行真实的命令
func (c *RealCommander) Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// defaultCommander 是默认的命令执行器
var defaultCommander Commander = &RealCommander{}
