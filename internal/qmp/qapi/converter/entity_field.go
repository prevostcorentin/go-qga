package converter

type EntityField struct {
	entityName string
	entityType string
}

func (entityField *EntityField) Name() string {
	return entityField.entityName
}

func (entityField *EntityField) TypeName() string {
	return entityField.entityType
}

func (_ *EntityField) Type() FieldType {
	return EntityFieldType
}
