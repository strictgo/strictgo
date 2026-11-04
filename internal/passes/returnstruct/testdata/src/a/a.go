package a

import (
	"contracts"
	"go/ast"
)

type visitor struct{}

func (visitor) Visit(ast.Node) ast.Visitor { return nil }

var _ ast.Visitor = visitor{}

type pointerVisitor struct{}

func (*pointerVisitor) Visit(ast.Node) ast.Visitor { return nil }

var _ ast.Visitor = (*pointerVisitor)(nil)

type complete struct{}

func (complete) Get() contracts.Value { return nil }
func (*complete) Close()              {}

var _ contracts.Contract = (*complete)(nil)

// Matching one method does not satisfy the whole external contract.
type incomplete struct{}

func (incomplete) Get() contracts.Value { return nil } // want "Get returns an interface"

type unrelated struct{}

func (unrelated) Get() any              { return nil }
func (unrelated) Visit(int) ast.Visitor { return nil } // want "Visit returns an interface"

func factory() ast.Visitor  { return nil } // want "factory returns an interface"
func allowed() (any, error) { return nil, nil }

// Interfaces owned by this package do not force external API constraints.
type local interface{ Local() contracts.Value }
type localImpl struct{}

func (localImpl) Local() contracts.Value { return nil } // want "Local returns an interface"
var _ local = localImpl{}

type embedded struct{}

func (embedded) Build() contracts.Value { return nil }
func (embedded) Close()                 {}

var _ contracts.Embedded = embedded{}

// Implementing an external contract does not exempt unrelated methods.
func (complete) Extra() contracts.Value { return nil } // want "Extra returns an interface"
