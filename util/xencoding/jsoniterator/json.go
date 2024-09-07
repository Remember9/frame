package jsoniterator

import (
	json0 "encoding/json"
	encoding "github.com/Remember9/frame/util/xencoding"
	jsoniter "github.com/json-iterator/go"
)

// Name is the name registered for the json codec.
const Name = "json"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

/*
var (

	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

)
*/
func init() {
	encoding.RegisterCodec(codec{})
	// json.RegisterExtension(new(EmitDefaultExtension))
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json0.Marshaler:
		return m.MarshalJSON()
	/*case proto.Message:
	return MarshalOptions.Marshal(m)*/
	default:
		return json.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json0.Unmarshaler:
		return m.UnmarshalJSON(data)
	/*case proto.Message:
	return UnmarshalOptions.Unmarshal(data, m)*/
	default:
		/*rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}*/
		return json.Unmarshal(data, m)
	}
}

func (codec) Name() string {
	return Name
}
