package elasticsearch

import (
	"encoding/json"
	"fmt"
)

const (
	// RangeScopeLoRo left open & right open
	RangeScopeLoRo string = "( )"
	// RangeScopeLoRc left open & right close
	RangeScopeLoRc string = "( ]"
	// RangeScopeLcRo left close & right open
	RangeScopeLcRo string = "[ )"
	// RangeScopeLcRc lect close & right close
	RangeScopeLcRc string = "[ ]"

	// LikeLevelHigh wildcard keyword
	LikeLevelHigh string = "high"
	// LikeLevelMiddle ngram(1,2)
	LikeLevelMiddle string = "middle"
	// LikeLevelLow match split word
	LikeLevelLow string = "low"
)

// QueryBody .
type QueryBody struct {
	Fields []string       `json:"fields"` // default:"*" _source，default = *
	Where  QueryBodyWhere `json:"where"`
}

// QueryBodyWhere .
type QueryBodyWhere struct {
	EQ    map[string]interface{}     `json:"eq"` //可能是数据或字符,[12,333,67] ["asd", "你好"]
	Or    map[string]interface{}     `json:"or"` //暂时不支持minimum should
	In    map[string][]interface{}   `json:"in"`
	Range map[string]string          `json:"range"` //[10,20)  (2018-05-10 00:00:00,2018-05-31 00:00:00]  (,30]
	Like  []QueryBodyWhereLike       `json:"like"`
	Combo []*QueryBodyWhereCombo     `json:"combo"` //混合与或
	Not   map[string]map[string]bool `json:"not"`   //对eq、in、range条件取反
}

// QueryBodyWhereLike .
type QueryBodyWhereLike struct {
	KWFields []string `json:"kw_fields"`
	KW       []string `json:"kw"`    //将kw的值使用空白间隔给query
	Or       bool     `json:"or"`    //default:"false"
	Level    string   `json:"level"` //默认default
}

// QueryBodyWhereCombo .
type QueryBodyWhereCombo struct {
	EQ       []map[string]interface{}   `json:"eq"`
	In       []map[string][]interface{} `json:"in"`
	Range    []map[string]string        `json:"range"`
	NotEQ    []map[string]interface{}   `json:"not_eq"`
	NotIn    []map[string][]interface{} `json:"not_in"`
	NotRange []map[string]string        `json:"not_range"`
	Min      struct {
		EQ       int `json:"eq"`
		In       int `json:"in"`
		Range    int `json:"range"`
		NotEQ    int `json:"not_eq"`
		NotIn    int `json:"not_in"`
		NotRange int `json:"not_range"`
		Min      int `json:"min"`
	} `json:"min"`
}

// ComboEQ .
func (cmb *QueryBodyWhereCombo) ComboEQ(eq []map[string]interface{}) *QueryBodyWhereCombo {
	cmb.EQ = append(cmb.EQ, eq...)
	return cmb
}

// ComboRange .
func (cmb *QueryBodyWhereCombo) ComboRange(r []map[string]string) *QueryBodyWhereCombo {
	cmb.Range = append(cmb.Range, r...)
	return cmb
}

// ComboIn .
func (cmb *QueryBodyWhereCombo) ComboIn(in []map[string][]interface{}) *QueryBodyWhereCombo {
	cmb.In = append(cmb.In, in...)
	return cmb
}

// ComboNotEQ .
func (cmb *QueryBodyWhereCombo) ComboNotEQ(eq []map[string]interface{}) *QueryBodyWhereCombo {
	cmb.NotEQ = append(cmb.NotEQ, eq...)
	return cmb
}

// ComboNotRange .
func (cmb *QueryBodyWhereCombo) ComboNotRange(r []map[string]string) *QueryBodyWhereCombo {
	cmb.NotRange = append(cmb.NotRange, r...)
	return cmb
}

// ComboNotIn .
func (cmb *QueryBodyWhereCombo) ComboNotIn(in []map[string][]interface{}) *QueryBodyWhereCombo {
	cmb.NotIn = append(cmb.NotIn, in...)
	return cmb
}

// MinEQ .
func (cmb *QueryBodyWhereCombo) MinEQ(min int) *QueryBodyWhereCombo {
	cmb.Min.EQ = min
	return cmb
}

// MinIn .
func (cmb *QueryBodyWhereCombo) MinIn(min int) *QueryBodyWhereCombo {
	cmb.Min.In = min
	return cmb
}

