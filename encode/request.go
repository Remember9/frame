package encode

import (
	"context"
	"net/http"
)

func DecodeNullRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}
