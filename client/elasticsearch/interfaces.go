package elasticsearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

// NewGenericClient create a ES Client
func NewGenericClient(connectConfig *Config) (GenericClient, error) {
	if connectConfig.Version == "" {
		connectConfig.Version = "v7"
	}
	switch connectConfig.Version {
	case "v7":
		return NewV7Client(connectConfig)
	default:
		return nil, fmt.Errorf("not supported ElasticSearch version: %v", connectConfig.Version)
	}
}

type (
	// GenericClient is a generic interface for all versions of ElasticSearch clients
	GenericClient interface {
		CreateIndex(ctx context.Context, index, mapping string) error
		DeleteIndex(ctx context.Context, index string) error
		PutMapping(ctx context.Context, index, root, key, valueType string) error
		Index(ctx context.Context, index, id, bodyJson string) (string, error)
		Delete(ctx context.Context, index, id string) (string, error)
		Update(ctx context.Context, index, id string, doc map[string]interface{}, p *UpdateParamsV7) (string, error)
		BulkIndex(ctx context.Context, bulkData []BulkIndexItemV7) (*elastic.BulkResponse, error)
		BulkDelete(ctx context.Context, bulkData []BulkDeleteItemV7) (*elastic.BulkResponse, error)
		BulkUpsert(ctx context.Context, upsertData *UpsertParamsV7) (*elastic.BulkResponse, error)
		CountByQuery(ctx context.Context, index, query string) (int64, error)
		Search(ctx context.Context, p *SearchParametersV7) (*elastic.SearchResult, error)
		SearchWithDSL(ctx context.Context, index, query string) (*elastic.SearchResult, error)
		QueryBody(sp *QueryBody) *elastic.BoolQuery
		Client() *elastic.Client
	}
)
