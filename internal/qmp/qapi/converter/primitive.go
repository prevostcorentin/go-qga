package converter

type Primitive struct {
	primitiveName string
	primitiveType PrimitiveType
}

func (primitive *Primitive) Name() string {
	return primitive.primitiveName
}

func (primitive *Primitive) TypeName() string {
	return string(primitive.primitiveType)
}

func (_ *Primitive) Type() FieldType {
	return PrimitiveFieldType
}

type PrimitiveType string

const (
	BooleanPrimitiveType PrimitiveType = "boolean"
	IntPrimitiveType                   = "integer"
	StringPrimitiveType                = "string"
)
