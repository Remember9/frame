// Package redis /**外部定义协程池方式队列
package redis

import (
	"errors"
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/xlog"
	"github.com/go-redis/redis/v8"
	"net"
)

type redisMq struct {
	stream     string //stream名
	Rdb        redis.Cmdable
	originName string //未加前缀队列名
}

// StreamName /**获取队列名称（需要加入根据预发和正式相关前缀处理）
func StreamName(name string) string {
	pre := "dev"
	if r := config.Get("app.env"); r == nil {
		panic("get config app.env nil")
	} else {
		if r.(string) != "" {
			pre = r.(string)
		}
	}
	if pre == "local" {
		prefix := ""
		if ip4, err := privateIPv4(); err != nil {
			xlog.Info("获取本地ip异常:" + err.Error())
		} else {
			prefix = ip4.String()
		}
		return "local_" + name + "_" + prefix
	}
	if pre == "prod" {
		return name
	}
	if pre == "pre" {
		return "pre_" + name
	}
	if pre == "lpt" {
		return "lpt_" + name
	}
	return "test_" + name
}
func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}
func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}
