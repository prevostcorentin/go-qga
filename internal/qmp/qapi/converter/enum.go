package converter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
)

type Enum struct {
	name   string
	values []string
}

func NewEnum(rawEnum *collector.Enum) *Enum {
	return &Enum{name: rawEnum.Name(), values: rawEnum.Data}
}

func (en *Enum) Name() string {
	return en.name
}

func (en *Enum) Values() []string {
	return en.values
}

func (en *Enum) Generate() ([]byte, error) {
	fset := token.NewFileSet()

	genDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(en.name),
				Type: ast.NewIdent("string"),
			},
		},
	}

	specs := make([]ast.Spec, len(en.values))

	for i, valueName := range en.values {
		specs[i] = &ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent(valueName)},
			Type:  ast.NewIdent(en.name),
			Values: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", valueName),
				},
			},
		}
	}

	constDecl := &ast.GenDecl{
		Tok:   token.CONST,
		Specs: specs,
	}

	file := &ast.File{
		Name:  ast.NewIdent("generated"),
		Decls: []ast.Decl{genDecl, constDecl},
	}

	var buffer bytes.Buffer
	if err := format.Node(&buffer, fset, file); err != nil {
		return nil, NewCodeGenerationError(err, AbstractSyntaxTreeErrorKind)
	}
	return buffer.Bytes(), nil
}
