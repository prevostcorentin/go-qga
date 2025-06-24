package collector

type Entity interface {
	Name() string
	Fields() map[string]any
	Type() string
}

type EntityType string

const (
	CommandType EntityType = "command"
	EnumType               = "enum"
	PragmaType             = "pragma"
	StructType             = "struct"
)

type FieldType string

const (
	IntFieldType     FieldType = "int"
	StringFieldType            = "str"
	BooleanFieldType           = "bool"
)
