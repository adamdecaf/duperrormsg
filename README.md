# duperrormsg

`duperrormsg` is a Go linter that detects duplicate error messages across your codebase.
Designed to help avoid the common problem of having identical error messages in different
code paths, which makes debugging difficult.

## Overview

When different error conditions in your code use the same error message, it becomes harder
to identify which part of the code generated an error. This linter helps you catch these
issues early by detecting duplicate error messages.

## Installation

```bash
go install github.com/adamdecaf/duperrormsg@latest
```

## Usage

```bash
# Run on current directory
duperrormsg ./...

# Run on specific packages
duperrormsg github.com/your/package/...

# Using with go vet
go vet -vettool=$(which duperrormsg) ./...
```

## Features

The linter detects duplicate error messages created through various methods:

- Standard library error creation:
  - `errors.New("message")`
  - `fmt.Errorf("message: %v", err)`

- Standard library logging:
  - `log.Printf("error message")`
  - `log.Fatalf("error message")`

- Custom error constructors:
  - Functions starting with `New` and containing `Error`
  - Other common error construction patterns

- Structured logging libraries:
  - Supports chained method calls like `logger.Info().Logf("message")`
  - Works with the moov-io/base/log package

## Error Normalization

The linter normalizes error messages to detect duplicates even when the format specifiers differ:

```go
fmt.Errorf("user %s not found", name)
fmt.Errorf("user %v not found", name)  // Detected as duplicate
```

## Examples

Here are some examples of issues that the linter will detect:

```go
// Duplicate error messages with different error constructors
errors.New("validation failed")
fmt.Errorf("validation failed")

// Duplicate error messages in different functions
func validateA() error {
    return errors.New("invalid format")
}

func validateB() error {
    return errors.New("invalid format")  // Duplicate!
}

// Duplicate error messages with different format specifiers
fmt.Errorf("user %s not found", username)
fmt.Errorf("user %v not found", id)  // Detected as duplicate
```

## Configuration

Currently, the linter doesn't support configuration, but future versions may include:

- Exclude patterns for certain errors
- Severity levels
- Scope limitations

## Contributing

Contributions are welcome! Here's how you can help:

1. File issues for bugs or feature requests
2. Submit pull requests for improvements
3. Help with documentation or tests

## License

[Apache 2 License](LICENSE)
