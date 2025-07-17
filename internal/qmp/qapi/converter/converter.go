package converter

import (
	"fmt"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"
)

func Convert(collectedEntities []qapi.Entity) ([]qapi.Entity, *ConvertError) {
	convertedEntities := make([]qapi.Entity, len(collectedEntities))
	for i, collectedOne := range collectedEntities {
		switch collectedOne.(type) {
		case *collector.Command:
			convertedEntities[i] = NewCommand(collectedOne.(*collector.Command))
		case *collector.Enum:
			convertedEntities[i] = NewEnum(collectedOne.(*collector.Enum))
		case *collector.Struct:
			convertedEntities[i] = NewStruct(collectedOne.(*collector.Struct))
		default:
			wrappedError := fmt.Errorf(`unknown entity type`)
			return convertedEntities, NewConvertError(wrappedError, UnknownEntityType)
		}
	}
	return convertedEntities, nil
}
