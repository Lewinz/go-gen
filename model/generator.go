package model

import "fmt"

// Generator 定义了代码生成器的接口
type Generator interface {
	// Generate 执行代码生成
	Generate() error
	// Validate 验证生成参数
	Validate() error
}

// BaseGenerator 提供了生成器的基本实现
type BaseGenerator struct {
	Type        string // 模型类型
	OutputDir   string // 输出目录
	TemplateDir string // 模板目录
	FileStyle   string // 文件命名风格
	VarStyle    string // 变量命名风格
}

// NewBaseGenerator 创建一个基础生成器
func NewBaseGenerator(typeName, outputDir, templateDir, fileStyle, varStyle string) *BaseGenerator {
	return &BaseGenerator{
		Type:        typeName,
		OutputDir:   outputDir,
		TemplateDir: templateDir,
		FileStyle:   fileStyle,
		VarStyle:    varStyle,
	}
}

// Validate 实现基本的参数验证
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
