# CRON

0.配置
```yaml
  cron:
    withSeconds: false
    concurrentDelay: -1
    immediatelyRun: false
    distributedTask: true
```

1.使用
```go
_, err := c.AddJob("10 * * * * *", xcron.FuncJob(s.biz.Group.Job.Test), time.Hour)
if err != nil {
xlog.Error("cron err", xlog.FieldErr(err))
}

_, err := c.AddFunc("10 * * * * *", s.biz.Group.Job.Test, time.Hour)
if err != nil {
    xlog.Error("cron err", xlog.FieldErr(err))
    return
}

c.Schedule(xcron.Every(time.Second*5), xcron.FuncJob(s.biz.Group.Job.Test), time.Second*10)
```
