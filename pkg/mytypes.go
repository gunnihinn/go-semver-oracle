package mytypes

import "go/ast"

// Go types

type Comparable interface {
	compare()
}

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
	orig   *ast.Node
}

func (a Array) compare() {}

// Ellipsis is a Go variadic function argument.
//
// We need to special-case it because it permits a chain of minor changes that
// we would otherwise signal as major ones:
//
//	func(a int) -> func(as ...int) -> func(as []int)
//       |      min                min     |
//       \-------------- maj --------------/
//
// Note that compositions of minor changes are not minor, which annoys my inner
// algebraist.
type Ellipsis struct {
	Values Comparable
	orig   *ast.Node
}

func (e Ellipsis) compare() {}

// Map is a Go map.
type Map struct {
	Key   Comparable
	Value Comparable
	orig  *ast.Node
}

func (m Map) compare() {}

// Struct is a Go struct literal.
type Struct struct {
	Fields map[string]Comparable
	orig   *ast.Node
}

func (s Struct) compare() {}

// Lambda is an anonymous function.
type Lambda struct {
	Inputs  []Comparable
	Outputs []Comparable
	orig    *ast.Node
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
	orig *ast.Node
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
	orig *ast.Node
}

func (v Var) compare() {}
func (v Var) declare() {}

// Const is a const declaration.
type Const struct {
	Name string
	Type Comparable
	orig *ast.Node
}

func (c Const) compare() {}
func (c Const) declare() {}

// Type is a type declaration.
type Type struct {
	Name string
	Type Comparable
	orig *ast.Node
}

func (t Type) compare() {}
func (t Type) declare() {}

// Func is a func declaration.
type Func struct {
	Name     string
	Receiver Comparable
	Lambda
	orig *ast.Node
}

func (f Func) compare() {}
func (f Func) declare() {}
