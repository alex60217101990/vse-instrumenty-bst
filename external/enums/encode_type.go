package enums

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type EncodingType int

const (
	Json EncodingType = iota + 1
	Yaml
)

var (
	_EncodingTypeNameToValue = map[string]EncodingType{
		"json": Json,
		"Json": Json,
		"Yaml": Yaml,
		"yaml": Yaml,
	}

	_EncodingTypeValueToName = map[EncodingType]string{
		Json: "json",
		Yaml: "yaml",
	}
)

func (r EncodingType) MarshalYAML() (interface{}, error) {
	s, ok := _EncodingTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid EncodingType: %d", r)
	}
	return s, nil
}

func (r *EncodingType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _EncodingTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid EventType %q", value.Value)
	}
	*r = v
	return nil
}

func (r EncodingType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _EncodingTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid EncodingType: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so ShirtSize satisfies json.Unmarshaler.
func (r *EncodingType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("EventType should be a string, got %s", data)
	}
	v, ok := _EncodingTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid EncodingType %q", s)
	}
	*r = v
	return nil
}

func (r EncodingType) Val() int {
	return int(r)
}

func (r EncodingType) String() string {
	return _EncodingTypeValueToName[r]
}

func (r *EncodingType) FromString(str string) {
	if val, ok := _EncodingTypeNameToValue[str]; ok {
		*r = val
	} else {
		*r = Yaml
	}
}
