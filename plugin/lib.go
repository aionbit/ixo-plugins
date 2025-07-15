package plugin

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

func DecodeInput(input, param any) error {
	return mapstructure.Decode(input, param)
}

func EncodeOutput(output any) (any, error) {
	d, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	var v any
	return v, json.Unmarshal(d, &v)
}
