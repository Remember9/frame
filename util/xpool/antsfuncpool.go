package xpool

import (
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
)

var ansPoolWithFunc sync.Map

// NewAntsPoolWithFunc 创建ans业务方法固定的协程池
//poolName为空时不放入全局map, n为协程数量，f为业务方法只有接收参数
//池子使用：resPool.Invoke(ages)
//协程池最好在业务方法开始前就使用NewAntsPoolWithNum生成好，避免在业务中动态生成
func NewAntsPoolWithFunc(poolName string, n int, f func(interface{}), options ...ants.Option) (resPool *ants.PoolWithFunc, topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("%s panic in NewAntsPoolWithFunc err:%#v", poolName, err)
		}
	}()
	p, err := ants.NewPoolWithFunc(n, f, options...)
	if err != nil {
		return nil, err
	}
	if poolName != "" {
		if err := AddAntsFuncPool(poolName, p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// GetAntsFuncPool 根据名称获取已注册协程池
//code返回用于附加判断出错时是否因为断言失败，特定情况使用该值
func GetAntsFuncPool(name string) (topPool *ants.PoolWithFunc, code int, topErr error) {
	defer func() {
		if err := recover(); err != nil {
			code = 500
			topErr = fmt.Errorf("%s panic in GetAntsFuncPool err:%#v", name, err)
		}
	}()
	if p, ok := ansPoolWithFunc.Load(name); !ok {
		return nil, 501, errors.New("该名字未注册协程池:" + name)
	} else {
		pool, ok2 := p.(*ants.PoolWithFunc)
		if !ok2 {
			return nil, 502, errors.New("获取协程池" + name + "成功，断言异常")
		}
		return pool, 0, nil
	}
}

func AddAntsFuncPool(name string, p *ants.PoolWithFunc) (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("%s panic in AddAntsFuncPool err:%#v", name, err)
		}
	}()
	if _, ok := ansPoolWithFunc.Load(name); ok {
		return errors.New("该名字已注册协程池:" + name)
	}
	ansPoolWithFunc.Store(name, p)
	return nil
}
