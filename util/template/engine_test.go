package template

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/lewinz/go-gen/util/naming"
	"github.com/stretchr/testify/assert"
)

// MockCommander 实现命令执行器接口用于测试
type MockCommander struct{}

// Command 返回一个模拟的命令
func (c *MockCommander) Command(name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		os.Exit(2)
	}

	cmd := args[0]
	switch cmd {
	case "git":
		switch args[1] {
		case "ls-remote":
			// 检查是否是有效的仓库 URL
			if args[2] == "invalid-url" || args[2] == "" {
				os.Exit(1)
			}
			os.Stdout.Write([]byte("abcdef1234567890 HEAD"))
		case "rev-parse":
			os.Stdout.Write([]byte("abcdef1234567890"))
		case "clone":
			// 模拟克隆成功
		}
	}
	os.Exit(0)
}

func TestNewEngine(t *testing.T) {
	engine := NewEngine(naming.StyleCamel)
	assert.NotNil(t, engine)
	assert.Equal(t, naming.StyleCamel, engine.fileStyle)
}

func TestGenerate(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "go-gen-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建模板目录
	templateDir := filepath.Join(tempDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	assert.NoError(t, err)

	// 创建模板文件
	templateFile := filepath.Join(templateDir, "model.tpl")
	err = os.WriteFile(templateFile, []byte(`package {{.PackageName}}

type {{.TypePascal}} struct {
	ID string
}`), 0644)
	assert.NoError(t, err)

	// 创建输出目录
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// 测试不同命名风格
	testCases := []struct {
		name     string
		style    naming.Style
		expected string
	}{
		{"snake", naming.StyleSnake, "user.go"},
		{"camel", naming.StyleCamel, "user.go"},
		{"pascal", naming.StylePascal, "User.go"},
		{"kebab", naming.StyleKebab, "user.go"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建引擎
			engine := NewEngine(tc.style)

			// 生成代码
			err := engine.Generate(templateDir, outputDir, "user")
			assert.NoError(t, err)

			// 检查文件是否存在
			outputFile := filepath.Join(outputDir, tc.expected)
			_, err = os.Stat(outputFile)
			assert.NoError(t, err)

			// 检查文件内容
			content, err := os.ReadFile(outputFile)
			assert.NoError(t, err)
			assert.Contains(t, string(content), "type User struct")
		})
	}
}

func TestGenerateWithInvalidTemplate(t *testing.T) {
	engine := NewEngine(naming.StyleSnake)
	err := engine.Generate("invalid/path", "output", "user")
	assert.Error(t, err)
}

