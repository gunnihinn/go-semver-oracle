package types

import (
	"go/ast"
)

// Id identifies a definition in a Go package.
// The only things that matter about this type are:
//
// - Each declaration in a Go package gets a unique one.
// - The type can be used as a map key.
type Id struct {
	name     string
	receiver string // name of receiver type
}

// Identifiable is an umbrella for different definitions.
type Identifiable interface {
	Identify() Id
}

// Result is the result of comparing identifiable types.
type Result int

const (
	Equal Result = iota // Equality.
	Minor               // Minor semver difference.
	Major               // Major semver difference.
)

// Cmp computes the semver difference from a to b.
func Cmp(a, b Identifiable) Result {
	// TODO: Use reflect to compare types
}

// Identifiable types, one per Go declaration type

// VarDef is the definition of a single variable.
type VarDef struct {
	decl *ast.GenDecl
	spec *ast.ValueSpec
}

func (def VarDef) Identify() Id {
	return Id{
		name:     def.spec.Names[0].Name,
		receiver: nil,
	}
}

// ConstDef is the definition of a single constant.
type ConstDef struct {
	decl *ast.GenDecl
	spec *ast.ValueSpec
}

func (def ConstDef) Identify() Id {
	return Id{
		name:     def.spec.Names[0].Name,
		receiver: nil,
	}
}

// TypeDef is the definition of a single type.
type TypeDef struct {
	decl *ast.GenDecl
	spec *ast.TypeSpec
}

func (def TypeDef) Identify() Id {
	return Id{
		name:     def.spec.Name,
		receiver: nil,
	}
}

// FuncDef is the definition of a single function.
type FuncDef struct {
	decl *ast.FuncDecl
}

func (def FuncDef) Identify() Id {
	var recv string
	if len(def.Recv.List) > 0 {
		recv = def.Recv.List[0].Type // TODO
	} else {
		recv = nil
	}

	return Id{
		name:     def.decl.Name,
		receiver: recv,
	}
}

// NoDef is a missing definition.
type NoDef struct {
	name string
}

func (def NoDef) Identify() Id {
	return Id{
		name: name,
	}
}
