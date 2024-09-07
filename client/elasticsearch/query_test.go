package elasticsearch

import "testing"

func TestNewRequest(t *testing.T) {
	res := NewRequest()

	t.Log(res.Params())
}
