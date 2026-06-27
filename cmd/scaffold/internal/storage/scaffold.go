package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type MigrationFile struct {
	UpFile   string
	DownFile string
}

const (
	createTable = `BEGIN;
CREATE TABLE %s
(
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
);
CREATE INDEX idx_%s__deleted_at ON %s(deleted_at);
COMMIT;
`
	dropTable = `DROP TABLE %s;
`
)

func ScaffoldTable(ctx context.Context, tableName string) error {
	tableMigrationFiles, err := createMigrationFile(ctx, fmt.Sprintf("create_%s_table", tableName))
	if err != nil {
		return err
	}
	err = writeToFile(tableMigrationFiles.UpFile, fmt.Sprintf(createTable, tableName, tableName, tableName))
	if err != nil {
		return fmt.Errorf("failed writing content to the create table (up) file: %w", err)
	}
	err = writeToFile(tableMigrationFiles.DownFile, fmt.Sprintf(dropTable, tableName))
	if err != nil {
		return fmt.Errorf("failed writing content to the create table (down) file: %w", err)
	}

	return nil
}

func ScaffoldMigration(ctx context.Context, migrationName string) error {
	_, err := createMigrationFile(ctx, migrationName)

	return err
}

const migrationScriptsRelPath = "internal/storage/database/migration/scripts"

func findMigrationScriptsDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	for {
		candidate := filepath.Join(dir, migrationScriptsRelPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find %q in any parent directory", migrationScriptsRelPath)
		}
		dir = parent
	}
}

func createMigrationFile(ctx context.Context, name string) (*MigrationFile, error) {
	if name == "" {
		return nil, errors.New("migration name cannot be empty")
	}
	scriptsDir, err := findMigrationScriptsDir()
	if err != nil {
		return nil, err
	}
	data, err := exec.CommandContext(
		ctx,
		"go", "tool", "migrate", "create",
		"-ext", "sql",
		"-dir", scriptsDir,
		name,
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create migration file: %w", err)
	}
	up, down := findMigration(string(data), "up"), findMigration(string(data), "down")
	if up == nil || down == nil {
		return nil, fmt.Errorf("failed to read migration file information")
	}

	return &MigrationFile{UpFile: *up, DownFile: *down}, nil
}

func findMigration(data string, migrationType string) *string {
	parts := strings.Split(data, "\n")
	for _, part := range parts {
		if strings.HasSuffix(part, fmt.Sprintf("%s.sql", migrationType)) {
			return &part
		}
	}

	return nil
}

func writeToFile(path string, content string) error {
	const defaultPermission = 0644
	err := os.WriteFile(path, []byte(content), defaultPermission)
	if err != nil {
		return fmt.Errorf("failed to write content to migration file: %w", err)
	}

	return nil
}
