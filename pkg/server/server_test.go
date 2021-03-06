package server_test

import (
	"fmt"
	"github.com/muratdemir0/faceit-task/pkg/server"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
	"time"
)

func TestGivenServerConfigWhenICallRunthenItShouldRunOnSpecifiedPort(t *testing.T) {
	freePort, err := freeport.GetFreePort()
	assert.Nil(t, err)
	port := fmt.Sprintf(":%d", freePort)

	s := server.New(port, []server.Handler{}, zap.NewExample())
	go s.Run()

	time.Sleep(50 * time.Millisecond)
	testEndpointURL := fmt.Sprintf("http://localhost%s/health", port)
	req, err := http.NewRequest(http.MethodGet, testEndpointURL, http.NoBody)
	assert.Nil(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGivenServerConfigWithInvalidPortWhenRunIsCalledThenServerShouldPanic(t *testing.T) {
	invalidPort := fmt.Sprintf(":%d", -1)

	s := server.New(invalidPort, []server.Handler{}, zap.NewExample())

	assert.Panics(t, func() { s.Run() })
}
