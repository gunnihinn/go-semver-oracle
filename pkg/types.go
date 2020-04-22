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

func max(as ...Result) Result {
	if len(as) == 0 {
		panic("Empty input")
	}

	m := as[0]
	for _, a := range as[1:] {
		if a > m {
			m = a
		}
	}

	return m
}

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
	return cmpExpr(a.spec.Type, b.spec.Type)
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

func cmpExpr(a, b ast.Expr) Result {
	/*
		*ast.Expr nodes can be:

		*BadExpr
		*Ident
		*Ellipsis
		*BasicLit
		*FuncLit
		*CompositeLit
		*ParenExpr
		*SelectorExpr
		*IndexExpr
		*SliceExpr
		*TypeAssertExpr
		*CallExpr
		*StarExpr
		*UnaryExpr
		*BinaryExpr
		*KeyValueExpr

		*ArrayType
		*StructType
		*FuncType
		*InterfaceType
		*MapType
		*ChanType

		Not all of those can appear in variable definitions.
	*/

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return Major
	}

	switch ta := a.(type) {
	case *ast.BasicLit:
		tb := b.(*ast.BasicLit)
		if ta.Kind == tb.Kind {
			// TODO: Compare ta.Value and tb.Value
			return Equal
		} else {
			return Major
		}

	case *ast.FuncLit:

	case *ast.ArrayType:
		tb := b.(*ast.ArrayType)
		return cmpExpr(ta.Elt, tb.Elt)

	case *ast.StructType:

	case *ast.FuncType:

	case *ast.InterfaceType:

	case *ast.MapType:
		tb := b.(*ast.MapType)
		return max(
			cmpExpr(ta.Key, tb.Key),
			cmpExpr(ta.Value, tb.Value),
		)

	case *ast.ChanType:
		tb := b.(*ast.ChanType)
		if ta.Dir > tb.Dir {
			return Major
		} else {
			vals := cmpExpr(ta.Value, tb.Value)

			var dirs Result
			if ta.Dir == tb.Dir {
				dirs = Equal
			} else {
				dirs = Minor
			}

			return max(dirs, vals)
		}

	default:
		panic(fmt.Sprintf("Unexpected type %s", ta))
	}

	return Equal
}
