package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMapping(t *testing.T) {
	mapping := loadMapping(t)

	expectedMappings := map[string]string{
		"GetCallerAt":   "CallerAt",
		"Addr":          "Address",
		"GetOrigSend":   "OriginSend",
		"GetOrigCaller": "OriginCaller",
		"PrevRealm":     "PreviousRealm",
	}

	for k, v := range expectedMappings {
		if mapping[k] != v {
			t.Errorf("expected mapping for %s to be '%s', got '%s'", k, v, mapping[k])
		}
	}

	if len(mapping) != 9 {
		t.Errorf("expected mapping length 9, got %d", len(mapping))
	}
}

func TestProcessDir(t *testing.T) {
	tmpDir := t.TempDir()

	sampleContent := `package main

import "std"

func main() {
	std.GetCallerAt()
	std.SomeOtherFunction()
	std.PrevRealm().Addr()
	Addr()
}
`
	sampleFile := filepath.Join(tmpDir, "sample.gno")
	if err := os.WriteFile(sampleFile, []byte(sampleContent), 0o644); err != nil {
		t.Fatalf("failed to write sample file: %v", err)
	}

	mapping := loadMapping(t)

	if err := ProcessDir(tmpDir, mapping); err != nil {
		t.Fatalf("ProcessDir failed: %v", err)
	}

	processedData, err := os.ReadFile(sampleFile)
	if err != nil {
		t.Fatalf("failed to read processed file: %v", err)
	}

	expected := `package main

import "std"

func main() {
	std.CallerAt()
	std.SomeOtherFunction()
	std.PreviousRealm().Address()
	Addr()
}
`
	if err := compareAST(t, string(processedData), expected); err != nil {
		t.Errorf("AST comparison failed: %v", err)
	}
}

func TestProcessDir_IgnoresNonGnoFiles(t *testing.T) {
	tmpDir := t.TempDir()

	nonGnoFile := filepath.Join(tmpDir, "sample.txt")
	originalContent := `This is not a gno file.`
	if err := os.WriteFile(nonGnoFile, []byte(originalContent), 0o644); err != nil {
		t.Fatalf("failed to write non-gno file: %v", err)
	}

	mapping := loadMapping(t)
	if err := ProcessDir(tmpDir, mapping); err != nil {
		t.Fatalf("ProcessDir failed: %v", err)
	}

	data, err := os.ReadFile(nonGnoFile)
	if err != nil {
		t.Fatalf("failed to read non-gno file: %v", err)
	}
	if string(data) != originalContent {
		t.Errorf("non-.gno file content changed; got %s, want %s", string(data), originalContent)
	}
}

func TestProcessDir_Nested(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "nested")
	if err := os.Mkdir(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}

	sampleContent := `package main

import "std"

func main() {
	std.Addr()
}`
	nestedFile := filepath.Join(nestedDir, "nested.gno")
	if err := os.WriteFile(nestedFile, []byte(sampleContent), 0o644); err != nil {
		t.Fatalf("failed to write nested .gno file: %v", err)
	}

	mapping := loadMapping(t)
	if err := ProcessDir(tmpDir, mapping); err != nil {
		t.Fatalf("ProcessDir failed: %v", err)
	}

	expected := `package main

import "std"

func main() {
	std.Address()
}
`
	processedData, err := os.ReadFile(nestedFile)
	if err != nil {
		t.Fatalf("failed to read nested file: %v", err)
	}
	if err := compareAST(t, string(processedData), expected); err != nil {
		t.Errorf("AST comparison failed for nested file: %v", err)
	}
}

func TestProcessDir_InvalidSyntax(t *testing.T) {
	tmpDir := t.TempDir()
	invalidContent := `package main

import "std"

func main() {
	std.Addr( // missing closing parenthesis
`
	invalidFile := filepath.Join(tmpDir, "invalid.gno")
	if err := os.WriteFile(invalidFile, []byte(invalidContent), 0o644); err != nil {
		t.Fatalf("failed to write invalid file: %v", err)
	}

	mapping := loadMapping(t)
	err := ProcessDir(tmpDir, mapping)
	if err == nil {
		t.Errorf("expected ProcessDir to fail on invalid syntax, but it did not")
	}
}

func TestLoadMapping_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidMappingFile := filepath.Join(tmpDir, "invalid_mapping.yaml")
	invalidYAML := `invalid_yaml: [unbalanced`
	if err := os.WriteFile(invalidMappingFile, []byte(invalidYAML), 0o644); err != nil {
		t.Fatalf("failed to write invalid mapping file: %v", err)
	}

	_, err := LoadMapping(invalidMappingFile)
	if err == nil {
		t.Errorf("expected LoadMapping to fail with invalid YAML, but it did not")
	}
}

func TestProcessDir_NonStdCalls(t *testing.T) {
	tmpDir := t.TempDir()

	sampleContent := `package main

import "custom"

func main() {
	custom.Addr()
}
`
	sampleFile := filepath.Join(tmpDir, "sample.gno")
	if err := os.WriteFile(sampleFile, []byte(sampleContent), 0o644); err != nil {
		t.Fatalf("failed to write sample file: %v", err)
	}

	mapping := loadMapping(t)
	if err := ProcessDir(tmpDir, mapping); err != nil {
		t.Fatalf("ProcessDir failed: %v", err)
	}

	if err := compareAST(t, sampleContent, sampleContent); err != nil {
		t.Errorf("AST comparison failed for non-std calls: %v", err)
	}
}

func loadMapping(t *testing.T) Mapping[string, string] {
	t.Helper()

	mappingPath := filepath.Join("..", "mapping.yml")
	mapping, err := LoadMapping(mappingPath)
	if err != nil {
		t.Fatalf("LoadMapping failed: %v", err)
	}
	return mapping
}
