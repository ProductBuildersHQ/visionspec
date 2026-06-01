package testgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParserParseFunctionalTests(t *testing.T) {
	content := `# Test Plan

## 3. Test Cases from PRD

### 3.1 Feature: User Login

| ID | Test Case | Input | Expected Output | Priority |
|----|-----------|-------|-----------------|----------|
| TC-001 | Valid login | username: test, password: pass | Login success | P0 |
| TC-002 | Invalid password | username: test, password: wrong | Error message | P1 |
`

	parser := NewParser()
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.FunctionalTests) != 2 {
		t.Errorf("Expected 2 functional tests, got %d", len(result.FunctionalTests))
	}

	if len(result.FunctionalTests) > 0 {
		tc := result.FunctionalTests[0]
		if tc.ID != "TC-001" {
			t.Errorf("Expected ID TC-001, got %s", tc.ID)
		}
		if tc.Priority != "P0" {
			t.Errorf("Expected priority P0, got %s", tc.Priority)
		}
	}
}

func TestParserParseAPITests(t *testing.T) {
	content := `# Test Plan

## 4. Technical Test Cases from TRD

### 4.1 API Testing

| Endpoint | Method | Test Scenario | Expected Response | Priority |
|----------|--------|---------------|-------------------|----------|
| /users | GET | Happy path | 200 OK with user list | P0 |
| /users | POST | Invalid input | 400 Bad Request | P1 |
`

	parser := NewParser()
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.APITests) != 2 {
		t.Errorf("Expected 2 API tests, got %d", len(result.APITests))
	}

	if len(result.APITests) > 0 {
		api := result.APITests[0]
		if api.Endpoint != "/users" {
			t.Errorf("Expected endpoint /users, got %s", api.Endpoint)
		}
		if api.Method != "GET" {
			t.Errorf("Expected method GET, got %s", api.Method)
		}
	}
}

func TestParserParseJourneyTests(t *testing.T) {
	content := `# Test Plan

## 5. User Journey Testing from UXD

### 5.1 Critical User Journeys

| Journey | Steps | Assertions | Priority |
|---------|-------|------------|----------|
| User Signup | 1. Visit signup; 2. Fill form; 3. Submit | Account created; Email sent | P0 |
`

	parser := NewParser()
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.JourneyTests) != 1 {
		t.Errorf("Expected 1 journey test, got %d", len(result.JourneyTests))
	}

	if len(result.JourneyTests) > 0 {
		journey := result.JourneyTests[0]
		if journey.Journey != "User Signup" {
			t.Errorf("Expected journey 'User Signup', got %s", journey.Journey)
		}
		if len(journey.Steps) < 3 {
			t.Errorf("Expected at least 3 steps, got %d", len(journey.Steps))
		}
	}
}

