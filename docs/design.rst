Design
======

The heart of this software is a diff function whose signature is something
like::

    func Diff(a *ast.Package, b *ast.Package) ([]Difference, error)

The struct in the return type contains at least::

    type Difference struct {
        Name    string
        Package string
        Left    ast.Decl
        Right   ast.Decl
    }

Telling the user what's actually different between ``Left`` and ``Right`` will
be punted to some other function. It can start by checking whether they're the
same type, and if so going case by case over the different types.

The ``Diff`` function should go through each package and build a
``map[Identifier]ast.Decl`` of names to nodes. Then we compare the maps and
construct ``Difference`` structs. We can't map bare names to declarations
because it's legal to define methods by the same name on different types in
the same package.

I don't think there's a ``Missing`` node type. Can we define our own?

If the packages ``a`` and ``b`` have different names we error out.

These are the declarations we have to handle:

Variable
    Node: ``ast.GenDecl``.
    Token: ``token.VAR``.
    Possibly parenthesized, so includes multiple variable definitions.
    Have an array of ``ast.ValueSpec`` nodes.
    See `the spec <https://golang.org/ref/spec#Variable_declarations>`_.
    
    Examples: ::

        var X int = 0
        var X, Y int = 1, 2 // vars have same type
        var (
            X int = 1
            Y float64 = 2.0
        )

Constant
    Node: ``ast.GenDecl``.
    Token: ``token.CONST``.
    The ``ast.Spec`` fields are ``ast.ValueSpec`` as for variables.
    The same comments and examples apply.
    See `the spec <https://golang.org/ref/spec#Constant_declarations>`_.

Type
    Node: ``ast.GenDecl``.
    Token: ``token.TYPE``.
    Possibly parenthesized, so includes multiple variable definitions.
    Have an array of ``ast.TypeSpec`` nodes.
    See `the spec <https://golang.org/ref/spec#Type_declarations>`_.

Function
    Node: ``ast.FuncDecl``.
    Its name is in the ``Name *ast.Ident`` field.
    If the function is a method, the ``Recv *FieldList`` is non-nil.
    If it is non-nil, the spec says it MUST contain one element.
    The ``Type FuncType`` field contains the ``Parameter`` and ``Result
    *FieldList``.
    Of those, we're only interested in each ``*Field``'s ``Type`` field.

    The ``Field`` can be an ``ast.Ident`` (a type or interface), a
    ``ast.StructType`` (literal struct), a ``ast.FuncType``, ``ast.MapType``,
    ``ast.ArrayType`` for each of those, or an ``ast.Ellipsis`` for variadic
    arguments.