// MinRange .
func (cmb *QueryBodyWhereCombo) MinRange(min int) *QueryBodyWhereCombo {
	cmb.Min.Range = min
	return cmb
}

// MinNotEQ .
func (cmb *QueryBodyWhereCombo) MinNotEQ(min int) *QueryBodyWhereCombo {
	cmb.Min.NotEQ = min
	return cmb
}

// MinNotIn .
func (cmb *QueryBodyWhereCombo) MinNotIn(min int) *QueryBodyWhereCombo {
	cmb.Min.NotIn = min
	return cmb
}

// MinNotRange .
func (cmb *QueryBodyWhereCombo) MinNotRange(min int) *QueryBodyWhereCombo {
	cmb.Min.NotRange = min
	return cmb
}

// MinAll .
func (cmb *QueryBodyWhereCombo) MinAll(min int) *QueryBodyWhereCombo {
	cmb.Min.Min = min
	return cmb
}

// QueryResult query result.
type QueryResult struct {
	Order  string          `json:"order"`
	Sort   string          `json:"sort"`
	Result json.RawMessage `json:"result"`

	Page *Page `json:"page"`
}

type Page struct {
	Pn    int   `json:"num"`
	Ps    int   `json:"size"`
	Total int64 `json:"total"`
}

// Request request to elastic
type Request struct {
	q *QueryBody
}

// NewRequest new a request every search query
func NewRequest() *Request {
	return &Request{
		q: &QueryBody{
			Fields: []string{},
		},
	}
}

// Fields add query fields
func (r *Request) Fields(fields ...string) *Request {
	r.q.Fields = append(r.q.Fields, fields...)
	return r
}

// WhereEq where equal
func (r *Request) WhereEq(field string, eq interface{}) *Request {
	if r.q.Where.EQ == nil {
		r.q.Where.EQ = make(map[string]interface{})
	}
	r.q.Where.EQ[field] = eq
	return r
}

// WhereOr where or
func (r *Request) WhereOr(field string, or interface{}) *Request {
	if r.q.Where.Or == nil {
		r.q.Where.Or = make(map[string]interface{})
	}
	r.q.Where.Or[field] = or
	return r
}

// WhereIn where in
func (r *Request) WhereIn(field string, in interface{}) *Request {
	if r.q.Where.In == nil {
		r.q.Where.In = make(map[string][]interface{})
	}
	switch v := in.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, string:
		r.q.Where.In[field] = append(r.q.Where.In[field], v)
	case []int:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int64:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []string:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int8:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int16:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int32:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint8:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint16:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint32:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint64:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	}
	return r
}

// WhereRange where range
func (r *Request) WhereRange(field string, start, end interface{}, scope string) *Request {
	if r.q.Where.Range == nil {
		r.q.Where.Range = make(map[string]string)
	}
	if start == nil {
		start = ""
	}
	if end == nil {
		end = ""
	}
	switch scope {
	case RangeScopeLoRo:
		r.q.Where.Range[field] = fmt.Sprintf("(%v,%v)", start, end)
	case RangeScopeLoRc:
		r.q.Where.Range[field] = fmt.Sprintf("(%v,%v]", start, end)
	case RangeScopeLcRo:
		r.q.Where.Range[field] = fmt.Sprintf("[%v,%v)", start, end)
	case RangeScopeLcRc:
		r.q.Where.Range[field] = fmt.Sprintf("[%v,%v]", start, end)
	}
	return r
}

// WhereNot where not
func (r *Request) WhereNot(typ string, fields ...string) *Request {
	if r.q.Where.Not == nil {
		r.q.Where.Not = make(map[string]map[string]bool)
	}
	if r.q.Where.Not[typ] == nil {
		r.q.Where.Not[typ] = make(map[string]bool)
	}
	for _, v := range fields {
		r.q.Where.Not[typ][v] = true
	}
	return r
}

// WhereLike where like
func (r *Request) WhereLike(fields, words []string, or bool, level string) *Request {
	if len(fields) == 0 || len(words) == 0 {
		return r
	}
	l := QueryBodyWhereLike{KWFields: fields, KW: words, Or: or, Level: level}
	r.q.Where.Like = append(r.q.Where.Like, l)
	return r
}

// WhereCombo where combo
func (r *Request) WhereCombo(cmb ...*QueryBodyWhereCombo) *Request {
	r.q.Where.Combo = append(r.q.Where.Combo, cmb...)
	return r
}

// Params get query parameters
func (r *Request) Params() *QueryBody {
	return r.q
}
