package xnet

import "testing"

func TestParseURL(t *testing.T) {
	url, err := ParseURL("http://esf.leju.com/bj/house/?a=b#c")
	if err != nil {
		t.Error(err)
	}
	t.Log(url.QueryString("a", "k"))
	t.Log(url.QueryString("av", "k"))
	t.Log(url.Query())
}
