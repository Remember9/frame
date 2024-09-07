package mongox

import (
	"context"
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/xlog"
	"go.mongodb.org/mongo-driver/bson"
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

func TestBuild(t *testing.T) {
	ctx := context.Background()
	client := Build(ctx, "mongo")
	var doc = bson.M{"test": "4353454354"}
	defer func() {
		if err := client.Close(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("test").Collection("golang_frame")
	result, err := collection.InsertOne(ctx, doc)

	if err != nil {
		t.Error(err)
	}
	t.Logf("insertID:%v", result.InsertedID)
}

func TestBuildCli(t *testing.T) {
	ctx := context.Background()
	client := BuildCli(ctx, "mongo")
	var doc = bson.M{"testcli": "121212121"}
	defer func() {
		if err := client.Close(ctx); err != nil {
			panic(err)
		}
	}()

	result, err := client.InsertOne(ctx, doc)

	if err != nil {
		t.Error(err)
	}
	t.Logf("insertID:%v", result.InsertedID)
}
