package converter

type Entity interface {
	Name() string
	Fields() []Field
	Type() EntityType
}

type EntityType string

const (
	CommandEntityType EntityType = "command"
	EnumEntityType               = "enum"
	StructEntityType             = "struct"
)
