package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPostgresInstallationSeed(t *testing.T) {
	seedPath := filepath.Join("..", "..", "database", "seed.sql")
	data, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatal(err)
	}
	seed := string(data)

	for _, required := range []string{
		"INSERT INTO public.users",
		"'admin@pearlnote.com'",
		"'demo@pearlnote.com'",
		"INSERT INTO public.configs",
		"INSERT INTO public.notebooks",
		"INSERT INTO public.notes",
		"ON CONFLICT DO NOTHING",
	} {
		if !strings.Contains(seed, required) {
			t.Errorf("installation seed is missing %q", required)
		}
	}

	for _, excludedTable := range []string{"sessions", "tokens", "email_logs", "suggestions", "reports"} {
		if strings.Contains(seed, "INSERT INTO public."+excludedTable) {
			t.Errorf("installation seed contains runtime-only table %s", excludedTable)
		}
	}
}
