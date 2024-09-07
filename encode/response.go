package encode

import (
	"context"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"esfgit.leju.com/golang/frame/util/xcast"
	"esfgit.leju.com/golang/frame/util/xerrors"
	"net/http"
)

type Response struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Failer interface {
	Failed() error
}

func JsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(Failer); ok && f.Failed() != nil {
		JsonError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(&Response{
		Code: SuccessCode,
		Msg:  "success",
		Data: response,
	})
}

func JsonError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case HealthError:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}

	var (
		code int64
		msg  string
	)
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		code = xcast.ToInt64(mysqlError.Number)
		if mysqlError.Number == 1062 {
			msg = "此数据已经存在"
		} else {
			msg = mysqlError.Message
		}
	} else {
		xerr := xerrors.ToCodeError(err)
		code = xerr.Code()
		msg = xerr.Error()
	}
	_ = json.NewEncoder(w).Encode(Response{
		Code: code,
		Msg:  msg,
	})
}
