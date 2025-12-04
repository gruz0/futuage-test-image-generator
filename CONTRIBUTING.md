# Contributing to FutuAge Test Image Generator

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Basic understanding of image processing

### Development Setup

1. **Fork and clone the repository**

```bash
git clone https://github.com/yourusername/futuage-test-image-generator.git
cd futuage-test-image-generator
```

2. **Install dependencies**

```bash
go mod download
```

3. **Build the project**

```bash
go build -o futuage-test-image-gen .
```

4. **Run tests**

```bash
go test ./...
```

5. **Test the CLI**

```bash
./futuage-test-image-gen generate --output ./test-output/
```

## Project Structure

```
futuage-test-image-generator/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── generate.go        # Generate command
│   └── list.go            # List command
├── internal/              # Internal packages
│   ├── config/            # Configuration handling
│   ├── generator/         # Image generation logic
│   ├── manifest/          # Manifest generation
│   └── filesystem/        # File operations
├── configs/               # Configuration files
└── main.go               # Entry point
```

## How to Contribute

### Reporting Bugs

1. Check if the bug is already reported in [Issues](https://github.com/gruz0/futuage-test-image-generator/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)
   - Sample output or screenshots

### Suggesting Features

1. Check existing [Issues](https://github.com/gruz0/futuage-test-image-generator/issues) for similar suggestions
2. Create a new issue with:
   - Clear use case description
   - Proposed solution
   - Alternative approaches considered
   - Impact on existing functionality

### Submitting Code

1. **Create a feature branch**

```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes**

Follow these guidelines:

- Write clear, readable code
- Add comments for complex logic
- Follow Go conventions and best practices
- Keep functions focused and small

3. **Test your changes**

```bash
# Run all tests
go test ./...

# Test the full generation
./futuage-test-image-gen generate --output ./test/
```

4. **Commit your changes**

Use clear commit messages:

```bash
git commit -m "Add feature: support for custom aspect ratios"
```

Follow the conventional commit format:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for adding tests
- `chore:` for maintenance tasks

5. **Push to your fork**

```bash
git push origin feature/your-feature-name
```

6. **Create a Pull Request**

- Provide a clear title and description
- Reference any related issues
- Include screenshots/examples if applicable
- Ensure CI checks pass

## Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golint` for linting
- Keep line length reasonable (~100 characters)

### Code Organization

- Keep packages focused on a single responsibility
- Use meaningful variable and function names
- Prefer composition over inheritance
- Write self-documenting code

### Error Handling

- Always check and handle errors
- Provide context with error messages
- Use `fmt.Errorf` with `%w` for error wrapping

### Testing

- Write tests for new features
- Maintain existing test coverage
- Use table-driven tests where appropriate

## Documentation

- Update README.md for user-facing changes
- Update CHANGELOG.md following Keep a Changelog format
- Add godoc comments for exported functions and types

## Performance

When making changes that affect performance:

- Benchmark before and after
- Document performance characteristics
- Avoid premature optimization
- Profile if investigating slowdowns

## Adding New Features

### Adding a New Platform Target

1. Update `configs/default.json`:

```json
"FACEBOOK_1_1": {
  "platform": "Facebook",
  "dimensions": [1200, 1200],
  "ratio": "1:1",
  "description": "Facebook (square)"
}
```

2. Test the generation:

```bash
./futuage-test-image-gen generate --output ./test/
```

### Adding a New Format

1. Update `internal/generator/encoder.go`
2. Add format to `configs/default.json`
3. Update documentation
4. Add tests

### Adding a New CLI Command

1. Create command file in `cmd/`
2. Register in `cmd/root.go`
3. Add tests
4. Update README

## Release Process

Maintainers will:

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create a git tag
4. Build release binaries
5. Create GitHub release
6. Update documentation

## Questions?

Feel free to:

- Open an issue for discussion
- Join our community channels
- Email the maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🎉
