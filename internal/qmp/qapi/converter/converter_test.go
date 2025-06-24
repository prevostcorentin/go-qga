package converter_test

import (
	"testing"

	"github.com/prevostcorentin/go-qga/internal/qmp/qapi"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
	. "github.com/prevostcorentin/go-qga/internal/qmp/qapi/converter"
)

func TestConvertEntities(t *testing.T) {
	collectedEntities := []qapi.Entity{
		&collector.Command{
			CommandName: "test-command",
			Arguments:   map[string]string{"argument": "str", "enum": "TestEnum"},
			Returns:     "TestStruct",
		},
		&collector.Enum{
			EnumName: "TestEnum",
			Data:     []string{"value1", "value2"},
		},
		&collector.Struct{
			StructName: "TestStruct",
			Data:       map[string]any{"argument": "TestEnum"},
		},
	}
	convertedEntities, conversionError := Convert(collectedEntities)
	if conversionError != nil {
		t.Fatalf("while converting entities: %v", conversionError)
	}
	if convertedEntities[0].Name() != "test-command" {
		t.Errorf(`wrong name "%v" for argument 0. expected "test-command"`, convertedEntities[0].Name())
	}
	var ok bool
	var command *Command
	if command, ok = convertedEntities[0].(*Command); !ok {
		t.Errorf(`wrong data type. expected "*Command"`)
	}
	commandArguments := command.Arguments()
	if commandArguments[0].Name() != "Argument" {
		t.Errorf(`Got first argument name "%v" when expecting "Argument"`, commandArguments[0].Name())
	}
	if commandArguments[0].Type() != "string" {
		t.Errorf(`Got first argument type "%v" when expecting "string"`, commandArguments[0].Type())
	}
	if commandArguments[1].Name() != "Enum" {
		t.Errorf(`Got second argument name "%v" when expecting "Enum"`, commandArguments[1].Name())
	}
	if commandArguments[1].Type() != "TestEnum" {
		t.Errorf(`Got second argument type "%v" when expecting "TestEnum"`, commandArguments[1].Type())
	}
	if command.ReturnsType() != "TestStruct" {
		t.Errorf(`Got return type "%v" when expecting "TestStruct"`, command.ReturnsType())
	}
}

func TestGenerateEnum(t *testing.T) {
	testEnum := NewEnum(&collector.Enum{
		EnumName: "TestEnum",
		Data:     []string{"value1", "value2"},
	})

	expectedCode := `package generated

type TestEnum string

const (
	value1 TestEnum = "value1"
	value2 TestEnum = "value2"
)
`
	bytes, err := testEnum.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != expectedCode {
		t.Logf("\n------------ Expected ------------\n%s\n------------ Generated ------------\n%s\n", expectedCode, string(bytes))
		t.Fatalf("generated code does not match what is expected")
	}
}

func TestGenerateStruct(t *testing.T) {
	testStruct := NewStruct(&collector.Struct{
		StructName: "TestStruct",
		Data:       map[string]any{"argument": "str"},
	})

	expectedCode := `package generated

type TestStruct struct {
	Argument string
}
`

	bytes, err := testStruct.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != expectedCode {
		t.Logf("\n------------ Expected ------------\n%s\n------------ Generated ------------\n%s\n", expectedCode, string(bytes))
		t.Fatalf("generated code does not match what is expected")
	}
}

func TestGenerateCommand(t *testing.T) {
	testCommand := NewCommand(&collector.Command{
		CommandName: "test-command",
		Arguments:   map[string]string{"argument": "TestEnum"},
		Returns:     "TestStruct",
	})

	expectedCode := `package generated

type TestCommand struct {
	argument *TestEnum
}

type testCommandArguments struct {
	Argument *TestEnum ` + "`json:\"argument\"`" + `
}

type testCommandResponse struct {
	Value *TestStruct ` + "`json:\"value\"`" + `
}

func NewTestCommand(argument *TestEnum) *TestCommand {
	return &TestCommand{ argument: argument }
}

func (command *TestCommand) Execute() string {
	return "test-command"
}

func (command *TestCommand) Arguments() any {
	return &testCommandArguments {
		Argument: command.argument
	}
}

func (command *TestCommand) Response() any {
	return &testCommandResponse{}
}
`

	bytes, err := testCommand.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != expectedCode {
		t.Logf("\n------------ Expected ------------\n%s\n------------ Generated ------------\n%s\n", expectedCode, string(bytes))
		t.Fatalf("generated code does not match what is expected")
	}
}
