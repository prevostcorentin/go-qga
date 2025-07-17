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
	if commandArguments[0].Name() != "argument" {
		t.Errorf(`wrong argument name "%v". Expected "argument"`, commandArguments[0].Name())
	}
	if commandArguments[0].Type() != StringFieldType {
		t.Errorf("wrong data type for argument 0. expected an string type")
	}
	if commandArguments[1].Name() != "TestEnum" {
		t.Errorf(`wrong entity type name "%v" for argument 1. expected "TestEnum"`, commandArguments[1].Name())
	}
	if commandArguments[1].Type() != CompositeFieldType {
		t.Errorf("wrong data type for argument 1. expected a compositie type")
	}
	if command.Returns().Name() != "TestStruct" {
		t.Errorf(`wrong entity name "%v" for command returns. expected "TestStruct"`, command.Returns().Name())
	}
	if command.Returns().Type() != CompositeFieldType {
		t.Errorf("wrong type for command returns. expected an entity")
	}
}
