package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	_ FileProcessor  = (*GnoFileProcessor)(nil)
	_ ASTTransformer = (*StdFunctionTransformer)(nil)
)

// GnoFileProcessor implements FileProcessor for .gno files.
type GnoFileProcessor struct {
	mapping     Mapping[string, string]
	transformer ASTTransformer
	parseMode   parser.Mode
}

// newGnoFileProcessor creates a new GnoFileProcessor with the given mapping and transformer.
func newGnoFileProcessor(mapping Mapping[string, string], transformer ASTTransformer) *GnoFileProcessor {
	return &GnoFileProcessor{
		mapping:     mapping,
		transformer: transformer,
		parseMode:   parser.ParseComments,
	}
}

// ShouldProcess checks if the given file should be processed.
func (p *GnoFileProcessor) ShouldProcess(path string) bool {
	return strings.HasSuffix(filepath.Ext(path), ".gno")
}

// ProcessFile processes a single .gno file.
func (p *GnoFileProcessor) ProcessFile(path string) error {
	fset := token.NewFileSet()
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	file, err := parser.ParseFile(fset, path, src, p.parseMode)
	if err != nil {
		return err
	}

	// Apply the transformation
	ast.Inspect(file, p.transformer.Transform)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), info.Mode())
}

// StdFunctionTransformer implements ASTTransformer for std package function calls.
type StdFunctionTransformer struct {
	mapping Mapping[string, string]
}

// newStdFunctionTransformer creates a new StdFunctionTransformer.
func newStdFunctionTransformer(mapping Mapping[string, string]) *StdFunctionTransformer {
	return &StdFunctionTransformer{mapping: mapping}
}

// Transform implements the AST transformation for std package function calls.
func (t *StdFunctionTransformer) Transform(n ast.Node) bool {
	if selExpr, ok := n.(*ast.SelectorExpr); ok {
		if isStdChain(selExpr.X) {
			if newName, exists := t.mapping[selExpr.Sel.Name]; exists {
				selExpr.Sel.Name = newName
			}
		}
	}
	return true
}

// DirectoryProcessor handles recursive directory processing.
type DirectoryProcessor struct {
	processor FileProcessor
}

// newDirectoryProcessor creates a new DirectoryProcessor.
func newDirectoryProcessor(processor FileProcessor) *DirectoryProcessor {
	return &DirectoryProcessor{processor: processor}
}

// ProcessDir recursively processes all matching files in the specified directory.
func (dp *DirectoryProcessor) ProcessDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && dp.processor.ShouldProcess(path) {
			if err := dp.processor.ProcessFile(path); err != nil {
				return fmt.Errorf("failed to process file %s: %w", path, err)
			}
		}

		return nil
	})
}

// isStdChain recursively checks whether the given expression is rooted in the "std" identifier.
func isStdChain(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name == "std"
	case *ast.SelectorExpr:
		return isStdChain(e.X)
	case *ast.CallExpr:
		return isStdChain(e.Fun)
	default:
		return false
	}
}
