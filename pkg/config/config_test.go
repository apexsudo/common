package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apexsudo/common/pkg/config"
)

type testConfig struct {
	Name string `json:"name" required:"true"`
	Port int    `json:"port"`
}

func TestLoadConfig_loadsFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"name":"test","port":8080}`), 0o600))
	t.Setenv("CONFIG_PATH", dir)

	cfg, err := config.LoadConfig[testConfig]("config.json")
	require.NoError(t, err)
	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, 8080, cfg.Port)
}

func TestLoadConfig_returnsWithDefaultVal(t *testing.T) {
	t.Parallel()
	_, err := config.LoadConfig[testConfig]("config.json")
	require.Error(t, err)
}

func TestLoadConfig_errorsOnMissingFile(t *testing.T) {
	t.Setenv("CONFIG_PATH", t.TempDir())

	_, err := config.LoadConfig[testConfig]("nonexistent.json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error loading config")
}

func TestLoadConfig_respectsConfigPathEnvVar(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir1, "config.json"), []byte(`{"name":"from-dir1"}`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir2, "config.json"), []byte(`{"name":"from-dir2"}`), 0o600))
	t.Setenv("CONFIG_PATH", dir2)

	cfg, err := config.LoadConfig[testConfig]("config.json")
	require.NoError(t, err)
	assert.Equal(t, "from-dir2", cfg.Name)
}
