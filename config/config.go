package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RDB struct {
	Dbfilename string
	Dir        string
}

type Config struct {
	RDB RDB
}

func NewConfig(dbfilename, dir string) (*Config, error) {
	if err := validateRDBFilename(dbfilename); err != nil {
		return nil, err
	}

	if err := validateRDBDir((dir)); err != nil {
		return nil, err
	}

	return &Config{
		RDB: RDB{
			Dbfilename: dbfilename,
			Dir:        dir,
		},
	}, nil
}

func (c *Config) GetConfig() map[string]string {
	return map[string]string{
		"dbfilename": c.RDB.Dbfilename,
		"dir":        c.RDB.Dir,
	}
}

func (c *Config) GetRDBPath() string {
	return filepath.Join(c.RDB.Dir, c.RDB.Dbfilename)
}

func validateRDBFilename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("dbfilename cannot be an empty string")
	}

	if strings.Contains(name, "/") || strings.Contains(name, `\`) {
		return fmt.Errorf("dbfilename must not contain path separators: %q", name)
	}

	// dump.rdb not ../dumb.rdb
	if name != filepath.Base(name) {
		return fmt.Errorf("dbfilename must be a simple filename, got: %q", name)
	}

	return nil
}

func validateRDBDir(path string) error {
	// 1. must be an abs path
	if !filepath.IsAbs(path) {
		return fmt.Errorf("dir path must be absolute")
	}

	// 2. must exist
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory doesn't exist")
		}

		return fmt.Errorf("error checking path %v", err)
	}

	// 3. must be a dir
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// 4. must be writable
	testFile := filepath.Join(path, ".redis_write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("directory is not writable")
	}

	f.Close()
	os.Remove(testFile)

	return nil
}
