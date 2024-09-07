package elasticsearch

import (
	"context"
	"fmt"
	"github.com/Remember9/frame/encode"
	"github.com/Remember9/frame/xlog"
	"github.com/olivere/elastic/v7"
	"net"
	"net/http"
	"strings"
	"time"
)

var _ GenericClient = (*elasticV7)(nil)

type (
	// elasticV7 implements Client
	elasticV7 struct {
		client *elastic.Client
		conf   *Config
	}

	// SearchParametersV7 holds all required and optional parameters for executing a search
	SearchParametersV7 struct {
		Index       string
		Query       elastic.Query
		From        int
		PageSize    int
		Sorter      []elastic.Sorter
		Aggregation []map[string]elastic.Aggregation
		SearchAfter []interface{}
	}

	queryParametersV7 struct {
		field     string
		whereKind string
		esQuery   elastic.Query
	}
	UpdateParamsV7 struct {
		DetectNoop      bool
		RetryOnConflict int
	}
	UpsertParamsV7 struct {
		Insert          bool
		RetryOnConflict int
		UpsertBody      []BulkUpsertItemV7
	}
	BulkUpsertItemV7 struct {
		IndexName string
		IndexID   string
		Doc       MapDataV7
	}
	BulkIndexItemV7 struct {
		IndexName string
		IndexID   string
		OpType    string // create,index
		Doc       MapDataV7
	}
	BulkDeleteItemV7 struct {
		IndexName string
		IndexID   string
	}
	MapDataV7 map[string]interface{}
)

// returns a new implementation of GenericClient
func NewV7Client(conf *Config, clientOptFuncs ...elastic.ClientOptionFunc) (GenericClient, error) {
	clientOptFuncs = append(clientOptFuncs,
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				DialContext:         (&net.Dialer{Timeout: conf.TransportTimeout}).DialContext,
				MaxIdleConns:        conf.TransportMaxIdleConns,
				MaxIdleConnsPerHost: conf.TransportMaxIdleConnsPerHost,
			},
			Timeout: conf.Timeout,
		}),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(128*time.Millisecond, 513*time.Millisecond))),
		elastic.SetURL(conf.URL...),
		elastic.SetDecoder(&elastic.NumberDecoder{}), // critical to ensure decode of int64 won't lose precise
		elastic.SetErrorLog(&WrapErrorLogger{}),
		elastic.SetInfoLog(&WrapInfoLogger{}),
	)
	if conf.Username != "" && conf.Password != "" {
		clientOptFuncs = append(clientOptFuncs, elastic.SetBasicAuth(conf.Username, conf.Password))
	}
	if conf.DisableSniff {
		clientOptFuncs = append(clientOptFuncs, elastic.SetSniff(false))
	}
	if conf.Compress && !conf.Debug {
		clientOptFuncs = append(clientOptFuncs, elastic.SetGzip(true))
	}
	if conf.Debug {
		clientOptFuncs = append(clientOptFuncs, elastic.SetTraceLog(&WrapTraceLogger{}))
	}
	client, err := elastic.NewClient(clientOptFuncs...)
	if err != nil {
		return nil, err
	}

	return &elasticV7{
		client: client,
		conf:   conf,
	}, nil
}

func (c *elasticV7) Client() *elastic.Client {
	return c.client
}

func (c *elasticV7) CreateIndex(ctx context.Context, index, mapping string) error {
	exists, err := c.client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		return encode.DataExistsError
	}
	_, err = c.client.CreateIndex(index).Body(mapping).Do(ctx)
	return err
}

func (c *elasticV7) DeleteIndex(ctx context.Context, index string) error {
	_, err := c.client.DeleteIndex(index).Do(ctx)
	return err
}

func (c *elasticV7) PutMapping(ctx context.Context, index, root, key, valueType string) error {
	body := make(map[string]interface{})
	if len(root) != 0 {
		body["properties"] = map[string]interface{}{
			root: map[string]interface{}{
				"properties": map[string]interface{}{
					key: map[string]interface{}{
						"type": valueType,
					},
				},
			},
		}
	} else {
		body["properties"] = map[string]interface{}{
			key: map[string]interface{}{
				"type": valueType,
			},
		}
	}
	_, err := c.client.PutMapping().Index(index).BodyJson(body).Do(ctx)
	return err
}

func (c *elasticV7) Index(ctx context.Context, index, id, bodyJson string) (string, error) {
	put, err := c.client.Index().Index(index).Id(id).BodyJson(bodyJson).Do(ctx)
	if err != nil {
		return "", err
	}
	return put.Result, err
}

func (c *elasticV7) Delete(ctx context.Context, index, id string) (string, error) {
	del, err := c.client.Delete().Index(index).Id(id).Do(ctx)
	if err != nil {
		return "", err
	}
	return del.Result, err
}

