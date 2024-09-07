package xpool

import (
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
)

var ansPool sync.Map

// NewAntsPool 创建ans协程池
//poolName为空时不放入全局map,当将池子放入全局池map时若map已存在同名池子将返回错误
//n为协程数量
//使用池子用：resPool.Submit(func() {})
//协程池最好在业务方法开始前就使用NewAntsPoolWithNum生成好，避免在业务中动态生成
func NewAntsPool(poolName string, n int, options ...ants.Option) (resPool *ants.Pool, topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("%s panic in NewAntsPoolWithNum err:%#v", poolName, err)
		}
	}()
	p, err := ants.NewPool(n, options...)
	if err != nil {
		return nil, err
	}
	if poolName != "" {
		if err := AddAntsPool(poolName, p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// GetAntsPool 根据名称获取已注册协程池
//code返回用于附加判断出错时是否因为断言失败，特定情况使用该值
func GetAntsPool(name string) (topPool *ants.Pool, code int, topErr error) {
	defer func() {
		if err := recover(); err != nil {
			code = 500
			topErr = fmt.Errorf("%s panic in GetAntsPool err:%#v", name, err)
		}
	}()
	if p, ok := ansPool.Load(name); !ok {
		return nil, 501, errors.New("该名字未注册协程池:" + name)
	} else {
		pool, ok2 := p.(*ants.Pool)
		if !ok2 {
			return nil, 502, errors.New("获取协程池" + name + "成功，断言异常")
		}
		return pool, 0, nil
	}
}

func AddAntsPool(name string, p *ants.Pool) (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("%s panic in AddAntsPool err:%#v", name, err)
		}
	}()
	if _, ok := ansPool.Load(name); ok {
		return errors.New("该名字已注册协程池:" + name)
	}
	ansPool.Store(name, p)
	return nil
}
