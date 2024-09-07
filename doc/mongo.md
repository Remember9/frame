# MongoDB

0.配置

```yaml
mongo:
  uri: mongodb://esf:esf@lejuesf.com:27017/admin?readPreference=secondaryPreferred
  maxPoolSize: 10
  database: test # cli模式需要
  coll: golang_frame # cli模式需要
```

1.client模式使用

```go
func test () {
    ctx := context.Background()
    client := mongo.Build(ctx, "mongo")
    var doc = bson.M{"test": "4353454354"}
    defer func() {
        if err := client.Close(ctx); err != nil {
            panic(err)
        }
    }()

    collection := client.Database("test").Collection("golang_frame")
    result, err := collection.InsertOne(ctx, doc)
}
```

2.cli模式使用

```go
func test () {
    ctx := context.Background()
    client := mongo.BuildCli(ctx, "mongo")
    var doc = bson.M{"testcli": "121212121"}
    defer func() {
        if err := client.Close(ctx); err != nil {
            panic(err)
        }
    }()

    result, err := client.InsertOne(ctx, doc)
}
```

3.更多

前往 https://github.com/qiniu/qmgo/blob/master/README_ZH.md
