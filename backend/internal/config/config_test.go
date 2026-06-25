package config

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadWithEnvironmentOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	key := base64.StdEncoding.EncodeToString([]byte(strings.Repeat("k", 32)))
	content := `
app: {env: development, base_url: http://localhost, timezone: Asia/Jakarta, log_level: info}
http: {addr: ":8080", allowed_origins: [http://localhost:3000]}
grpc: {addr: ":9090"}
database: {url: postgres://localhost/crm, max_conns: 10, min_conns: 1}
redis: {url: redis://localhost:6379/0}
auth: {jwt_secret: "01234567890123456789012345678901", jwt_ttl: 24h, bcrypt_cost: 12}
security: {encryption_key: "` + key + `"}
telegram: {worker_interval: 5s, worker_batch_size: 10}
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HTTP_ADDR", ":8181")
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTP.Addr != ":8181" {
		t.Fatalf("expected env override, got %s", cfg.HTTP.Addr)
	}
	if cfg.Auth.JWTTTL.Hours() != 24 {
		t.Fatalf("unexpected JWT TTL %s", cfg.Auth.JWTTTL)
	}
}
