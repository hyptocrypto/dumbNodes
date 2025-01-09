package types

type Request struct {
	Method      string            `json:"method" msgpack:"method"`
	Headers     map[string]string `json:"headers" msgpack:"headers"`
	Destination string            `json:"destination" msgpack:"destination"`
	Data        []byte            `json:"data" msgpack:"data"`
}

type Response struct {
	Headers map[string]string `json:"headers" msgpack:"headers"`
	Source  string            `json:"source" msgpack:"source"`
	Data    []byte            `json:"data" msgpack:"data"`
}
