package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lewinz/go-gen/util/naming"
)

const (
	cacheDir = ".go-gen"
)

// Engine 模板处理引擎
type Engine struct {
	fileStyle naming.Style
}

// NewEngine 创建一个模板处理引擎
func NewEngine(fileStyle naming.Style) *Engine {
	return &Engine{
		fileStyle: fileStyle,
	}
}

// TemplateData 模板数据
type TemplateData struct {
	Type        string // 模型类型
	TypeSnake   string // 蛇形命名
	TypeCamel   string // 驼峰命名
	TypePascal  string // 帕斯卡命名
	TypeKebab   string // 短横线命名
	PackageName string // 包名
}

// Generate 生成代码文件
func (e *Engine) Generate(templateDir, outputDir, typeName string) error {
	// 如果是 git 仓库，先克隆或使用缓存
	if isGitRepo(templateDir) {
		cachedDir, err := getCachedTemplate(templateDir)
		if err != nil {
			return fmt.Errorf("get cached template: %w", err)
		}
		templateDir = cachedDir
	}

	// 准备模板数据
	data := &TemplateData{
		Type:        typeName,
		TypeSnake:   naming.NewConverter(naming.StyleSnake).Convert(typeName),
		TypeCamel:   naming.NewConverter(naming.StyleCamel).Convert(typeName),
		TypePascal:  naming.NewConverter(naming.StylePascal).Convert(typeName),
		TypeKebab:   naming.NewConverter(naming.StyleKebab).Convert(typeName),
		PackageName: filepath.Base(outputDir),
	}

	// 遍历模板目录
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理 .tpl 文件
		if !strings.HasSuffix(info.Name(), ".tpl") {
			return nil
		}

		// 读取模板文件
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("parse template %s: %w", path, err)
		}

		// 生成输出文件名
		outputName := strings.TrimSuffix(info.Name(), ".tpl") // 去掉 .tpl 后缀
		// 组合类型名和模板名，并确保它们之间有分隔符
		outputName = typeName + "_" + outputName // 例如：user + _ + model = user_model
		// 根据指定的命名风格转换
		converter := naming.NewConverter(e.fileStyle)
		outputName = converter.Convert(outputName) // 例如：user_model -> userModel
		outputName = outputName + ".go"            // 添加 .go 后缀
		outputPath := filepath.Join(outputDir, outputName)

		// 创建输出文件
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("create output file %s: %w", outputPath, err)
		}
		defer outputFile.Close()

		// 渲染模板
		if err := tmpl.Execute(outputFile, data); err != nil {
			return fmt.Errorf("execute template %s: %w", path, err)
		}

		return nil
	})
}

// isGitRepo checks if the path is a git repository URL
func isGitRepo(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "git@")
}

// getRepoHash gets the latest commit hash of a repository
func getRepoHash(repoURL string) (string, error) {
	cmd := defaultCommander.Command("git", "ls-remote", repoURL, "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git ls-remote: %w", err)
	}
	parts := strings.Fields(string(output))
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid git ls-remote output")
	}
	return parts[0], nil
}

// getCurrentHash gets the current commit hash of a repository
func getCurrentHash(repoPath string) (string, error) {
	cmd := defaultCommander.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// getCachedTemplate gets or creates a cached template
func getCachedTemplate(repoURL string) (string, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	// Create cache directory if not exists
	cachePath := filepath.Join(homeDir, cacheDir)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}

	// Get repository name and hash
	repoName := getRepoName(repoURL)
	repoHash, err := getRepoHash(repoURL)
	if err != nil {
		return "", fmt.Errorf("get repo hash: %w", err)
	}

	// Check if cached version exists
	cachedPath := filepath.Join(cachePath, repoName)
	if _, err := os.Stat(cachedPath); err == nil {
		// Check if hash matches
		currentHash, err := getCurrentHash(cachedPath)
		if err == nil && currentHash == repoHash {
			return cachedPath, nil
		}
		// Remove old cache if hash doesn't match
		os.RemoveAll(cachedPath)
	}

	// Clone repository
	cmd := defaultCommander.Command("git", "clone", repoURL, cachedPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone: %w", err)
	}

	return cachedPath, nil
}

// getRepoName extracts repository name from URL
func getRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	lastPart := parts[len(parts)-1]
	return strings.TrimSuffix(lastPart, ".git")
}
