package valkey

import (
	"context"
	"log/slog"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/valkey-io/valkey-go"
)

func New(cfg *config.Valkey) (valkey.Client, error) {
	vkey, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{cfg.Address},
		Password:    cfg.Password,
		SelectDB:    cfg.DB,
	})
	if err != nil {
		return nil, err
	}

	cmd := vkey.B().Ping().Build()
	if res := vkey.Do(context.Background(), cmd); res.Error() != nil {
		return nil, err
	}

	slog.Info("Connected to Valkey", "addr", cfg.Address)

	return vkey, nil
}
