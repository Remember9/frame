package http

import (
	encoding "github.com/Remember9/frame/util/xencoding"
	"github.com/Remember9/frame/util/xencoding/form"
	"net/http"
)

// BindQuery bind vars parameters to target.
func BindQuery(req *http.Request, target interface{}) error {
	return encoding.GetCodec(form.Name).Unmarshal([]byte(req.URL.Query().Encode()), target)
}

// BindForm bind form parameters to target.
func BindForm(req *http.Request, target interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	return encoding.GetCodec(form.Name).Unmarshal([]byte(req.Form.Encode()), target)
}

func Bind(req *http.Request, target interface{}) error {
	return defaultRequestDecoder(req, target)
}

func EncodeResponse(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return DefaultResponseEncoder(w, r, v)
}
