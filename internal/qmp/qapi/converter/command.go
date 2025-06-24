package converter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"unicode"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
)

type Command struct {
	name        string
	arguments   []*Field
	returnsType string
}

func NewCommand(rawCommand *collector.Command) *Command {
	commandArguments, argumentsCounter := make([]*Field, len(rawCommand.Arguments)), 0
	for argumentName, argumentType := range rawCommand.Arguments {
		commandArguments[argumentsCounter] = NewField(argumentName, argumentType)
		argumentsCounter = argumentsCounter + 1
	}
	return &Command{name: rawCommand.Name(), arguments: commandArguments, returnsType: rawCommand.Returns}
}

func (command *Command) Name() string {
	return command.name
}

func (command *Command) safeName() string {
	safeNameBuilder := strings.Builder{}
	capitalizeLetter := true // The first letter has to be capitalized
	for _, c := range command.name {
		if c == '-' {
			capitalizeLetter = true
		} else {
			letter := rune(c)
			if capitalizeLetter {
				letter = unicode.ToUpper(letter)
			}
			safeNameBuilder.WriteRune(letter)
			capitalizeLetter = false
		}
	}
	return safeNameBuilder.String()
}

func (command *Command) Generate() ([]byte, error) {
	fset := token.NewFileSet()

	commandStructDecl := command.generateMainStruct()
	commandArgumentsDecl := command.generateArgumentsStruct()
	commandResponseDecl := command.generateResponseStruct()
	constructorDecl := command.generateConstructor()

	file := &ast.File{
		Name:  ast.NewIdent("generated"),
		Decls: []ast.Decl{commandStructDecl, commandArgumentsDecl, commandResponseDecl, constructorDecl},
	}

	var buffer bytes.Buffer
	if err := format.Node(&buffer, fset, file); err != nil {
		return nil, NewCodeGenerationError(err, AbstractSyntaxTreeErrorKind)
	}
	return buffer.Bytes(), nil
}

func (command *Command) generateMainStruct() *ast.GenDecl {
	structFields := command.transformArguments(func(argument *Field) *ast.Field {
		return &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(argument.Name())},
			Type:  ast.NewIdent(argument.Type()),
		}
	})
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(command.safeName()),
				Type: &ast.StructType{
					Fields: structFields,
				},
			},
		},
	}
}
func (command *Command) generateArgumentsStruct() *ast.GenDecl {
	structCapitalizedFields := command.transformArguments(func(argument *Field) *ast.Field {
		return &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(argument.CapitalizedName())},
			Type:  ast.NewIdent(argument.Type()),
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", argument.Name()),
			},
		}
	})
	commandSafeName := command.safeName()
	argumentsStructName := fmt.Sprintf("%c%sArguments", command.name[0], commandSafeName[1:])
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(argumentsStructName),
				Type: &ast.StructType{
					Fields: structCapitalizedFields,
				},
			},
		},
	}
}

func (command *Command) generateResponseStruct() *ast.GenDecl {
	commandResponseFuncName := fmt.Sprintf("%c%sResponse", command.name[0], command.safeName()[1:])
	return &ast.GenDecl{
		Tok: token.FUNC,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(commandResponseFuncName),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("Value")},
								Type:  ast.NewIdent(command.ReturnsType()),
								Tag:   &ast.BasicLit{Value: "`json:\"value\"`"},
							},
						},
					},
				},
			},
		},
	}
}

func (command *Command) generateConstructor() *ast.FuncDecl {
	parameters := command.transformArguments(func(argument *Field) *ast.Field {
		return &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(argument.Name())},
			Type:  ast.NewIdent(argument.Type()),
		}
	})
	assignations := make([]*ast.KeyValueExpr, len(parameters))
	for i, param := range parameters.List {
		assignations[i] = &ast.KeyValueExpr{Key: param.Names[0], Value: params.Names[0]}
	}
	commandConstructorName := fmt.Sprintf("New%s", command.safeName())
	return &ast.FuncDecl{
		Name: ast.NewIdent(commandConstructorName),
		Type: &ast.FuncType{
			Params: parameters,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{&ast.CompositeLit{
						Type: ast.NewIdent(command.safeName()),
						Elts: []ast.Expr { &ast.KeyValueExpr{} }
					}},
				},
			},
		},
	}
}

func (command *Command) transformArguments(transformCallback func(*Field) *ast.Field) *ast.FieldList {
	transformedArguments := make([]*ast.Field, len(command.arguments))
	for i, argument := range command.arguments {
		transformedArguments[i] = transformCallback(argument)
	}
	return &ast.FieldList{List: transformedArguments}
}

func (command *Command) Arguments() []*Field {
	return command.arguments
}

func (command *Command) ReturnsType() string {
	var returnsType string
	switch command.returnsType {
	case "string":
	case "int":
	case "float32":
	case "uint16":
	case "bool":
		returnsType = command.returnsType
	default:
		returnsType = "*" + command.returnsType
	}
	return returnsType
}
