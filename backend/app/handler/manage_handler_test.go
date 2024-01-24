package handler

import (
	"context"
	//	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	//"time"

	//	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// type RedisClient interface {
// 	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
// 	Get(ctx context.Context, key string) *redis.StringCmd
// 	// Include other methods you're using...
// }

// MockRedisClient - мок для redis.Client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	args := m.Called(ctx, cursor, match, count)
	return args.Get(0).(*redis.ScanCmd)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.DurationCmd)
}

func NewHandlerMock(database RedisClient) Handler {
	return Handler{Database: database}
}

func TestAllKeys(t *testing.T) {
	mockRedis := new(MockRedisClient)
	handler := NewHandlerMock(mockRedis)

	// Мокуємо необхідні виклики redis.Client для тесту AllKeys

	// Приклад мокування Scan
	mockRedis.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&redis.ScanCmd{})

	// Приклад мокування Get
	// mockRedis.On("Get", mock.Anything, mock.Anything).
	// 	Return(&redis.StringCmd{
	// 		Cmd: &redis.Cmd{},
	// 		Val: "some_value",
	// 	})

	req, err := http.NewRequest("GET", "/keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.AllKeys(rr, req)

	// Перевірка статусу відповіді
	assert.Equal(t, http.StatusOK, rr.Code)

	// Додаткові перевірки, які вам можуть знадобитися...
}

// Аналогічно напишіть тести для інших методів, таких як GroupKeys, DeleteKey, DeleteByGroup, DeleteAllKeys, GetKey.
