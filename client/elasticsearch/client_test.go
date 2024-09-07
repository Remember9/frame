package elasticsearch

import (
	"context"
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/util/xcast"
	"github.com/Remember9/frame/util/xstring"
	"github.com/Remember9/frame/xlog"
	"testing"
)

func init() {
	err := config.InitTest()
	if err != nil {
		panic(err)
	}
	err = xlog.Build()
	if err != nil {
		panic(err)
	}
}

func TestElasticV7_Update(t *testing.T) {
	client, err := Build("elasticsearch")
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	rep, err := client.Update(ctx, "testindex", "5", map[string]interface{}{"name": "test_55"}, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(rep)
}

func TestElasticV7_BulkIndex(t *testing.T) {
	client, err := Build("elasticsearch")
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	var bulkData []BulkIndexItemV7
	for i := 1; i < 10; i++ {
		id := xcast.ToString(i)
		bulkData = append(bulkData, BulkIndexItemV7{
			IndexName: "testindex",
			IndexID:   id,
			OpType:    "index",
			Doc: MapDataV7{
				"name": "test_" + id,
			},
		})
	}
	rep, err := client.BulkIndex(ctx, bulkData)
	if err != nil {
		t.Error(err)
	}
	t.Log(xstring.Json(rep))
}

func TestElasticV7_BulkDelete(t *testing.T) {
	client, err := Build("elasticsearch")
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	var bulkData []BulkDeleteItemV7
	for i := 1; i < 3; i++ {
		id := xcast.ToString(i)
		bulkData = append(bulkData, BulkDeleteItemV7{
			IndexName: "testindex",
			IndexID:   id,
		})
	}
	rep, err := client.BulkDelete(ctx, bulkData)
	if err != nil {
		t.Error(err)
	}
	t.Log(xstring.Json(rep))
}

func TestElasticV7_BulkUpsert(t *testing.T) {
	client, err := Build("elasticsearch")
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	var bulkData []BulkUpsertItemV7
	for i := 1; i < 5; i++ {
		id := xcast.ToString(i)
		bulkData = append(bulkData, BulkUpsertItemV7{
			IndexName: "testindex",
			IndexID:   id,
			Doc: MapDataV7{
				"name": "u_test_" + id,
			},
		})
	}
	rep, err := client.BulkUpsert(ctx, &UpsertParamsV7{
		RetryOnConflict: 15,
		Insert:          false,
		UpsertBody:      bulkData,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(rep.Errors)
}

func TestSearch(t *testing.T) {
	client, err := Build("elasticsearch")
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	// 实现:
	// (aid=122 or id=677) && (tid in (1,2,3,21) or tid_type in (1,2,3)) && (id > 10) &&
	// (aid=88 or fid=99) && (mid in (11,33) or id in (22,33)) && (2 < cid <= 10  || sid > 10)
	cmbA := &QueryBodyWhereCombo{}
	cmbA.ComboEQ([]map[string]interface{}{
		{"aid": 122},
		{"id": 677},
	}).ComboIn([]map[string][]interface{}{
		{"tid": {1, 2, 3, 21}},
		{"tid_type": {1, 2, 3}},
	}).ComboRange([]map[string]string{
		{"id": "(10,)"},
	}).MinIn(1).MinEQ(1).MinRange(1).MinAll(1)

	cmbB := &QueryBodyWhereCombo{}
	cmbB.ComboEQ([]map[string]interface{}{
		{"aid": 88},
		{"fid": 99},
	}).ComboIn([]map[string][]interface{}{
		{"mid": {11, 33}},
		{"id": {22, 33}},
	}).ComboRange([]map[string]string{
		{"cid": "(2,4]"},
		{"sid": "(10,)"},
	}).MinEQ(1).MinIn(2).MinRange(1).MinAll(1)
	params := NewRequest().Fields("aid", "id").WhereEq("txt", "hh").WhereCombo(cmbA, cmbB).Params()
	result, err := client.Search(ctx, &SearchParametersV7{
		Index:    "testindex",
		Query:    client.QueryBody(params),
		From:     0,
		PageSize: 2,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(xstring.Json(result))
}

func BenchmarkParallelEs(b *testing.B) {
	client, err := Build("elasticsearch")
	if err != nil {
		b.Error(err)
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := client.SearchWithDSL(context.Background(), "testindex", "{}")
			if err != nil {
				b.Error(err)
			}
			b.Log(result)
		}
	})
}
