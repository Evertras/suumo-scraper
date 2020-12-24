package esdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type indexMappingProperties struct {
	Properties map[string]indexMappingPropertyType `json:"properties"`
}

type indexMappingPropertyType struct {
	Type string `json:"type"`
}

func genIndexMappingBody(mappings map[string]string) (io.Reader, error) {
	properties := make(map[string]indexMappingPropertyType)

	for key, value := range mappings {
		properties[key] = indexMappingPropertyType{
			Type: value,
		}
	}

	m := indexMappingProperties{
		Properties: properties,
	}

	raw, err := json.Marshal(m)

	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return bytes.NewReader(raw), nil
}
