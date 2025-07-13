# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
Blindspot is a Go-based finite state machine generator that creates state transition diagrams from JSON rule definitions. The tool helps identify blind spots in state transitions and supports multiple output formats.

## Common Commands

### Building and Running
```bash
# Build the CLI
go build -o blindspot cmd/cli/blindspot/main.go

# Run with input file
./blindspot -input example/stringlist.json -output mermaid

# Run with iteration limit (safe mode)
./blindspot -input example/stringlist.json -output mermaid --limit 1000

# Install globally
go install github.com/yuukiiwai/blindspot/cmd/cli/blindspot@latest
```

### Testing
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/std-impl/stringlist/

# Run tests with verbose output
go test -v ./...
```

### Development
```bash
# Format code
go fmt ./...

# Run linter (if available)
go vet ./...

# Check module dependencies
go mod tidy
```

## Architecture

### Core Components
- **pkg/core/**: Business logic and interface definitions
  - `Generator`: Core state machine generation engine with BFS traversal
  - `Node`: Interface for state representation with ID and resource management
  - `Edge`: Represents state transitions with rules
  - `EdgeRule`: Defines transition conditions and effects
  - `Formatter`: Interface for output generation (mermaid, dot, visjs)
  - `Parser`: Interface for input parsing

- **pkg/std-impl/**: Concrete implementations
  - `output/`: Output formatters (mermaid.go, dot.go, visjs.go)
  - `stringlist/`: String list-based parser implementation

### Key Design Patterns
- Interface-based architecture for extensibility
- Clean separation between core logic and implementations
- Factory pattern for formatter and parser creation
- BFS algorithm for state space exploration with loop detection
- Iteration limit safety mechanism to prevent infinite loops

### Input Format
JSON structure with:
- `start_resources`: Initial state resources
- `edge_rules`: Array of transition rules with:
  - `name`: Rule identifier
  - `action`: "create" or "delete"
  - `rule`: Resources to modify
  - `fire_condition`: Prerequisites for rule activation
  - `block_condition`: Conditions that prevent rule execution

### Output Formats
- **mermaid**: Mermaid.js flowchart format
- **dot**: Graphviz DOT format (requires github.com/goccy/go-graphviz)
- **visjs**: vis.js network format

## Code Style Guidelines
- Japanese comments and documentation
- CamelCase naming (Go standard)
- Interfaces start with uppercase (e.g., `Formatter`, `Parser`)
- Structs start with uppercase (e.g., `DotFormatter`, `MermaidFormatter`)
- Use `log/slog` for structured logging with debug, info, warn, error levels
- Resource management with proper cleanup (especially for go-graphviz)

## Implementation Patterns

### Adding New Output Formatters
1. Create new file in `pkg/std-impl/output/`
2. Implement `core.Formatter` interface
3. Add `Format(generator *core.Generator) (string, error)` method
4. Create node ID generation function (`getXxxNodeID`)
5. Create node label generation function (`getXxxNodeLabel`)
6. Handle special character conversion (comma, space, hyphen, dot → underscore)
7. Add to formatter selection in `cmd/cli/blindspot/main.go`

### Adding New Parsers
1. Create new package in `pkg/std-impl/`
2. Implement `core.Parser` interface
3. Add to parser selection logic in CLI

### Error Handling
- Use panic for type mismatches (see EdgeRule comments)
- Return proper errors for normal cases
- Use defer for resource cleanup
- Log debug information extensively

## Testing
- Test files use `*_test.go` format
- Follow existing test patterns in `stringlist_test.go`
- Include expected output validation
- Test both success and error cases

## Dependencies
- External: `github.com/goccy/go-graphviz` (DOT format output)
- Standard library preferred for other functionality
- Go 1.24.3 required

## Safety Features
### Iteration Limit Mode
⚠️ **Critical**: Running without `--limit` is dangerous and may cause infinite loops that can freeze or crash the system.

- Use `--limit` flag to prevent infinite loops in state generation
- Generator accepts optional limit parameter in NewGenerator constructor
- When limit is exceeded, generation stops with error log
- Interactive confirmation prompts for safety awareness
- **Always specify --limit in production environments**