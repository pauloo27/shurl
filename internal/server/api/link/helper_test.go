package link_test

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/valkey-io/valkey-go"
)

func mockValkey() valkey.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress:  []string{s.Addr()},
		DisableCache: true,
	})
	if err != nil {
		panic(err)
	}

	return client
}
