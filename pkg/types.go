package types

import (
	"fmt"
	"go/ast"
	"reflect"
)

// ID identifies a definition in a Go package.
// The only things that matter about this type are:
//
// - Each declaration in a Go package gets a unique one.
// - The type can be used as a map key.
type ID struct {
	name     string
	receiver string // name of receiver type
}

func eqID(a, b ID) bool {
	return a.name == b.name && a.receiver == b.receiver
}

// Identifiable is an umbrella for different definitions.
type Identifiable interface {
	Identify() ID
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
	if !eqID(a.Identify(), b.Identify()) {
		// TODO: Don't panic?
		panic("Can't Cmp different a and b")
	}

	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	// From nothing...
	if ta == reflect.TypeOf(NoDef{}) {
		if tb == reflect.TypeOf(NoDef{}) {
			// ... to nothing
			// How would we ever end up here?
			return Equal
		} else {
			// ... to something
			return Minor
		}
	}

	if ta != tb {
		// Type changes are major because Go has no type hierarchy or inference
		// so, for example, int -> int64 is a major change even though int fits
		// in int64. In a better type system this might only be a minor change,
		// but here Go's "volgt rusl" system makes things easy for us.
		return Major
	}

	// Same type, not nothing. Have to go type by type.
	switch ta {
	case reflect.TypeOf(VarDef{}):
		x := a.(VarDef)
		y := b.(VarDef)
		return cmpVarDef(x, y)

	case reflect.TypeOf(ConstDef{}):
		x := a.(ConstDef)
		y := b.(ConstDef)
		return cmpConstDef(x, y)

	case reflect.TypeOf(TypeDef{}):
		x := a.(TypeDef)
		y := b.(TypeDef)
		return cmpTypeDef(x, y)

	case reflect.TypeOf(FuncDef{}):
		x := a.(FuncDef)
		y := b.(FuncDef)
		return cmpFuncDef(x, y)

	default:
		// We should have dealt with NoDef already and there aren't other types
		// that satisfy Identifiable, but just in case.
		panic(fmt.Sprintf("Unexpected type %s", ta))
	}
}

// Identifiable types, one per Go declaration type

// VarDef is the definition of a single variable.
type VarDef struct {
	decl *ast.GenDecl
	spec *ast.ValueSpec
}

func (def VarDef) Identify() ID {
	return ID{
		name: def.spec.Names[0].Name,
	}
}

func cmpVarDef(a, b VarDef) Result {
	// TODO
	return Equal
}

// ConstDef is the definition of a single constant.
type ConstDef struct {
	decl *ast.GenDecl
	spec *ast.ValueSpec
}

func (def ConstDef) Identify() ID {
	return ID{
		name: def.spec.Names[0].Name,
	}
}

func cmpConstDef(a, b ConstDef) Result {
	// TODO
	return Equal
}

// TypeDef is the definition of a single type.
type TypeDef struct {
	decl *ast.GenDecl
	spec *ast.TypeSpec
}

func (def TypeDef) Identify() ID {
	return ID{
		name: def.spec.Name.Name,
	}
}

func cmpTypeDef(a, b TypeDef) Result {
	// TODO
	return Equal
}

// FuncDef is the definition of a single function.
type FuncDef struct {
	decl *ast.FuncDecl
}

func (def FuncDef) Identify() ID {
	var recv string
	if len(def.decl.Recv.List) > 0 {
		recv = def.decl.Recv.List[0].Type // TODO
	} else {
		recv = ""
	}

	return ID{
		name:     def.decl.Name.Name,
		receiver: recv,
	}
}

func cmpFuncDef(a, b FuncDef) Result {
	// TODO
	return Equal
}

// NoDef is a missing definition.
type NoDef struct {
	name string
}

func (def NoDef) Identify() ID {
	return ID{
		name: def.name,
	}
}
