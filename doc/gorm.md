# GORM

0.配置

```yaml
mysql:
  logLevel: info # silent error warn info
  dsn: leju:leju@tcp(10.208.0.102:7306)/base_esf_leju_com?charset=utf8mb4&parseTime=True&loc=Local&readTimeout=1s&timeout=1s&writeTimeout=3s
  dsnReplicas:
    - leju:leju@tcp(10.208.0.102:7306)/base_esf_leju_com?charset=utf8mb4&parseTime=True&loc=Local&readTimeout=1s&timeout=1s&writeTimeout=3s
    - leju:leju@tcp(10.208.0.102:7306)/base_esf_leju_com?charset=utf8mb4&parseTime=True&loc=Local&readTimeout=1s&timeout=1s&writeTimeout=3s
  connMaxLifeTime: 30s
  maxIdleConns: 50
  maxOpenConns: 100
```

1.使用

```go
db := gorm.Build("mysql")
ctx := context.Background()
err = db.WithContext(ctx).Where("`citycode` = ?", citycode).First(&res).Error
```

2.更多

前往 https://gorm.io/zh_CN/docs/ 
