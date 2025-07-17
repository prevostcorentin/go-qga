package collector_test

import (
	"testing"

	. "github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
)

func TestReadQapiJson(t *testing.T) {
	qapi_schema := `[
{ "command": "test-command", "data":{ "id": "int", "enum": "TestEnum" }, "returns": "TestStruct" },
{ "struct": "TestStruct", "data": { "argument": "str" } },
{ "enum": "TestEnum", "data": [ "value1", "value2" ] }
]`

	entities, collectError := Collect([]byte(qapi_schema))
	if collectError != nil {
		t.Fatalf("while collecting entities: %v", collectError)
	}

	if entities[0].Name() != "test-command" {
		t.Errorf(`entity "test-command" is missing`)
	}
	var ok bool
	var command *Command
	if command, ok = entities[0].(*Command); !ok {
		t.Errorf(`wrong entity type. expected "*Command"`)
	}
	if _, ok := command.Arguments["id"]; !ok {
		t.Errorf(`wrong field name "id" not found.`)
	}
	if command.Arguments["id"] != "int" {
		t.Errorf(`wrong field type "%v" for "id". expected "int"`, command.Arguments["id"])
	}
	if command.Arguments["enum"] != "TestEnum" {
		t.Errorf(`wrong field type "%v" for "enum". expected "TestStruct"`, command.Arguments["enum"])
	}
	if command.Returns != "TestStruct" {
		t.Errorf(`wrong retunrs type "%v". expected "TestStruct"`, command.Returns)
	}
	if entities[1].Name() != "TestStruct" {
		t.Errorf(`wrong entity name "%v". expected "TestStruct"`, entities[1].Name())
	}
	var st *Struct
	if st, ok = entities[1].(*Struct); !ok {
		t.Errorf(`wrong entity type. expected "*Struct"`)
	}
	if _, valuePresent := st.Data["argument"]; !valuePresent {
		t.Errorf(`"argument" key is missing from struct`)
	}
	if st.Data["argument"] != "str" {
		t.Errorf(`wrong field type "%v". expected "string"`, st.Data["argument"])
	}
	if entities[2].Name() != "TestEnum" {
		t.Errorf(`wrong entity name "%v". expected "TestEnum"`, entities[2].Name())
	}
	var en *Enum
	if en, ok = entities[2].(*Enum); !ok {
		t.Errorf(`wrong entity type. expected "*Enum"`)
	}
	if en.Data[0] != "value1" {
		t.Errorf(`wrong field value "%v". expected "value1"`, en.Data[0])
	}
	if en.Data[1] != "value2" {
		t.Errorf(`wrong field value "%v". expected "value2"`, en.Data[1])
	}
}
