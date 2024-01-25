package handler

import (
	"encoding/json"
	//"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"micro-manager-redis/app/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) GetAllKeys(pattern string) ([]repository.Keys, error) {
	args := m.Called(pattern)
	return args.Get(0).([]repository.Keys), args.Error(1)
}

func (m *MockRedisRepository) GroupKeys(pattern, separator string) ([]repository.SplitKeys, error) {
	args := m.Called(pattern, separator)
	return args.Get(0).([]repository.SplitKeys), args.Error(1)
}

func (m *MockRedisRepository) GetKey(key string) (repository.Keys, error) {
	args := m.Called(key)
	return args.Get(0).(repository.Keys), args.Error(1)
}

func (m *MockRedisRepository) DeleteKey(key string) {
	m.Called(key)
}

func (m *MockRedisRepository) DeleteAllKeys() {
	m.Called()
}

func (m *MockRedisRepository) DeleteByGroup(pattern string) error {
	args := m.Called(pattern)
	return args.Error(0)
}

// Implement other methods from RedisRepositoryInterface...

func Test_AllKeys(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("GetAllKeys", mock.Anything).
		Return([]repository.Keys{{Key: "key1", Value: "value1", Expire: 0}}, nil)

	req, err := http.NewRequest("GET", "/keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.AllKeys(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body or any other expectations
	var result []repository.Keys
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "key1", result[0].Key)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func Test_AllKeys_Range(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("GetAllKeys", mock.Anything).
		Return([]repository.Keys{
			{Key: "key1", Value: "value1", Expire: 0},
			{Key: "key2", Value: "value2", Expire: 0},
			{Key: "key3", Value: "value3", Expire: 0},
		}, nil)

	req, err := http.NewRequest("GET", "/keys?range=[0,1]", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.AllKeys(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body or any other expectations
	var result []repository.Keys
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "key1", result[0].Key)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func Test_AllKeys_Filter(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("GetAllKeys", mock.Anything).
		Return([]repository.Keys{
			{Key: "key1", Value: "value1", Expire: 0},
		}, nil)

	req, err := http.NewRequest("GET", `/keys?filter={"key":"key1"}`, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.AllKeys(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body or any other expectations
	var result []repository.Keys
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "key1", result[0].Key)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func TestGroupKeys(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("GroupKeys", mock.Anything, mock.Anything).
		Return([]repository.SplitKeys{{Key: "key1", Separator: ":"}}, nil)

	req, err := http.NewRequest("GET", "/keys/group", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GroupKeys(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body or any other expectations
	var result []repository.SplitKeys
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "key1", result[0].Key)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func Test_GroupKeys_Range(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("GroupKeys", mock.Anything, mock.Anything).
		Return([]repository.SplitKeys{
			{Key: "key1", Separator: ":"},
			{Key: "key2", Separator: ":"},
			{Key: "key3", Separator: ":"},
		}, nil)

	req, err := http.NewRequest("GET", "/keys/group?range=[1,2]", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GroupKeys(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body or any other expectations
	var result []repository.SplitKeys
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "key2", result[0].Key)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func TestGetKey(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("GetKey", mock.Anything).
		Return(repository.Keys{Key: "key1", Value: "value1", Expire: 0}, nil)

	req, err := http.NewRequest("GET", "/keys/key1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetKey(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body or any other expectations
	var result repository.Keys
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "key1", result.Key)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func TestDeleteKey(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("DeleteKey", mock.Anything).
		Return()

	req, err := http.NewRequest("DELETE", "/keys/key1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.DeleteKey(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func TestDeleteAllKeys(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("DeleteAllKeys").
		Return()

	req, err := http.NewRequest("DELETE", "/keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.DeleteAllKeys(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}

func TestDeleteByGroup(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	handler := NewHandler(mockRepo)

	// Mock repository method
	mockRepo.On("DeleteByGroup", mock.Anything).
		Return(nil)

	req, err := http.NewRequest("DELETE", "/keys/group", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.DeleteByGroup(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the mocked repository method was called with the expected argument
	mockRepo.AssertExpectations(t)
}
