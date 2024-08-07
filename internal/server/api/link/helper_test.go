package link_test

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func mockRedis() *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}
