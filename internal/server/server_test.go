package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bobgromozeka/metrics/internal/server/storage"
)

const Key = "key"

func TestUpdateJSON_BadRequest(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", nil)
	httpW := httptest.NewRecorder()

	stor := storage.NewMemory()
	server := New(
		stor, StartupConfig{
			HashKey: Key,
		},
	)

	server.ServeHTTP(httpW, req)

	response, _ := io.ReadAll(httpW.Result().Body)
	httpW.Result().Body.Close()

	assert.Equal(t, "Bad request: EOF\n", string(response))
	assert.Equal(t, http.StatusBadRequest, httpW.Code)
}

func TestUpdateJSON_WrongMetricsType(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id": "id", "type":"random"}`))
	httpW := httptest.NewRecorder()

	stor := storage.NewMemory()
	server := New(
		stor, StartupConfig{
			HashKey: Key,
		},
	)

	server.ServeHTTP(httpW, req)

	response, _ := io.ReadAll(httpW.Result().Body)
	httpW.Result().Body.Close()

	assert.Equal(t, "Wrong metrics type\n", string(response))
	assert.Equal(t, http.StatusBadRequest, httpW.Code)
}

func TestUpdateJSON_CounterType(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id": "id","type":"counter","delta":22}`))
	httpW := httptest.NewRecorder()

	stor := storage.NewMemory()
	server := New(
		stor, StartupConfig{
			HashKey: Key,
		},
	)

	stor.AddCounter(context.Background(), "id", 20)

	server.ServeHTTP(httpW, req)

	response, _ := io.ReadAll(httpW.Result().Body)
	httpW.Result().Body.Close()

	assert.Equal(t, "application/json", httpW.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":"id","type":"counter","delta":42}`+"\n", string(response))
	assert.Equal(t, http.StatusOK, httpW.Code)
}

func TestUpdateJSON_GaugeType(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id": "id","type":"gauge","value":33}`))
	httpW := httptest.NewRecorder()

	stor := storage.NewMemory()
	server := New(
		stor, StartupConfig{
			HashKey: Key,
		},
	)

	stor.SetGauge(context.Background(), "id", 123.123)

	server.ServeHTTP(httpW, req)

	response, _ := io.ReadAll(httpW.Result().Body)
	httpW.Result().Body.Close()

	assert.Equal(t, "application/json", httpW.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":"id","type":"gauge","value":33}`+"\n", string(response))
	assert.Equal(t, http.StatusOK, httpW.Code)
}
