package config

import (
	"testing"
)

func init() {
	err := InitTest()
	if err != nil {
		panic(err)
	}
}

func TestGetConfig(t *testing.T) {
	c := GetConfig()
	t.Log(c.GetIntSlice("test.int-slice"))
}

func TestGet(t *testing.T) {
	res := Get("app.env")
	t.Logf("%v - %T", res, res)

	res = GetString("test.string")
	t.Logf("%v - %T", res, res)

	res = GetStringSlice("test.string-slice")
	t.Logf("%v - %T", res, res)

	res = GetBool("test.bool")
	t.Logf("%v - %T", res, res)

	res = GetInt("test.int")
	t.Logf("%v - %T", res, res)

	res = GetInt32("test.int32")
	t.Logf("%v - %T", res, res)

	res = GetInt64("test.int64")
	t.Logf("%v - %T", res, res)

	res = GetIntSlice("test.int-slice")
	t.Logf("%v - %T", res, res)

	res = GetFloat64("test.float64")
	t.Logf("%v - %T", res, res)

	res = GetTime("test.time")
	t.Logf("%v - %T", res, res)

	res = GetDuration("test.duration")
	t.Logf("%v - %T", res, res)
}

func TestUnmarshal(t *testing.T) {
	var c *AppConfig
	err := UnmarshalKey("app", &c)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)

	err = Unmarshal(&c)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)

}

func TestDebug(t *testing.T) {
	Debug()
}