func (c *elasticV7) Update(ctx context.Context, index, id string, doc map[string]interface{}, p *UpdateParamsV7) (string, error) {
	updateRequest := c.client.Update().Index(index).Id(id).Doc(doc)
	if p != nil && p.DetectNoop {
		updateRequest.DetectNoop(true)
	}
	if p != nil && p.RetryOnConflict > 0 {
		updateRequest.RetryOnConflict(p.RetryOnConflict)
	} else {
		updateRequest.RetryOnConflict(15)
	}
	res, err := updateRequest.Do(ctx)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

func (c *elasticV7) BulkIndex(ctx context.Context, bulkData []BulkIndexItemV7) (*elastic.BulkResponse, error) {
	bulkRequest := c.client.Bulk()
	for _, b := range bulkData {
		request := elastic.NewBulkIndexRequest().Index(b.IndexName).Id(b.IndexID).Doc(b.Doc)
		if b.OpType == "create" || b.OpType == "index" {
			request.OpType(b.OpType)
		}
		bulkRequest.Add(request)
	}
	return bulkRequest.Do(ctx)
}

func (c *elasticV7) BulkDelete(ctx context.Context, bulkData []BulkDeleteItemV7) (*elastic.BulkResponse, error) {
	bulkRequest := c.client.Bulk()
	for _, b := range bulkData {
		request := elastic.NewBulkDeleteRequest().Index(b.IndexName).Id(b.IndexID)
		bulkRequest.Add(request)
	}
	return bulkRequest.Do(ctx)
}

func (c *elasticV7) BulkUpsert(ctx context.Context, upsertData *UpsertParamsV7) (*elastic.BulkResponse, error) {
	bulkRequest := c.client.Bulk()
	for _, b := range upsertData.UpsertBody {
		request := elastic.NewBulkUpdateRequest().Index(b.IndexName).Id(b.IndexID).Doc(b.Doc)
		if upsertData.Insert {
			request.DocAsUpsert(true)
		}
		if upsertData.RetryOnConflict > 0 {
			request.RetryOnConflict(upsertData.RetryOnConflict)
		} else {
			request.RetryOnConflict(15)
		}
		bulkRequest.Add(request)
	}
	return bulkRequest.Do(ctx)
}

func (c *elasticV7) CountByQuery(ctx context.Context, index, query string) (int64, error) {
	return c.client.Count(index).BodyString(query).Do(ctx)
}

func (c *elasticV7) Search(ctx context.Context, p *SearchParametersV7) (*elastic.SearchResult, error) {
	searchService := c.client.Search(p.Index).
		Query(p.Query).
		From(p.From).
		SortBy(p.Sorter...)

	if p.PageSize != 0 {
		searchService.Size(p.PageSize)
	}
	if len(p.Aggregation) > 0 {
		for _, aggMap := range p.Aggregation {
			for name, agg := range aggMap {
				searchService.Aggregation(name, agg)
			}
		}
	}
	if len(p.SearchAfter) != 0 {
		searchService.SearchAfter(p.SearchAfter...)
	}
	return searchService.Do(ctx)
}

func (c *elasticV7) SearchWithDSL(ctx context.Context, index, query string) (*elastic.SearchResult, error) {
	return c.client.Search(index).Source(query).Do(ctx)
}

func (c *elasticV7) QueryBody(sp *QueryBody) (mixedQuery *elastic.BoolQuery) {
	mixedQuery = elastic.NewBoolQuery()
	var queries []*queryParametersV7
	nestedQueries := map[string]*elastic.BoolQuery{} // key: path  value: boolQuery
	// fields
	if len(sp.Fields) == 0 {
		sp.Fields = []string{}
	}
	// where
	if &sp.Where == nil {
		sp.Where = QueryBodyWhere{} // 要给个默认值
	}
	// where - eq
	for k, v := range sp.Where.EQ {
		queries = append(queries, &queryParametersV7{
			field:     k,
			whereKind: "eq",
			esQuery:   elastic.NewTermQuery(k, v),
		})
	}
	// where - or
	for k, v := range sp.Where.Or {
		queries = append(queries, &queryParametersV7{
			field:     k,
			whereKind: "or",
			esQuery:   elastic.NewTermQuery(k, v),
		})
	}
	// where - in
	for k, v := range sp.Where.In {
		if len(v) > 1024 {
			xlog.Errorf("where in 超过1024 error(%v)", v)
			continue
		}
		queries = append(queries, &queryParametersV7{
			field:     k,
			whereKind: "in",
			esQuery:   elastic.NewTermsQuery(k, v...),
		})
	}
	// where - range
	ranges, err := c.queryBasicRange(sp.Where.Range)
	if err != nil {
		xlog.Error("Es", xlog.FieldErr(err))
	}
	for k, v := range ranges {
		queries = append(queries, &queryParametersV7{
			field:     k,
			whereKind: "range",
			esQuery:   v,
		})
	}
	// where - combo
	for _, v := range sp.Where.Combo {
		// 外面用bool+should+minimum包裹
		combo := elastic.NewBoolQuery()
		// 里面每个子项也是bool+should+minimum
		cmbEQ := elastic.NewBoolQuery()
		cmbIn := elastic.NewBoolQuery()
		cmbRange := elastic.NewBoolQuery()
		cmbNotEQ := elastic.NewBoolQuery()
		cmbNotIn := elastic.NewBoolQuery()
		cmbNotRange := elastic.NewBoolQuery()
		// 所有的minimum
		if v.Min.Min == 0 {
			v.Min.Min = 1
		}
		if v.Min.EQ == 0 {
			v.Min.EQ = 1
		}
		if v.Min.In == 0 {
			v.Min.In = 1
		}
		if v.Min.Range == 0 {
			v.Min.Range = 1
		}
		if v.Min.NotEQ == 0 {
			v.Min.NotEQ = 1
		}
		if v.Min.NotIn == 0 {
			v.Min.NotIn = 1
		}
		if v.Min.NotRange == 0 {
			v.Min.NotRange = 1
		}
		// 子项should
		for _, vEQ := range v.EQ {
			for eqK, eqV := range vEQ {
				cmbEQ.Should(elastic.NewTermQuery(eqK, eqV))
			}
		}
		for _, vIn := range v.In {
			for inK, inV := range vIn {
				cmbIn.Should(elastic.NewTermsQuery(inK, inV...))
			}
		}
		for _, vRange := range v.Range {
			ranges, _ := c.queryBasicRange(vRange)
			for _, rangeV := range ranges {
				cmbRange.Should(rangeV)
			}
		}
		for _, notEQ := range v.NotEQ {
			for k, v := range notEQ {
				cmbNotEQ.Should(elastic.NewTermQuery(k, v))
			}
		}
		for _, notIn := range v.NotIn {
			for k, v := range notIn {
				cmbNotIn.Should(elastic.NewTermsQuery(k, v...))
			}
		}
		for _, notRange := range v.NotRange {
			ranges, _ := c.queryBasicRange(notRange)
			for _, v := range ranges {
				cmbNotRange.Should(v)
			}
		}
		// 子项minimum
		if len(v.EQ) > 0 {
			combo.Should(cmbEQ.MinimumNumberShouldMatch(v.Min.EQ))
		}
		if len(v.In) > 0 {
			combo.Should(cmbIn.MinimumNumberShouldMatch(v.Min.In))
		}
		if len(v.Range) > 0 {
			combo.Should(cmbRange.MinimumNumberShouldMatch(v.Min.Range))
		}
		if len(v.NotEQ) > 0 {
			combo.MustNot(elastic.NewBoolQuery().Should(cmbNotEQ.MinimumNumberShouldMatch(v.Min.NotEQ)))
		}
		if len(v.NotIn) > 0 {
			combo.MustNot(elastic.NewBoolQuery().Should(cmbNotIn.MinimumNumberShouldMatch(v.Min.NotIn)))
		}
		if len(v.NotRange) > 0 {
			combo.MustNot(elastic.NewBoolQuery().Should(cmbNotRange.MinimumNumberShouldMatch(v.Min.NotRange)))
		}
		// 合并子项
		mixedQuery.Filter(combo.MinimumNumberShouldMatch(v.Min.Min))
	}
	// where - like
	like, err := c.queryBasicLike(sp.Where.Like)
	if err != nil {
		xlog.Error("Es", xlog.FieldErr(err))
	}
	for _, v := range like {
		queries = append(queries, &queryParametersV7{
			whereKind: "like",
			esQuery:   v,
		})
	}
	// mixedQuery
	for _, q := range queries {
		// like  TODO like的map型字段也要支持must not和 nested
		if q.field == "" && q.whereKind == "like" {
			mixedQuery.Must(q.esQuery)
			continue
		}
		if q.field == "" {
			continue
		}
		// prepare nested 一个DSL只能出现一个nested，不然会有问题
		if mapField := strings.Split(q.field, "."); len(mapField) > 1 && mapField[0] != "" {
			if _, ok := nestedQueries[mapField[0]]; !ok {
				nestedQueries[mapField[0]] = elastic.NewBoolQuery()
			}
			if bl, ok := sp.Where.Not[q.whereKind][q.field]; ok && bl {
				nestedQueries[mapField[0]].MustNot(q.esQuery)
				continue
			}
			nestedQueries[mapField[0]].Must(q.esQuery)
			continue
		}
		// must not
		if bl, ok := sp.Where.Not[q.whereKind][q.field]; ok && bl {
			mixedQuery.MustNot(q.esQuery)
			continue
		}
		// should
		if q.whereKind == "or" {
			mixedQuery.Should(q.esQuery)
			mixedQuery.MinimumShouldMatch("1") // 暂时为1
			continue
		}
		// default
		mixedQuery.Filter(q.esQuery)
	}
	// insert nested
	for k, n := range nestedQueries {
		mixedQuery.Must(elastic.NewNestedQuery(k, n))
	}

	return
}

func (c *elasticV7) queryBasicRange(rangeMap map[string]string) (rangeQuery map[string]*elastic.RangeQuery, err error) {
	rangeQuery = make(map[string]*elastic.RangeQuery)
	for k, v := range rangeMap {
		if r := strings.Trim(v, " "); r != "" {
			if rs := []rune(r); len(rs) > 3 {
				firstStr := string(rs[0:1])
				endStr := string(rs[len(rs)-1:])
				rangeStr := strings.Trim(v, "[]() ")
				FromTo := strings.Split(rangeStr, ",")
				if len(FromTo) != 2 {
					err = fmt.Errorf("sp.Where.Range Fromto err")
					continue
				}
				rQuery := elastic.NewRangeQuery(k)
				rc := 0
				if firstStr == "(" && strings.Trim(FromTo[0], " ") != "" {
					rQuery.Gt(strings.Trim(FromTo[0], " "))
					rc++
				}
				if firstStr == "[" && strings.Trim(FromTo[0], " ") != "" {
					rQuery.Gte(strings.Trim(FromTo[0], " "))
					rc++
				}
				if endStr == ")" && strings.Trim(FromTo[1], " ") != "" {
					rQuery.Lt(strings.Trim(FromTo[1], " "))
					rc++
				}
				if endStr == "]" && strings.Trim(FromTo[1], " ") != "" {
					rQuery.Lte(strings.Trim(FromTo[1], " "))
					rc++
				}
				if rc == 0 {
					continue
				}
				rangeQuery[k] = rQuery
			} else {
				// 范围格式有问题
				err = fmt.Errorf("sp.Where.Range range format err. error(%v)", v)
				continue
			}
		}
	}
	return
}

func (c *elasticV7) queryBasicLike(likeMap []QueryBodyWhereLike) (likeQuery []elastic.Query, err error) {
	for _, v := range likeMap {
		if len(v.KW) == 0 {
			continue
		}
		switch v.Level {
		case LikeLevelHigh:
			var kw []string
			r := []rune(v.KW[0])
			for i := 0; i < len(r); i++ {
				if k := string(r[i : i+1]); !strings.ContainsAny(k, "~[](){}^?:\"\\/!+-=&* ") { // 去掉特殊符号
					kw = append(kw, k)
				} else if len(kw) > 1 && kw[len(kw)-1:][0] != "*" {
					kw = append(kw, "*", " ", "*")
				}
			}
			if len(kw) == 0 || strings.Join(kw, "") == "* *" {
				continue
			}
			qs := elastic.NewQueryStringQuery("*" + strings.Trim(strings.Join(kw, ""), "* ") + "*").AllowLeadingWildcard(true) // 默认是or
			if !v.Or {
				qs.DefaultOperator("AND")
			}
			for _, v := range v.KWFields {
				qs.Field(v)
			}
			likeQuery = append(likeQuery, qs)
		case LikeLevelMiddle:
			// 单个字要特殊处理
			if r := []rune(v.KW[0]); len(r) == 1 && len(v.KW) == 1 {
				qs := elastic.NewQueryStringQuery("*" + string(r[:]) + "*").AllowLeadingWildcard(true) // 默认是or
				if !v.Or {
					qs.DefaultOperator("AND")
				}
				for _, v := range v.KWFields {
					qs.Field(v)
				}
				likeQuery = append(likeQuery, qs)
				continue
			}
			// 自定义analyzer时，multi_match无法使用minimum_should_match，默认为至少一个满足，导致结果集还是很大
			// ngram(2,2)
			for _, kw := range v.KW {
				rn := []rune(kw)
				for i := 0; i+1 < len(rn); i++ {
					kwStr := string(rn[i : i+2])
					for _, kwField := range v.KWFields {
						likeQuery = append(likeQuery, elastic.NewTermQuery(kwField, kwStr))
					}
				}
			}
		case "", LikeLevelLow:
			qs := elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("90%") // 默认是and
			if v.Or {
				qs.Operator("OR")
			}
			likeQuery = append(likeQuery, qs)
		}
	}
	return
}
