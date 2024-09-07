# Redis

0.配置

```yaml
# version >= 4.0
# 标注模式
redis:
  addrs:
    - mem.lejuesf.com:7516
  mode: stub
  debug: true
# 集群模式
redis-cluster:
  addrs:
    - 127.0.0.1:6379
    - 127.0.0.1:6380
  mode: cluster
```

1.使用

```go
package doc
import (
    "context"
    "github.com/Remember9/frame/client/redis"
)
var ctx = context.Background()
func (eng *Engine) exampleForRedisStub() (err error) {
	//build redisStub
	redisStub := redis.Build("redis")

    if err := redisStub.Stub().Ping(ctx).Err(); err != nil {
        xlog.Error("start redis", xlog.Any("err", err))
    }
	
	// set string
	setRes := redisStub.Set(ctx, "frame-redis", "redisStub", time.Second*5)
	xlog.Info("redisStub set string", xlog.Any("res", setRes))
	// get string
	getRes := redisStub.Get(ctx, "frame-redis")
	xlog.Info("redisStub get string", xlog.Any("res", getRes))
	return
}

func (eng *Engine) exampleForRedisClusterStub() (err error) {
	//build redisClusterStub
	redisClusterStub := redis.Build("redis-cluster")

    if err := redisClusterStub.Cluster().Ping(ctx).Err(); err != nil {
        xlog.Error("start redis", xlog.Any("err", err))
    }
	
	// set string
	setRes := redisClusterStub.Set(ctx, "frame-redisCluster", "redisClusterStub", time.Second*5)
	xlog.Info("redisClusterStub set string", xlog.Any("res", setRes))
	// get string
	getRes := redisClusterStub.Get(ctx, "frame-redisCluster")
	xlog.Info("redisClusterStub get string", xlog.Any("res", getRes))
	return
}

```
