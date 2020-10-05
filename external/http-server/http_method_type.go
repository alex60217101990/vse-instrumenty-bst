package http_server

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type HttpMethods int

const (
	GET HttpMethods = iota + 1
	POST
)

var (
	_HttpMethodsNameToValue = map[string]HttpMethods{
		"GET":  GET,
		"POST": POST,
		"get":  GET,
		"post": POST,
		"Get":  GET,
		"Post": POST,
	}

	_HttpMethodsValueToName = map[HttpMethods]string{
		GET:  "GET",
		POST: "POST",
	}
)

func (m HttpMethods) MarshalYAML() (interface{}, error) {
	s, ok := _HttpMethodsValueToName[m]
	if !ok {
		return nil, fmt.Errorf("invalid HttpMethods: %d", m)
	}
	return s, nil
}

func (m *HttpMethods) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _HttpMethodsNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid HttpMethods %q", value.Value)
	}
	*m = v
	return nil
}

func (m HttpMethods) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(m).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _HttpMethodsValueToName[m]
	if !ok {
		return nil, fmt.Errorf("invalid HttpMethods: %d", m)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so ShirtSize satisfies json.Unmarshaler.
func (m *HttpMethods) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("HttpMethods should be a string, got %s", data)
	}
	v, ok := _HttpMethodsNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid HttpMethods %q", s)
	}
	*m = v
	return nil
}

func (m HttpMethods) Val() int {
	return int(m)
}

func (m HttpMethods) String() string {
	return _HttpMethodsValueToName[m]
}
