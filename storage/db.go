package storage

import (
	"database/sql"
	"fmt"
	"os"

	"be-links/models"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}
	
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS links (
		id VARCHAR(255) PRIMARY KEY,
		universal_link TEXT NOT NULL,
		ios_store TEXT NOT NULL,
		android_store TEXT NOT NULL,
		title TEXT,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		click_count INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_links_id ON links(id);

	-- Migration: Add new columns if they don't exist (for existing databases)
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='links' AND column_name='universal_link') THEN
			ALTER TABLE links ADD COLUMN universal_link TEXT;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='links' AND column_name='ios_store') THEN
			ALTER TABLE links ADD COLUMN ios_store TEXT;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='links' AND column_name='android_store') THEN
			ALTER TABLE links ADD COLUMN android_store TEXT;
		END IF;
	END $$;
	`

	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) CreateLink(link *models.Link) error {
	query := `
		INSERT INTO links (id, universal_link, ios_store, android_store, title, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := db.conn.Exec(query, link.ID, link.UniversalLink, link.IOSStore, link.AndroidStore,
		link.Title, link.Description, link.CreatedAt, link.UpdatedAt)
	return err
}

func (db *DB) GetLink(id string) (*models.Link, error) {
	query := `
		SELECT id, universal_link, ios_store, android_store, title, description, created_at, updated_at, click_count
		FROM links WHERE id = $1
	`
	
	link := &models.Link{}
	err := db.conn.QueryRow(query, id).Scan(
		&link.ID, &link.UniversalLink, &link.IOSStore, &link.AndroidStore,
		&link.Title, &link.Description, &link.CreatedAt, &link.UpdatedAt, &link.ClickCount,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return link, err
}

func (db *DB) GetAllLinks() ([]*models.Link, error) {
	query := `
		SELECT id, universal_link, ios_store, android_store, title, description, created_at, updated_at, click_count
		FROM links ORDER BY created_at DESC
	`
	
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var links []*models.Link
	for rows.Next() {
		link := &models.Link{}
		err := rows.Scan(
			&link.ID, &link.UniversalLink, &link.IOSStore, &link.AndroidStore,
			&link.Title, &link.Description, &link.CreatedAt, &link.UpdatedAt, &link.ClickCount,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	
	return links, rows.Err()
}

func (db *DB) IncrementClickCount(id string) error {
	query := `UPDATE links SET click_count = click_count + 1 WHERE id = $1`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}