func TestGenerateWithInvalidOutput(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "go-gen-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建模板目录
	templateDir := filepath.Join(tempDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	assert.NoError(t, err)

	// 创建模板文件
	templateFile := filepath.Join(templateDir, "model.tpl")
	err = os.WriteFile(templateFile, []byte(`invalid template`), 0644)
	assert.NoError(t, err)

	// 测试无效输出目录
	engine := NewEngine(naming.StyleSnake)
	err = engine.Generate(templateDir, "/invalid/output", "user")
	assert.Error(t, err)
}

func TestIsGitRepo(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{"http", "http://github.com/user/repo", true},
		{"https", "https://github.com/user/repo", true},
		{"ssh", "git@github.com:user/repo", true},
		{"local", "/path/to/repo", false},
		{"relative", "./repo", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isGitRepo(tc.path)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetRepoName(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected string
	}{
		{"http", "http://github.com/user/repo", "repo"},
		{"https", "https://github.com/user/repo", "repo"},
		{"ssh", "git@github.com:user/repo", "repo"},
		{"with-git", "https://github.com/user/repo.git", "repo"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getRepoName(tc.url)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetCachedTemplate(t *testing.T) {
	// 保存原始的命令执行器
	oldCommander := defaultCommander
	defer func() { defaultCommander = oldCommander }()

	// 替换为 mock 版本
	defaultCommander = &MockCommander{}

	// 创建临时目录作为 home 目录
	homeDir, err := os.MkdirTemp("", "go-gen-home-*")
	assert.NoError(t, err)
	defer os.RemoveAll(homeDir)

	// 设置环境变量
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	testCases := []struct {
		name        string
		repoURL     string
		setupCache  bool
		expectError bool
	}{
		{
			name:        "new repo",
			repoURL:     "https://github.com/user/repo",
			setupCache:  false,
			expectError: false,
		},
		{
			name:        "existing repo same hash",
			repoURL:     "https://github.com/user/repo",
			setupCache:  true,
			expectError: false,
		},
		{
			name:        "invalid repo url",
			repoURL:     "invalid-url",
			setupCache:  false,
			expectError: true,
		},
		{
			name:        "empty repo url",
			repoURL:     "",
			setupCache:  false,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupCache {
				// 创建缓存目录
				cacheDir := filepath.Join(homeDir, ".go-gen", "repo")
				err := os.MkdirAll(cacheDir, 0755)
				assert.NoError(t, err)

				// 创建 .git 目录
				gitDir := filepath.Join(cacheDir, ".git")
				err = os.MkdirAll(gitDir, 0755)
				assert.NoError(t, err)

				// 创建 HEAD 文件
				headFile := filepath.Join(gitDir, "HEAD")
				err = os.WriteFile(headFile, []byte("ref: refs/heads/main"), 0644)
				assert.NoError(t, err)
			}

			result, err := getCachedTemplate(tc.repoURL)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestGenerateWithGitTemplate(t *testing.T) {
	// 保存原始的命令执行器
	oldCommander := defaultCommander
	defer func() { defaultCommander = oldCommander }()

	// 替换为 mock 版本
	defaultCommander = &MockCommander{}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "go-gen-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建输出目录
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// 创建缓存目录
	homeDir, err := os.MkdirTemp("", "go-gen-home-*")
	assert.NoError(t, err)
	defer os.RemoveAll(homeDir)

	// 设置环境变量
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	// 创建缓存目录结构
	cacheDir := filepath.Join(homeDir, ".go-gen", "repo")
	err = os.MkdirAll(cacheDir, 0755)
	assert.NoError(t, err)

	// 创建 .git 目录
	gitDir := filepath.Join(cacheDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	assert.NoError(t, err)

	// 创建 HEAD 文件
	headFile := filepath.Join(gitDir, "HEAD")
	err = os.WriteFile(headFile, []byte("ref: refs/heads/main"), 0644)
	assert.NoError(t, err)

	// 创建模板文件
	templateFile := filepath.Join(cacheDir, "model.tpl")
	err = os.WriteFile(templateFile, []byte(`package {{.PackageName}}

type {{.TypePascal}} struct {
	ID string
}`), 0644)
	assert.NoError(t, err)

	// 测试使用 git 仓库作为模板
	engine := NewEngine(naming.StyleSnake)
	err = engine.Generate("https://github.com/user/repo", outputDir, "user")
	assert.NoError(t, err)

	// 验证输出文件是否存在
	outputFile := filepath.Join(outputDir, "user.go")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)

	// 验证文件内容
	content, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "type User struct")
}

func TestGenerateWithInvalidTemplateContent(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "go-gen-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建模板目录
	templateDir := filepath.Join(tempDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	assert.NoError(t, err)

	// 创建无效的模板文件
	templateFile := filepath.Join(templateDir, "model.tpl")
	err = os.WriteFile(templateFile, []byte(`{{.InvalidField}}`), 0644)
	assert.NoError(t, err)

	// 创建输出目录
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// 测试生成
	engine := NewEngine(naming.StyleSnake)
	err = engine.Generate(templateDir, outputDir, "user")
	assert.Error(t, err)
}

func TestGenerateWithEmptyTemplate(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "go-gen-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建模板目录
	templateDir := filepath.Join(tempDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	assert.NoError(t, err)

	// 创建空模板文件
	templateFile := filepath.Join(templateDir, "model.tpl")
	err = os.WriteFile(templateFile, []byte(``), 0644)
	assert.NoError(t, err)

	// 创建输出目录
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// 测试生成
	engine := NewEngine(naming.StyleSnake)
	err = engine.Generate(templateDir, outputDir, "user")
	assert.NoError(t, err)

	// 验证输出文件是否存在
	outputFile := filepath.Join(outputDir, "user.go")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}
