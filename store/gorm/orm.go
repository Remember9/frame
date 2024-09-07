package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// Open ...
func OpenMysql(dbConfig *Config) (*gorm.DB, error) {
	inner, err := gorm.Open(mysql.Open(dbConfig.DSN), &gorm.Config{
		DryRun:                 dbConfig.DryRun,
		SkipDefaultTransaction: dbConfig.SkipDefaultTransaction,
		Logger:                 NewLogger(UnmarshalText(dbConfig.LogLevel), dbConfig.SlowThreshold),
	})
	if err != nil {
		return nil, err
	}

	var dsnReplicas []gorm.Dialector
	if len(dbConfig.DSNReplicas) > 0 {
		for _, v := range dbConfig.DSNReplicas {
			dsnReplicas = append(dsnReplicas, mysql.Open(v))
		}
	}
	err = inner.Use(
		dbresolver.Register(dbresolver.Config{Replicas: dsnReplicas}).
			SetConnMaxLifetime(dbConfig.ConnMaxLifetime).
			SetMaxIdleConns(dbConfig.MaxIdleConns).
			SetMaxOpenConns(dbConfig.MaxOpenConns),
	)
	if err != nil {
		return nil, err
	}

	return inner, err
}
