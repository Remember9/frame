使用如下配置，对 resty 进行配置：

```yaml
resty:
  debug: true   # 是否开启 debug 模式
  retryCount: 2 # 重试次数
  timeout: 0s   # 超时时间。为 0 时，表示不设置超时时间
```