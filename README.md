# go-gen

A flexible code generation tool for Go projects that supports multiple generators and naming conventions.

[English](README.md) | [中文](README_zh.md)

## Features

- Multiple code generators support (MongoDB, MySQL, etc.)
- Customizable naming conventions
- Template-based code generation
- Extensible architecture

## Installation

```bash
go install github.com/lewinz/go-gen@latest
```

## Usage

### Basic Usage

```bash
# Generate MongoDB model
go-gen model mongo --type user --dir ./internal/model --template ./template
```

### Command Options

```bash
# Required flags
--type string     Model type name (e.g., user, product)
--dir string      Output directory
--template string Template directory

# Optional flags
--file-style string   File naming style (snake|camel|pascal|kebab) (default "snake")
```

### Naming Conventions

The tool supports four naming conventions in templates:

- `{{.TypeSnake}}`: snake_case (e.g., user_profile)
- `{{.TypeCamel}}`: camelCase (e.g., userProfile)
- `{{.TypePascal}}`: PascalCase (e.g., UserProfile)
- `{{.TypeKebab}}`: kebab-case (e.g., user-profile)

### Examples

1. Generate a MongoDB model with default naming:
```bash
go-gen model mongo --type user --dir ./internal/model --template ./template
```

2. Generate with custom file naming:
```bash
go-gen model mongo \
  --type user \
  --dir ./internal/model \
  --template ./template \
  --file-style camel
```

## Templates

### Template Files

The tool uses template files (`.tpl`) to generate code. Here's how the template files map to the generated files:

```
template/
└── mongo/
    ├── model.tpl      # Generates: {type_snake}.go
    │                  # Example: user.go, product.go
    │                  # Contains: struct definition and CRUD operations
    │
    └── model_test.tpl # Generates: {type_snake}_test.go
                       # Example: user_test.go, product_test.go
                       # Contains: unit tests for the model
```

### Template Variables

The following variables are available in templates:

- `{{.TypeSnake}}`: Type name in snake_case (e.g., user_profile)
- `{{.TypeCamel}}`: Type name in camelCase (e.g., userProfile)
- `{{.TypePascal}}`: Type name in PascalCase (e.g., UserProfile)
- `{{.TypeKebab}}`: Type name in kebab-case (e.g., user-profile)
- `{{.PackageName}}`: Package name for the generated file

### Example Output

For a template input:
```go
// {{.TypePascal}} is a MongoDB model
type {{.TypePascal}} struct {
    ID        string    `bson:"_id,omitempty"`
    CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}
```

With type "user_profile", it generates:
```go
// UserProfile is a MongoDB model
type UserProfile struct {
    ID        string    `bson:"_id,omitempty"`
    CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.