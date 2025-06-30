package converter

type Field interface {
	Name() string
	Type() FieldType
	TypeName() string
}

type FieldType string

const (
	PrimitiveFieldType FieldType = "primitive"
	EntityFieldType              = "entity"
)
