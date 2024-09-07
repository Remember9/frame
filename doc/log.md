# 日志

0.配置

```yaml
log:
  prefix: frame-srv
  level: debug # debug info warn error panic
  pathDir: ./logs/
  development: true # false 的时候会写入文件
```

1.使用

```go
xlog.Info("init", xlog.String("mod", "app"), xlog.String("event", "init"))
xlog.Warn("xx warn", xlog.String("key", "value"))
xlog.Error("xx error", xlog.FieldErr(err))    
xlog.Panic(err.Error())

xlog.Infof(str, args...)
```
