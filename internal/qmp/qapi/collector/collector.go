package collector

import (
	"encoding/json"
	"fmt"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/qapi"
)

func Collect(bytes []byte) ([]qapi.Entity, *CollectError) {
	var root []json.RawMessage
	if unmarshalError := json.Unmarshal(bytes, &root); unmarshalError != nil {
		return nil, NewCollectError(unmarshalError, MalformedSchema)
	}
	return processEntities(root)
}

func processEntities(root []json.RawMessage) ([]qapi.Entity, *CollectError) {
	entities := make([]qapi.Entity, len(root))
	for i, rawItem := range root {
		var item map[string]json.RawMessage
		if unmarshalError := json.Unmarshal(rawItem, &item); unmarshalError != nil {
			return nil, NewCollectError(unmarshalError, MalformedSchema)
		}
		var typeDetectionError error
		entities[i], typeDetectionError = detectEntityType(item)
		if typeDetectionError != nil {
			return nil, NewCollectError(typeDetectionError, MalformedSchema)
		}
		if unmarshalError := json.Unmarshal(rawItem, &entities[i]); unmarshalError != nil {
			return nil, NewCollectError(unmarshalError, MalformedSchema)
		}
	}
	return entities, nil
}

func detectEntityType(item map[string]json.RawMessage) (qapi.Entity, error) {
	var entity qapi.Entity
	switch {
	case item["command"] != nil:
		entity = &Command{}
	case item["enum"] != nil:
		entity = &Enum{}
	case item["struct"] != nil:
		entity = &Struct{}
	default:
		return nil, fmt.Errorf("invalid entity type")
	}
	return entity, nil
}
