package redis

import "github.com/go-redis/redis/v8"

//Redis client
type Redis struct {
	Config *Config
	Client redis.Cmdable
}

// Cluster try to get a redis.ClusterClient
func (r *Redis) Cluster() *redis.ClusterClient {
	if c, ok := r.Client.(*redis.ClusterClient); ok {
		return c
	}
	return nil
}

//Stub try to get a redis.Client
func (r *Redis) Stub() *redis.Client {
	if c, ok := r.Client.(*redis.Client); ok {
		return c
	}
	return nil
}
