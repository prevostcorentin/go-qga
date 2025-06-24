package converter

import "strings"

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

func (field *Field) CapitalizedName() string {
	return strings.ToUpper(string(field.name[0])) + field.name[1:len(field.name)]
}

func (field *Field) Type() string {
	switch field.typeName {
	case "str":
		return "string"
	case "int":
		return "int"
	case "number":
		return "float32"
	case "uint16":
		return "uint16"
	case "bool":
		return "bool"
	default:
		return "*" + field.typeName
	}
}
