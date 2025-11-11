package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const (
	migrationsPath  = "migrations"
	migrationsTable = "schema_migrations"
	colorReset      = "\033[0m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[33m"
	colorRed        = "\033[31m"
	colorBlue       = "\033[34m"
)

type Migration struct {
	Version  string
	FileName string
	UpSQL    string
	DownSQL  string
}

func main() {
	cmd := flag.String("cmd", "", "Command to run: up, down, status")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("%sError connecting to database: %v%s\n", colorRed, err, colorReset)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("%sError pinging database: %v%s\n", colorRed, err, colorReset)
	}

	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("%sError creating migrations table: %v%s\n", colorRed, err, colorReset)
	}

	switch *cmd {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("%sError running migrations: %v%s\n", colorRed, err, colorReset)
		}
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("%sError rolling back migration: %v%s\n", colorRed, err, colorReset)
		}
	case "status":
		if err := showStatus(db); err != nil {
			log.Fatalf("%sError showing status: %v%s\n", colorRed, err, colorReset)
		}
	default:
		log.Fatalf("%sInvalid command. Use: up, down, or status%s\n", colorRed, colorReset)
	}
}

func ensureMigrationsTable(db *sql.DB) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`, migrationsTable)

	_, err := db.Exec(query)
	return err
}

func getMigrations() ([]Migration, error) {
	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading migrations directory: %w", err)
	}

	migrationsMap := make(map[string]*Migration)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".sql") {
			continue
		}

		parts := strings.Split(fileName, "_")
		if len(parts) < 2 {
			continue
		}
		version := parts[0]

		content, err := os.ReadFile(filepath.Join(migrationsPath, fileName))
		if err != nil {
			return nil, fmt.Errorf("error reading migration file %s: %w", fileName, err)
		}

		if migrationsMap[version] == nil {
			migrationsMap[version] = &Migration{
				Version:  version,
				FileName: fileName,
			}
		}

		if strings.HasSuffix(fileName, ".up.sql") {
			migrationsMap[version].UpSQL = string(content)
		} else if strings.HasSuffix(fileName, ".down.sql") {
			migrationsMap[version].DownSQL = string(content)
		}
	}

	var migrations []Migration
	for _, m := range migrationsMap {
		migrations = append(migrations, *m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT version FROM %s", migrationsTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

func migrateUp(db *sql.DB) error {
	migrations, err := getMigrations()
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	hasNewMigrations := false

	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		hasNewMigrations = true
		fmt.Printf("%s→ Running migration %s...%s\n", colorBlue, migration.Version, colorReset)

		conn, err := db.Conn(context.Background())
		if err != nil {
			return fmt.Errorf("error getting connection: %w", err)
		}
		defer conn.Close()

		if _, err := conn.ExecContext(context.Background(), migration.UpSQL); err != nil {
			return fmt.Errorf("error executing migration %s: %w", migration.Version, err)
		}

		if _, err := conn.ExecContext(
			context.Background(),
			fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", migrationsTable),
			migration.Version,
		); err != nil {
			return fmt.Errorf("error recording migration %s: %w", migration.Version, err)
		}

		fmt.Printf("%s✓ Migration %s applied successfully%s\n", colorGreen, migration.Version, colorReset)
	}

	if !hasNewMigrations {
		fmt.Printf("%s✓ No new migrations to apply%s\n", colorGreen, colorReset)
	} else {
		fmt.Printf("\n%s✓ All migrations completed successfully!%s\n", colorGreen, colorReset)
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	migrations, err := getMigrations()
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	var lastMigration *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if applied[migrations[i].Version] {
			lastMigration = &migrations[i]
			break
		}
	}

	if lastMigration == nil {
		fmt.Printf("%s✓ No migrations to rollback%s\n", colorYellow, colorReset)
		return nil
	}

	fmt.Printf("%s→ Rolling back migration %s...%s\n", colorBlue, lastMigration.Version, colorReset)

	conn, err := db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("error getting connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(context.Background(), lastMigration.DownSQL); err != nil {
		return fmt.Errorf("error executing down migration %s: %w", lastMigration.Version, err)
	}

	if _, err := conn.ExecContext(
		context.Background(),
		fmt.Sprintf("DELETE FROM %s WHERE version = $1", migrationsTable),
		lastMigration.Version,
	); err != nil {
		return fmt.Errorf("error removing migration record %s: %w", lastMigration.Version, err)
	}

	fmt.Printf("%s✓ Migration %s rolled back successfully%s\n", colorGreen, lastMigration.Version, colorReset)

	return nil
}

func showStatus(db *sql.DB) error {
	migrations, err := getMigrations()
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	fmt.Printf("\n%sMigration Status:%s\n", colorBlue, colorReset)
	fmt.Println(strings.Repeat("=", 60))

	if len(migrations) == 0 {
		fmt.Printf("%sNo migrations found%s\n", colorYellow, colorReset)
		return nil
	}

	appliedCount := 0
	for _, migration := range migrations {
		status := ""
		if applied[migration.Version] {
			status = fmt.Sprintf("%s✓ Applied%s", colorGreen, colorReset)
			appliedCount++
		} else {
			status = fmt.Sprintf("%s✗ Pending%s", colorYellow, colorReset)
		}

		// Extract migration name from filename
		nameParts := strings.Split(migration.FileName, "_")
		name := strings.TrimSuffix(strings.Join(nameParts[1:], "_"), ".up.sql")
		name = strings.TrimSuffix(name, ".down.sql")

		fmt.Printf("%-10s %-40s %s\n", migration.Version, name, status)
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n%sTotal: %d migrations | Applied: %d | Pending: %d%s\n",
		colorBlue, len(migrations), appliedCount, len(migrations)-appliedCount, colorReset)

	if appliedCount > 0 {
		lastApplied := ""
		for i := len(migrations) - 1; i >= 0; i-- {
			if applied[migrations[i].Version] {
				lastApplied = migrations[i].Version
				break
			}
		}
		fmt.Printf("%sCurrent version: %s%s\n\n", colorGreen, lastApplied, colorReset)
	} else {
		fmt.Printf("%sNo migrations applied yet%s\n\n", colorYellow, colorReset)
	}

	return nil
}
