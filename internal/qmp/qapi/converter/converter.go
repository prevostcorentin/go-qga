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
		switch collectedOne := collectedOne.(type) {
		case *collector.Command:
			convertedEntities[i] = NewCommand(collectedOne)
		case *collector.Enum:
			convertedEntities[i] = NewEnum(collectedOne)
		case *collector.Struct:
			convertedEntities[i] = NewStruct(collectedOne)
		default:
			wrappedError := fmt.Errorf(`unknown entity type`)
			return convertedEntities, NewConvertError(wrappedError, UnknownEntityType)
		}
	}
	return convertedEntities, nil
}
