package api

import (
	"os"
	"testing"

	"github.com/cihanalici/api/db/sqlc"
	"github.com/cihanalici/api/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store sqlc.Store) *Server {
	config := util.Config{}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
