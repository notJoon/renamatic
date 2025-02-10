package internal

import "go/ast"

// Mapping represents a generic mapping from one string to another.
type Mapping[K comparable, V any] map[K]V

// FileProcessor defines the interface for processing source files.
type FileProcessor interface {
	ProcessFile(path string) error
	ShouldProcess(path string) bool
}

// ASTTransformer defines the interface for AST transformations.
type ASTTransformer interface {
	Transform(node ast.Node) bool
}
