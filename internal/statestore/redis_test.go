package statestore

import (
	"testing"

	utilTesting "github.com/FairlySadPanda/open-match-but-kube-got-shreked/internal/util/testing"
	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/pkg/pb"
	"github.com/stretchr/testify/require"
)

func TestNewMutex(t *testing.T) {
	cfg, closer := createRedis(t, false, "")
	defer closer()
	service := New(cfg)
	require.NotNil(t, service)
	defer service.Close()
	ctx := utilTesting.NewContext(t)

	mutex := service.NewMutex("key")

	err := mutex.Lock(ctx)
	require.NoError(t, err)

	err = service.CreateBackfill(ctx, &pb.Backfill{
		Id: "222",
	}, nil)
	require.NoError(t, err)

	b, err := mutex.Unlock(ctx)
	require.NoError(t, err)
	require.True(t, b)

}
