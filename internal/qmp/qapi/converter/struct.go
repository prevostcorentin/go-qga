package converter

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
)

type Struct struct {
	name   string
	fields []*Field
}

func NewStruct(rawStruct *collector.Struct) *Struct {
	structFields, fieldCount := make([]*Field, len(rawStruct.Data)), 0
	for name, fieldType := range rawStruct.Data {
		structFields[fieldCount] = NewField(name, fieldType.(string))
		fieldCount = fieldCount + 1
	}
	return &Struct{name: rawStruct.Name(), fields: structFields}
}

func (st *Struct) Name() string {
	return st.name
}

func (st *Struct) Fields() []*Field {
	return st.fields
}

func (st *Struct) Generate() ([]byte, error) {
	fset := token.NewFileSet()

	structFields := make([]*ast.Field, len(st.fields))
	for i, field := range st.fields {
		structFields[i] = &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(field.Name())},
			Type:  ast.NewIdent(field.Type()),
		}
	}

	genDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(st.Name()),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: structFields,
					},
				},
			},
		},
	}

	file := &ast.File{
		Name:  ast.NewIdent("generated"),
		Decls: []ast.Decl{genDecl},
	}

	var buffer bytes.Buffer
	if err := format.Node(&buffer, fset, file); err != nil {
		return nil, NewCodeGenerationError(err, AbstractSyntaxTreeErrorKind)
	}
	return buffer.Bytes(), nil
}
