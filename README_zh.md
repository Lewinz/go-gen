# go-gen

一个灵活的 Go 项目代码生成工具，支持多种生成器和命名规范。

[English](README.md) | [中文](README_zh.md)

## 特性

- 支持多种代码生成器（MongoDB、MySQL 等）
- 可自定义命名规范
- 基于模板的代码生成
- 可扩展的架构设计

## 安装

```bash
go install github.com/lewinz/go-gen@latest
```

## 使用方法

### 基本用法

```bash
# 生成 MongoDB 模型
go-gen model mongo --type user --dir ./internal/model --template ./template
```

### 命令行选项

```bash
# 必需参数
--type string     模型类型名称（例如：user、product）
--dir string      输出目录
--template string 模板目录

# 可选参数
--file-style string   文件命名风格（snake|camel|pascal|kebab）（默认 "snake"）
```

### 命名规范

工具在模板中支持四种命名规范：

- `{{.TypeSnake}}`: 下划线命名（例如：user_profile）
- `{{.TypeCamel}}`: 驼峰命名（例如：userProfile）
- `{{.TypePascal}}`: 帕斯卡命名（例如：UserProfile）
- `{{.TypeKebab}}`: 短横线命名（例如：user-profile）

### 使用示例

1. 使用默认命名生成 MongoDB 模型：
```bash
go-gen model mongo --type user --dir ./internal/model --template ./template
```

2. 使用自定义文件命名生成：
```bash
go-gen model mongo \
  --type user \
  --dir ./internal/model \
  --template ./template \
  --file-style camel
```

## 模板说明

### 模板文件

工具使用模板文件（`.tpl`）来生成代码。以下是模板文件与生成文件的对应关系：

```
template/
└── mongo/
    ├── model.tpl      # 生成: {type_snake}.go
    │                  # 示例: user.go, product.go
    │                  # 包含: 结构体定义和 CRUD 操作
    │
    └── model_test.tpl # 生成: {type_snake}_test.go
                       # 示例: user_test.go, product_test.go
                       # 包含: 模型的单元测试
```

### 模板变量

模板中可以使用以下变量：

- `{{.TypeSnake}}`: 下划线命名（例如：user_profile）
- `{{.TypeCamel}}`: 驼峰命名（例如：userProfile）
- `{{.TypePascal}}`: 帕斯卡命名（例如：UserProfile）
- `{{.TypeKebab}}`: 短横线命名（例如：user-profile）
- `{{.PackageName}}`: 生成文件的包名

### 示例输出

模板输入：
```go
// {{.TypePascal}} 是 MongoDB 模型
type {{.TypePascal}} struct {
    ID        string    `bson:"_id,omitempty"`
    CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}
```

使用类型 "user_profile" 时，生成：
```go
// UserProfile 是 MongoDB 模型
type UserProfile struct {
    ID        string    `bson:"_id,omitempty"`
    CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}
```

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支（`git checkout -b feature/amazing-feature`）
3. 提交您的更改（`git commit -m '添加某个特性'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 开启一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。 