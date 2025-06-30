package converter

import "github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"

func Convert(collectedEntities []collector.Entity) ([]Entity, *ConvertError) {
	convertedEntities := make([]Entity, len(collectedEntities))
	for i, collectedOne := range collectedEntities {
		convertedEntities[i] := 
		for _, field := range collectedOne.Fields() {
			if field != collector.StringFieldType ||
			   field != collector.BooleanFieldType ||
			   field != collector.IntFieldType { // This is an entity here
				
			}
		}
	}
}