func TestGoGenerator(t *testing.T) {
	cases := []TestCase{
		{
			ID:       "TC-001",
			Title:    "Valid Login",
			Type:     TestTypeFunctional,
			Input:    "username: test, password: pass",
			Expected: "Login success",
			Priority: "P0",
		},
		{
			ID:       "API-001",
			Title:    "GET /users returns list",
			Type:     TestTypeAPI,
			Input:    "GET /users",
			Expected: "200 OK",
			Priority: "P0",
		},
	}

	tmpDir, err := os.MkdirTemp("", "testgen-go-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := &GoGenerator{}
	opts := GenerateOptions{
		OutputDir:   tmpDir,
		PackageName: "mytest",
		GroupBy:     "type",
	}

	result, err := gen.Generate(cases, opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result.TotalTests != 2 {
		t.Errorf("Expected 2 total tests, got %d", result.TotalTests)
	}

	// Check files exist
	for _, f := range result.Files {
		if _, err := os.Stat(f.Path); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", f.Path)
		}

		content, err := os.ReadFile(f.Path)
		if err != nil {
			t.Errorf("Failed to read %s: %v", f.Path, err)
			continue
		}

		// Check package declaration
		if !strings.Contains(string(content), "package mytest") {
			t.Errorf("Expected package declaration in %s", f.Path)
		}

		// Check it has test functions
		if !strings.Contains(string(content), "func Test") {
			t.Errorf("Expected test functions in %s", f.Path)
		}
	}
}

func TestTypeScriptGenerator(t *testing.T) {
	cases := []TestCase{
		{
			ID:       "TC-001",
			Title:    "Valid Login",
			Type:     TestTypeFunctional,
			Input:    "username: test",
			Expected: "Login success",
			Priority: "P0",
		},
	}

	tmpDir, err := os.MkdirTemp("", "testgen-ts-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := &TypeScriptGenerator{}
	opts := GenerateOptions{
		OutputDir: tmpDir,
		GroupBy:   "type",
	}

	result, err := gen.Generate(cases, opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result.TotalTests != 1 {
		t.Errorf("Expected 1 total test, got %d", result.TotalTests)
	}

	// Check file content
	for _, f := range result.Files {
		content, err := os.ReadFile(f.Path)
		if err != nil {
			t.Errorf("Failed to read %s: %v", f.Path, err)
			continue
		}

		if !strings.Contains(string(content), "describe(") {
			t.Errorf("Expected describe block in %s", f.Path)
		}
		if !strings.Contains(string(content), "test(") {
			t.Errorf("Expected test block in %s", f.Path)
		}
	}
}

func TestPythonGenerator(t *testing.T) {
	cases := []TestCase{
		{
			ID:       "TC-001",
			Title:    "Valid Login",
			Type:     TestTypeFunctional,
			Input:    "username: test",
			Expected: "Login success",
			Priority: "P0",
		},
	}

	tmpDir, err := os.MkdirTemp("", "testgen-py-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := &PythonGenerator{}
	opts := GenerateOptions{
		OutputDir: tmpDir,
		GroupBy:   "type",
	}

	result, err := gen.Generate(cases, opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result.TotalTests != 1 {
		t.Errorf("Expected 1 total test, got %d", result.TotalTests)
	}

	// Check conftest.py was generated
	conftestPath := filepath.Join(tmpDir, "conftest.py")
	if _, err := os.Stat(conftestPath); os.IsNotExist(err) {
		t.Error("Expected conftest.py to be generated")
	}

	// Check test file content
	for _, f := range result.Files {
		if strings.HasSuffix(f.Path, "conftest.py") {
			continue
		}

		content, err := os.ReadFile(f.Path)
		if err != nil {
			t.Errorf("Failed to read %s: %v", f.Path, err)
			continue
		}

		if !strings.Contains(string(content), "class Test") {
			t.Errorf("Expected test class in %s", f.Path)
		}
		if !strings.Contains(string(content), "def test_") {
			t.Errorf("Expected test method in %s", f.Path)
		}
	}
}

func TestParsedTPDAllTestCases(t *testing.T) {
	tpd := &ParsedTPD{
		FunctionalTests: []TestCase{
			{ID: "TC-001", Title: "Test 1", Type: TestTypeFunctional},
		},
		APITests: []APITestCase{
			{Endpoint: "/users", Method: "GET", Scenario: "Happy path"},
		},
		JourneyTests: []JourneyTestCase{
			{Journey: "User Signup", Steps: []string{"Step 1", "Step 2"}},
		},
	}

	all := tpd.AllTestCases()
	if len(all) != 3 {
		t.Errorf("Expected 3 total test cases, got %d", len(all))
	}
}

func TestRegistry(t *testing.T) {
	// Check all generators are registered
	available := Available()
	if len(available) < 3 {
		t.Errorf("Expected at least 3 generators, got %d", len(available))
	}

	// Check we can get each generator
	for _, lang := range []string{"go", "ts", "py"} {
		gen, err := Get(lang)
		if err != nil {
			t.Errorf("Failed to get %s generator: %v", lang, err)
			continue
		}
		if gen.Language() != lang {
			t.Errorf("Expected language %s, got %s", lang, gen.Language())
		}
	}

	// Check unknown generator returns error
	_, err := Get("unknown")
	if err == nil {
		t.Error("Expected error for unknown generator")
	}
}
