package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockchainInfo(t *testing.T) {
	setup(t)

	tMockState.CommitTestBlocks(10)

	w := httptest.NewRecorder()
	r := new(http.Request)

	tHTTPServer.BlockchainHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "10")
}

func TestNetworkInfo(t *testing.T) {
	setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)

	tHTTPServer.NetworkHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "Peers")
	assert.Contains(t, w.Body.String(), "ID")
	//	fmt.Println(w.Body.String())
}
