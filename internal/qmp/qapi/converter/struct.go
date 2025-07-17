package converter

import "github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"

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
