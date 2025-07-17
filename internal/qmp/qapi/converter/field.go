package converter

type Field struct {
	name     string
	typeName string
}

func NewField(name string, typeName string) *Field {
	return &Field{name: name, typeName: typeName}
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) Type() string {
	return field.typeName
}

type FieldType string

const (
	IntegerFieldType   FieldType = "int"
	StringFieldType              = "str"
	BoolFieldType                = "bool"
	CompositeFieldType           = "to_resolve"
)
