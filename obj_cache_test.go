package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

func TestIsExpired(t *testing.T) {
	domain := DomainModel("xxx")

	cache := NewCache(context.Background(), domain)
	defer cache.Close()

	now := time.Now()
	require.False(t, cache.isTimeExpired(now.UnixNano(), 1, now))
	require.True(t, cache.isTimeExpired(now.UnixNano(), 1, now.Add(time.Second*time.Duration(2))))
}

func TestInMemoryCache(t *testing.T) {
	require := require.New(t)

	domain := DomainModel("xxx")

	cache := NewCache(context.Background(), domain)
	defer cache.Close()

	dto := NewDTO(domain, []byte("1"))
	require.NoError(cache.Set(dto))

	fetchedDTO, errGet1 := cache.Get(dto.Key)
	t.Log(*fetchedDTO)
	t.Log(dto.Key)

	require.NoError(errGet1)
	require.Nil(deep.Equal(dto, fetchedDTO), fmt.Sprintf("dto: %#v\n", *dto))

	errDel := cache.Delete(dto.Key)
	require.NoError(errDel)

	shouldBeNilDTO, errGet2 := cache.Get(dto.Key)
	require.Error(errGet2)
	require.Nil(shouldBeNilDTO)
}
