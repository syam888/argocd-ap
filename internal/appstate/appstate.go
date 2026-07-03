package appstate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Dir            string
	CurrentAppFile string
}

func New() *Config {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".config", "ap")
	return &Config{
		Dir:            dir,
		CurrentAppFile: filepath.Join(dir, "current_app"),
	}
}

func (c *Config) EnsureDir() error {
	return os.MkdirAll(c.Dir, 0o755)
}

func (c *Config) Set(app string) error {
	if err := c.EnsureDir(); err != nil {
		return err
	}
	return os.WriteFile(c.CurrentAppFile, []byte(strings.TrimSpace(app)+"\n"), 0o644)
}

func (c *Config) Get() (string, error) {
	data, err := os.ReadFile(c.CurrentAppFile)
	if err != nil {
		return "", fmt.Errorf("no app selected. Run: ap select")
	}
	app := strings.TrimSpace(string(data))
	if app == "" {
		return "", fmt.Errorf("no app selected. Run: ap select")
	}
	return app, nil
}

func (c *Config) Clear() error {
	if err := os.Remove(c.CurrentAppFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
