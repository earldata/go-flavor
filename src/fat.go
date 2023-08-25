package main

import (
	"go/ast"
	"go/token"
)

type fatVisitor struct {
	fat int
}

func Fat(fn ast.Node) int {
	visitor := fatVisitor{
		fat: 1,
	}
	ast.Walk(&visitor, fn)
	return visitor.fat
}

func (v *fatVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
		v.fat++
	case *ast.CommClause:
		if n.Comm != nil {
			v.fat++
		}
	case *ast.CaseClause:
		if n.List != nil {
			v.fat++
		}
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			v.fat++
		}
	}
	return v
}
