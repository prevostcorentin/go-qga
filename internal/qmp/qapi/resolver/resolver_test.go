package resolver_test

import (
	"testing"

	"github.com/prevostcorentin/go-qga/internal/qmp/qapi"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/converter"
	. "github.com/prevostcorentin/go-qga/internal/qmp/qapi/resolver"
)

func TestConvertEntities(t *testing.T) {
	convertedEntities := []qapi.Entity{
		converter.NewCommand(&collector.Command{
			CommandName: "test-command",
			Arguments:   map[string]string{"argument": "str", "enum": "TestEnum"},
			Returns:     "TestStruct",
		}),
		converter.NewEnum(&collector.Enum{
			EnumName: "TestEnum",
			Data:     []string{"value1", "value2"},
		}),
		converter.NewStruct(&collector.Struct{
			StructName: "TestStruct",
			Data:       map[string]any{"argument": "TestEnum"},
		}),
	}
	resolvedCommands, resolveError := Resolve(convertedEntities)
	if resolveError != nil {
		t.Fatalf("while resolving entities: %v", resolveError)
	}
	command := resolvedCommands[0]
	if command.Name() != "test-command" {
		t.Errorf(`wrong name "%v" for command". expected "test-command".`)
	}
	arguments := command.Arguments()
	var ok bool
	var primitive *Primitive
	if ok, primitive = arguments["argument"].(*Primitive); !ok {
		t.Errorf("wrong argument type. should be a Primitive")
	}
	if primitive.Type() != StringType {
		t.Errorf(`wrong primitive type "%v". shoud be "%v"`, primitive.Type(), StringType)
	}
	var composite *Composite
	if composite, ok = arguments["enum"].(*Composite); !ok {
		t.Errorf(`wrong data type. expected "Composite"`)
	}
	if composite.Name() != "TestEnum" {
		t.Errorf(`wrong composite name "%v". expected "TestEnum".`, composite.Name())
	}
	var en *Enum
	if en, ok = composite.(*Enum); !ok {
		t.Errorf(`wrong composite sub-type. expected "Enum"`)
	}
	compositeData := composite.Data()
	if compositeData[0] != "value1" || compositeData[1] != "value2" {
		t.Errorf(`wrong enum data ["%v", "%v"]. expected ["value1", "value2"]`)
	}
	returns := command.Returns()
	if ok, composite = returns.(*Composite); !ok {
		t.Errorf("wrong returns type. expected Composite")
	}
	var st *Struct
	if ok, st = composite.(*Struct); !ok {
		t.Errorf("wrong composite type. expected Struct")
	}
	if st.Name() != "TestStruct" {
		t.Errorf(`wrong struct name "%v". expected "TestStruct"`)
	}
	structData := st.Data()
	if ok, composite = structData["argument"].(*Composite); !ok {
		t.Errorf(`wrong data argument data type. expected "Composite"`)
	}
	if structData.Data["argument"] != en {
		t.Errorf("enum from struct should be the same as enum from command arguments")
	}
}
