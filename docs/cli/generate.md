# generate

Generate artifacts from specifications.

## Synopsis

```bash
visionspec generate <subcommand> [flags]
```

## Subcommands

### tests

Generate test stubs from Test Plan Document (TPD).

```bash
visionspec generate tests [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `-p, --project` | Project name |
| `--lang` | Target language: go, ts, py (default: go) |
| `--framework` | Test framework: testing, testify, jest, vitest, pytest |
| `--output` | Output directory (default: ./generated-tests) |
| `--group-by` | Grouping strategy: type, file, priority (default: type) |

## Examples

### Generate Go Tests

```bash
# Generate Go test stubs using standard testing package
visionspec generate tests -p myproject --lang go

# Generate with testify framework
visionspec generate tests -p myproject --lang go --framework testify

# Specify output directory
visionspec generate tests -p myproject --lang go --output ./tests/generated
```

### Generate TypeScript Tests

```bash
# Generate Jest test stubs
visionspec generate tests -p myproject --lang ts --framework jest

# Generate Vitest test stubs
visionspec generate tests -p myproject --lang ts --framework vitest
```

### Generate Python Tests

```bash
# Generate pytest test stubs
visionspec generate tests -p myproject --lang py
```

## TPD Sections Parsed

The generator extracts test cases from these TPD sections:

| Section | Test Type |
|---------|-----------|
| Section 3: Functional Tests | Unit and functional tests |
| Section 4.1: API Tests | API endpoint tests |
| Section 5.1: Journey Tests | E2E and integration tests |

## Output Structure

### Go Output

```
generated-tests/
├── functional_test.go     # Functional test stubs
├── api_test.go            # API test stubs
└── journey_test.go        # Journey/E2E test stubs
```

### TypeScript Output

```
generated-tests/
├── functional.test.ts     # Functional test stubs
├── api.test.ts            # API test stubs
└── journey.test.ts        # Journey/E2E test stubs
```

### Python Output

```
generated-tests/
├── conftest.py            # pytest fixtures
├── test_functional.py     # Functional test stubs
├── test_api.py            # API test stubs
└── test_journey.py        # Journey/E2E test stubs
```

## Traceability

Generated tests include source references for traceability:

```go
func TestTC001_UserLogin(t *testing.T) {
    // Source: REQ-AUTH-001
    // TPD: Section 3, Row 1
    t.Skip("TODO: Implement test")
}
```

## See Also

- [synthesize](synthesize.md) - Generate TPD from PRD/TRD
- [eval](eval.md) - Evaluate specifications
