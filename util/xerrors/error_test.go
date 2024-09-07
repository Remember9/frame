package xerrors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	t.Log(New(1001, "自定义增加错误码的错误处理"))
}

func TestToCodeError(t *testing.T) {
	err := New(1001, "自定义增加错误码的错误处理")
	r := ToCodeError(err)
	t.Log(r.Code(), r.Error())

	err = errors.New("原生错误处理")
	r = ToCodeError(err)
	t.Log(r.Code(), r.Error())
}
