package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

func JSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(response)
}

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	//jsoni.RegisterExtension(new(jsoniterator.EmitDefaultExtension))
}

func JSONiterResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	return jsoni.NewEncoder(w).Encode(response)
}
func XMLResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(200)
	return xml.NewEncoder(w).Encode(response)
}

func StringResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, err = w.Write([]byte(response.(string)))
	if err != nil {
		return err
	}
	return nil
}

func ReturnResponse(code int, contentType string) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, data interface{}) error {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(code)
		_, err := w.Write(data.([]byte))
		if err != nil {
			return err
		}
		return nil
	}
}

func ExcelResponse(fileName string) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, data interface{}) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Disposition", "attachment;filename="+fileName)
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		w.WriteHeader(200)
		_, err := w.Write(data.([]byte))
		if err != nil {
			return err
		}
		return nil
	}
}
