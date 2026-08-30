package tests

import (
	"os"
	"sync"
	"testing"

	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/service"
)

var integrationOnce sync.Once

func requireMongoIntegration(t *testing.T) {
	t.Helper()
	url := os.Getenv("LEANOTE_INTEGRATION_MONGO_URL")
	if url == "" {
		t.Skip("set LEANOTE_INTEGRATION_MONGO_URL to run database integration tests")
	}
	integrationOnce.Do(func() {
		db.Init(url, "")
		service.InitService()
		service.ConfigS.InitGlobalConfigs()
	})
}
