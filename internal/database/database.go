package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func InitDB(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Database connection established")
	return &DB{db}, nil
}

func (db *DB) RunMigrations() error {
	schema, err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %w", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("error executing migrations: %w", err)
	}

	if err := db.ensureQuestionPosition(); err != nil {
		return fmt.Errorf("error ensuring question position column: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// ensureQuestionPosition adds the questions.position column on databases created
// before it existed (MySQL has no ADD COLUMN IF NOT EXISTS), then backfills any
// rows still at 0 so every question has a distinct order.
func (db *DB) ensureQuestionPosition() error {
	var exists int
	err := db.QueryRow(`SELECT COUNT(*) FROM information_schema.columns
	    WHERE table_schema = DATABASE() AND table_name = 'questions' AND column_name = 'position'`).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		if _, err := db.Exec(`ALTER TABLE questions ADD COLUMN position INT NOT NULL DEFAULT 0`); err != nil {
			return err
		}
	}

	// Backfill: if every question is still at position 0 (fresh column), number
	// them by id so the existing order is preserved.
	var maxPos int
	if err := db.QueryRow(`SELECT COALESCE(MAX(position), 0) FROM questions`).Scan(&maxPos); err != nil {
		return err
	}
	if maxPos == 0 {
		if _, err := db.Exec(`SET @r := 0`); err != nil {
			return err
		}
		if _, err := db.Exec(`UPDATE questions SET position = (@r := @r + 1) ORDER BY id`); err != nil {
			return err
		}
	}
	return nil
}

// RunSeed loads internal/database/seed.sql. Statements are run one-by-one so
// drivers without multi-statement support still work. It will skip silently if
// the seed has already been applied (detected by counting rows in questions).
func (db *DB) RunSeed() error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM questions WHERE question_type IN ('acomodare','educatie_experienta','abilitati_personale')").Scan(&count); err != nil {
		return fmt.Errorf("error checking seed state: %w", err)
	}
	if count > 0 {
		log.Println("Seed already applied, skipping")
		return nil
	}

	raw, err := os.ReadFile("internal/database/seed.sql")
	if err != nil {
		return fmt.Errorf("error reading seed file: %w", err)
	}

	for _, stmt := range splitSQLStatements(string(raw)) {
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("error executing seed statement: %w\n--- statement ---\n%s", err, stmt)
		}
	}

	log.Println("Database seed applied successfully")
	return nil
}

// splitSQLStatements strips line comments and splits on ';'. Good enough for
// the curated seed file in this repo (no semicolons inside string literals).
func splitSQLStatements(src string) []string {
	var clean []string
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}
		clean = append(clean, line)
	}
	joined := strings.Join(clean, "\n")

	parts := strings.Split(joined, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func (db *DB) Close() error {
	return db.DB.Close()
}
