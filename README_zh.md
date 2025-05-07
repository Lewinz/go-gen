# go-gen

一个灵活的 Go 项目代码生成工具，支持多种生成器和命名规范。

[English](README.md) | [中文](README_zh.md)

## 特性

- MongoDB 模型生成
- 可自定义命名规范
- 基于模板的代码生成
- 跨平台支持（Linux、macOS、Windows）
- 支持多种 CPU 架构（amd64、arm64）

## 快速开始

### 安装

```bash
# 使用 Go 安装
go install github.com/lewinz/go-gen@latest
```

### 基本使用

```bash
# 生成 MongoDB 模型
go-gen model mongo --type user --dir ./internal/model
```

## 详细使用说明

### 命令行选项

```bash
# 必需参数
--type string     模型类型名称（例如：user、product）
--dir string      输出目录

# 可选参数
--template string 模板目录或 Git 仓库 URL（默认：git@github.com:Lewinz/go-gen.git）
--file-style string   文件命名风格（snake|camel|pascal|kebab）（默认为 "snake"）
```

### 命名规范

工具在模板中支持四种命名规范：

- `{{.TypeSnake}}`: 蛇形命名（例如：user_profile）
- `{{.TypeCamel}}`: 驼峰命名（例如：userProfile）
- `{{.TypePascal}}`: 帕斯卡命名（例如：UserProfile）
- `{{.TypeKebab}}`: 短横线命名（例如：user-profile）

### 示例

1. 使用默认命名和模板生成 MongoDB 模型：
```bash
go-gen model mongo --type user --dir ./internal/model
```

2. 使用自定义文件命名：
```bash
go-gen model mongo \
  --type user \
  --dir ./internal/model \
  --file-style camel
```

3. 使用自定义模板：
```bash
go-gen model mongo \
  --type user \
  --dir ./internal/model \
  --template ./template
```

4. 从 Git 模板仓库生成：
```bash
go-gen model mongo \
  --type user \
  --dir ./internal/model \
  --template https://github.com/your-org/go-templates
```

## 模板

### 模板文件

工具使用模板文件（`.tpl`）来生成代码。以下是模板文件与生成文件的对应关系：

```
template/
└── mongo/
    ├── model.tpl      # 生成：{type_snake}.go
    │                  # 示例：user.go、product.go
    │                  # 包含：结构体定义和 CRUD 操作
    │
    └── model_test.tpl # 生成：{type_snake}_test.go
                       # 示例：user_test.go、product_test.go
                       # 包含：模型的单元测试
```

### 模板变量

模板中可用的变量：

- `{{.TypeSnake}}`: 蛇形命名的类型名（例如：user_profile）
- `{{.TypeCamel}}`: 驼峰命名的类型名（例如：userProfile）
- `{{.TypePascal}}`: 帕斯卡命名的类型名（例如：UserProfile）
- `{{.TypeKebab}}`: 短横线命名的类型名（例如：user-profile）
- `{{.PackageName}}`: 生成文件的包名

## 高级用法

### 使用自定义模板

1. 创建模板目录：
```bash
mkdir -p templates/mongo
```

2. 创建模板文件：
```bash
# templates/mongo/model.tpl
package {{.PackageName}}

type {{.TypePascal}} struct {
    ID        string    `bson:"_id,omitempty"`
    CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}
```

3. 生成代码：
```bash
go-gen model mongo --type user --dir ./internal/model --template ./templates
```

### 使用 Git 模板

你可以使用 Git 仓库中的模板：

```bash
go-gen model mongo \
  --type user \
  --dir ./internal/model \
  --template https://github.com/your-org/go-templates
```

工具会：
1. 克隆仓库
2. 在本地缓存
3. 用于代码生成

## 贡献

1. Fork 本仓库
2. 创建你的特性分支（`git checkout -b feature/amazing-feature`）
3. 提交你的更改（`git commit -m '添加一些很棒的特性'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 开启一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。 