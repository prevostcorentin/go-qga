package collector_test

import (
	"testing"

	. "github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
)

func TestReadQapiJson(t *testing.T) {
	qapi_schema := `[
{ "command": "test-command", "data":{ "id": "int" }, "returns": "int" },
{ "struct": "TestStruct", "data": { "value": "int" } },
{ "enum": "TestEnum", "data": [ "value1", "value2" ] }
]`

	entities, collectError := Collect([]byte(qapi_schema))
	if collectError != nil {
		t.Fatalf("while collecting entities: %v", collectError)
	}

	command := entities[0]
	if command.Name() != "test-command" {
		t.Errorf(`entity "test-command" is missing`)
	}
	if command.Type() != string(CommandType) {
		t.Errorf(`wrong entity type "%v". expected "%s"`, command.Type(), CommandType)
	}
	if _, ok := command.Fields()["id"]; !ok {
		t.Errorf(`wrong field name "id" not found.`)
	}
	if command.Fields()["id"] != "int" {
		t.Errorf(`wrong field type "%v" for "%v". expected "int"`, command.Fields()["value"], command.Fields()["value"])
	}
	st := entities[1]
	if st.Name() != "TestStruct" {
		t.Errorf(`wrong entity name "%v". expected "TestStruct"`, st.Name())
	}
	if st.Type() != StructType {
		t.Errorf(`wrong entity type "%v". expected "%s"`, st.Type(), StructType)
	}
	if _, valuePresent := st.Fields()["value"]; !valuePresent {
		t.Errorf(`"value" key is missing from struct`)
	}
	if st.Fields()["value"] == IntFieldType {
		t.Errorf(`wrong field type "%v". expected "%s"`, st.Fields()["value"], IntFieldType)
	}
	en := entities[2]
	if en.Name() != "TestEnum" {
		t.Errorf(`wrong entity name "%v". expected "TestEnum"`, en.Name())
	}
	if en.Type() != EnumType {
		t.Errorf(`wrong entity type "%v". expected "%s"`, en.Type(), EnumType)
	}
	enumData := en.(*Enum).Data()
	if enumData[0] != "value1" {
		t.Errorf(`wrong field value "%v". expected "value1"`, enumData[0])
	}
	if enumData[1] != "value2" {
		t.Errorf(`wrong field value "%v". expected "value2"`, enumData[1])
	}
}
