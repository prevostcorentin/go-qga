package converter_test

import (
	"testing"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
	. "github.com/prevostcorentin/go-qga/internal/qmp/qapi/converter"
)

func TestConvertEntities(t *testing.T) {
	collectedEntities := []collector.Entity{
		&collector.Command{
			CommandName:    "test-command",
			CommandData:    map[string]any{"argument": "str", "enum": "TestEnum"},
			CommandReturns: "TestStruct",
		},
		&collector.Enum{
			EnumName: "TestEnum",
			EnumData: []string{"value1", "value2"},
		},
		&collector.Struct{
			StructName: "TestStruct",
			StructData: map[string]any{"argument": "TestEnum"},
		},
	}
	convertedEntities, conversionError := Convert(collectedEntities)
	if conversionError != nil {
		t.Fatalf("while converting entities: %v", conversionError)
	}
	command := convertedEntities[0]
	commandFields := command.Fields()
	if commandFields[0].Name() != "argument" {
		t.Errorf(`wrong name "%v" for argument 0. expected "argument"`, commandFields[0].Name())
	}
	if _, isPrimitive := commandFields[0].(*Primitive); !isPrimitive {
		t.Errorf("wrong data type for argument 0. expected a primitive type.")
	}
	primitive := commandFields[0].(*Primitive)
	if primitive.TypeName() != string(StringPrimitiveType) {
		t.Errorf(`wrong primitive type "%v" for argument 0. expected "%s"`, string(StringPrimitiveType))
	}
	if _, isEntity := commandFields[1].(*EntityField); !isEntity {
		t.Errorf("wrong data type for argument 1. expected an entity type")
	}
	entity := commandFields[1].(*EntityField)
	if entity.Name() != "TestEnum" {
		t.Errorf(`wrong entity type name "%v" for argument 1. expected "TestEnum"`, entity.Type())
	}
	commandReturns := command.(*Command).Returns()
	if commandReturns.Type() != EntityFieldType {
		t.Errorf("wrong type for command returns. expected an entity")
	}
	entity = commandReturns.(*EntityField)
	if entity.Name() != "TestStruct" {
		t.Errorf(`wrong entity name "%v" for command returns. expected "TestStruct"`)
	}
}
