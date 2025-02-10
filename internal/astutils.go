package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// CompareAble is an interface for types that can compare two values of type T.
type CompareAble[T any] interface {
	// Compare compares two values of type T and returns an error if they are not considered equal.
	Compare(got, want T) error
}

// TypeComparator encapsulates the comparison logic for a specific AST node type.
// It implements the CompareAble interface for ast.Node values.
type TypeComparator[T ast.Node] struct {
	// compareFunc performs the actual comparison between two nodes of type T.
	compareFunc func(got, want T) error
}

// Compare asserts that both got and want are of type T and then invokes compareFunc.
// It returns an error if the types do not match or if the compareFunc detects a difference.
func (tc TypeComparator[T]) Compare(got, want ast.Node) error {
	gotT, ok := got.(T)
	if !ok {
		return fmt.Errorf("type mismatch: got %T, want %T", got, want)
	}
	wantT, ok := want.(T)
	if !ok {
		return fmt.Errorf("type mismatch: got %T, want %T", got, want)
	}
	return tc.compareFunc(gotT, wantT)
}

// compareAST parses two Go source code strings and compares their ASTs.
// It returns an error if any discrepancies are found.
func compareAST(t *testing.T, got, want string) error {
	t.Helper()

	fset := token.NewFileSet()

	gotAst, err := parser.ParseFile(fset, "", got, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing got code: %v", err)
	}

	wantAst, err := parser.ParseFile(fset, "", want, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing want code: %v", err)
	}

	// Use a TypeComparator for *ast.File nodes.
	fileComparator := TypeComparator[*ast.File]{
		compareFunc: func(got, want *ast.File) error {
			if got.Name.Name != want.Name.Name {
				return fmt.Errorf("package name mismatch: got %s, want %s", got.Name.Name, want.Name.Name)
			}
			return compareDecls(got.Decls, want.Decls)
		},
	}

	if err := fileComparator.Compare(gotAst, wantAst); err != nil {
		return fmt.Errorf("AST comparison failed: %v", err)
	}

	return nil
}

// compareNode compares two ast.Node values by delegating to a type-specific comparator
// based on the dynamic type of the node.
func compareNode(got, want ast.Node) error {
	if got == nil || want == nil {
		if got != want {
			return fmt.Errorf("one node is nil, other is not")
		}
		return nil
	}

	switch got.(type) {
	case *ast.File:
		fileComparator := TypeComparator[*ast.File]{
			compareFunc: func(got, want *ast.File) error {
				if got.Name.Name != want.Name.Name {
					return fmt.Errorf("package name mismatch: got %s, want %s", got.Name.Name, want.Name.Name)
				}
				return compareDecls(got.Decls, want.Decls)
			},
		}
		return fileComparator.Compare(got, want)

	case *ast.Ident:
		identComparator := TypeComparator[*ast.Ident]{
			compareFunc: func(got, want *ast.Ident) error {
				if got.Name != want.Name {
					return fmt.Errorf("identifier mismatch: got %s, want %s", got.Name, want.Name)
				}
				return nil
			},
		}
		return identComparator.Compare(got, want)

	case *ast.CallExpr:
		callExprComparator := TypeComparator[*ast.CallExpr]{
			compareFunc: func(got, want *ast.CallExpr) error {
				if err := compareNode(got.Fun, want.Fun); err != nil {
					return fmt.Errorf("function comparison failed: %v", err)
				}
				if err := compareExprs(got.Args, want.Args); err != nil {
					return fmt.Errorf("arguments comparison failed: %v", err)
				}
				return nil
			},
		}
		return callExprComparator.Compare(got, want)

	case *ast.SelectorExpr:
		selectorComparator := TypeComparator[*ast.SelectorExpr]{
			compareFunc: func(got, want *ast.SelectorExpr) error {
				if err := compareNode(got.X, want.X); err != nil {
					return fmt.Errorf("selector X comparison failed: %v", err)
				}
				if err := compareNode(got.Sel, want.Sel); err != nil {
					return fmt.Errorf("selector Sel comparison failed: %v", err)
				}
				return nil
			},
		}
		return selectorComparator.Compare(got, want)

	default:
		return nil
	}
}

// compareSlice is a generic helper function that compares two slices element-by-element
// using the provided comparator function.
func compareSlice[T any](got, want []T, cmp func(got, want T) error) error {
	if len(got) != len(want) {
		return fmt.Errorf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range got {
		if err := cmp(got[i], want[i]); err != nil {
			return fmt.Errorf("comparison failed at index %d: %v", i, err)
		}
	}
	return nil
}

// compareDecls compares two slices of ast.Decl.
func compareDecls(got, want []ast.Decl) error {
	return compareSlice(got, want, func(gotDecl, wantDecl ast.Decl) error {
		switch gotD := gotDecl.(type) {
		case *ast.FuncDecl:
			wantD, ok := wantDecl.(*ast.FuncDecl)
			if !ok {
				return fmt.Errorf("declaration type mismatch: expected *ast.FuncDecl but got %T", wantDecl)
			}
			if gotD.Name.Name != wantD.Name.Name {
				return fmt.Errorf("function name mismatch: got %s, want %s", gotD.Name.Name, wantD.Name.Name)
			}
			return compareNode(gotD.Body, wantD.Body)
		case *ast.GenDecl:
			wantD, ok := wantDecl.(*ast.GenDecl)
			if !ok {
				return fmt.Errorf("declaration type mismatch: expected *ast.GenDecl but got %T", wantDecl)
			}
			if gotD.Tok != wantD.Tok {
				return fmt.Errorf("token mismatch: got %v, want %v", gotD.Tok, wantD.Tok)
			}
			return compareSpecs(gotD.Specs, wantD.Specs)
		default:
			return nil
		}
	})
}

// compareExprs compares two slices of ast.Expr.
func compareExprs(got, want []ast.Expr) error {
	return compareSlice(got, want, func(gotExpr, wantExpr ast.Expr) error {
		return compareNode(gotExpr, wantExpr)
	})
}

// compareSpecs compares two slices of ast.Spec.
func compareSpecs(got, want []ast.Spec) error {
	return compareSlice(got, want, func(gotSpec, wantSpec ast.Spec) error {
		switch gotS := gotSpec.(type) {
		case *ast.ImportSpec:
			wantS, ok := wantSpec.(*ast.ImportSpec)
			if !ok {
				return fmt.Errorf("spec type mismatch: got %T, want %T", gotSpec, wantSpec)
			}
			if gotS.Path.Value != wantS.Path.Value {
				return fmt.Errorf("import path mismatch: got %s, want %s", gotS.Path.Value, wantS.Path.Value)
			}
			return nil
		default:
			return nil
		}
	})
}
