package mytypes

import (
	"fmt"
	"go/ast"
)

// Go types

type Comparable interface {
	compare()
}

// Null is nothing.
type Null struct{}

func (n Null) compare() {}

// Ident is a Go identifier.
type Ident struct {
	Name string
	orig ast.Node
}

func (n Ident) compare() {}

// Primitive is a Go primitive.
type Primitive int

const (
	pBool Primitive = iota

	pString

	pInt
	pInt8
	pInt16
	pInt32
	pInt64

	pUint
	pUint8
	pUint16
	pUint32
	pUint64
	pUintptr

	pByte

	pRune

	pFloat32
	pFloat64

	pComplex64
	pComplex128

	pError // convenient lie
)

func (p Primitive) compare() {}

// Array is a Go array.
type Array struct {
	Values Comparable
	orig   ast.Node
}

func (a Array) compare() {}

// Ellipsis is a Go variadic function argument.
type Ellipsis struct {
	Values Comparable
	orig   ast.Node
}

func (e Ellipsis) compare() {}

// TODO: Do we need a Placeholder? It'd be a struct with an ast.Node and a
// string type name. We'd use it when parsing functions whose arguments are
// user-defined types we haven't seen yet.

// Map is a Go map.
type Map struct {
	Key   Comparable
	Value Comparable
	orig  ast.Node
}

func (m Map) compare() {}

// Struct is a Go struct literal.
type Struct struct {
	Fields map[string]Comparable
	orig   ast.Node
}

func (s Struct) compare() {}

// Lambda is an anonymous function.
type Lambda struct {
	Inputs  []Comparable
	Outputs []Comparable
	orig    ast.Node
}

func (l Lambda) compare() {}

// Interface is a Go interface.
type Interface struct {
	Methods map[string]Lambda
}

func (i Interface) compare() {}

// Channel is a Go channel.
type Channel struct {
	Type Comparable
	Dir  Direction
	orig ast.Node
}

func (c Channel) compare() {}

type Direction int

const (
	Send        Direction = 1 << 0
	Receive               = 1 << 1
	SendReceive           = Send | Receive
)

// Declarations.

// Declaration is a Go constant, variable, type or function declaration.
type Declaration interface {
	declare()
}

// Var is a var declaration.
type Var struct {
	Name string
	Type Comparable
	orig ast.Node
}

func (v Var) compare() {}
func (v Var) declare() {}

// Const is a const declaration.
type Const struct {
	Name string
	Type Comparable
	orig ast.Node
}

func (c Const) compare() {}
func (c Const) declare() {}

// Type is a type declaration.
type Type struct {
	Name string
	Type Comparable
	orig ast.Node
}

func (t Type) compare() {}
func (t Type) declare() {}

// Func is a func declaration.
type Func struct {
	Name     string
	Receiver Comparable
	Lambda
	orig ast.Node
}

func (f Func) compare() {}
func (f Func) declare() {}

func Parse(p *ast.Package) ([]Declaration, error) {
	decls := make([]Declaration, 0)

	for _, file := range p.Files {
		for _, d := range file.Decls {
			switch decl := d.(type) {
			case *ast.GenDecl:
				for _, s := range decl.Specs {
					switch spec := s.(type) {
					case *ast.ValueSpec:
						t, err := parseExpr(spec.Type)
						if err != nil {
							return decls, err
						}

						for _, name := range spec.Names {
							decls = append(decls, Var{
								Name: name.Name,
								Type: t,
								orig: name,
							})
						}

					case *ast.TypeSpec:
						t, err := parseExpr(spec.Type)
						if err != nil {
							return decls, err
						}

						decls = append(decls, Type{
							Name: spec.Name.Name,
							Type: t,
							orig: spec,
						})

					case *ast.ImportSpec:
						return nil, fmt.Errorf("Unexpected import")

					default:
						return nil, fmt.Errorf("Unexpected GenDecl type %T", spec)
					}
				}

			case *ast.FuncDecl:
				var r Comparable
				if decl.Recv == nil {
					r = Null{}
				} else {
					var err error
					r, err = parseExpr(decl.Recv.List[0].Type)
					if err != nil {
						return decls, err
					}
				}

				lambda, err := parseLambda(decl.Type)
				if err != nil {
					return decls, err
				}

				decls = append(decls, Func{
					Name:     decl.Name.Name,
					Receiver: r,
					orig:     decl,
					Lambda:   lambda,
				})

			case *ast.BadDecl:
				// TODO
				panic("Not implemented")

			default:
				return nil, fmt.Errorf("Unexpected declaration type %T", decl)
			}
		}
	}

	return decls, nil
}

func parseExpr(e ast.Expr) (Comparable, error) {
	/*
		ast.Expr nodes can be:

		BadExpr
		Ident
		Ellipsis
		BasicLit
		FuncLit
		CompositeLit
		ParenExpr
		SelectorExpr
		IndexExpr
		SliceExpr
		TypeAssertExpr
		CallExpr
		StarExpr
		UnaryExpr
		BinaryExpr
		KeyValueExpr

		ArrayType
		StructType
		FuncType
		InterfaceType
		MapType
		ChanType

		Not all of those can appear in variable definitions.
	*/

	switch te := e.(type) {
	case *ast.Ident:
		return Ident{
			Name: te.Name,
			orig: te,
		}, nil

	case *ast.BasicLit:

	case *ast.FuncLit:

	case *ast.ArrayType:

	case *ast.StructType:

	case *ast.FuncType:

	case *ast.InterfaceType:

	case *ast.MapType:

	case *ast.ChanType:

	default:
		return Null{}, fmt.Errorf("Unexpected type %s", te)
	}

	return Null{}, nil
}

func parseLambda(e ast.Expr) (Lambda, error) {
	return Lambda{}, nil
}
