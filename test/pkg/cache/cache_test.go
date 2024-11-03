package cache_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/hhankj2u/translators/pkg/cache"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "translators_test")
	os.MkdirAll(cacheDir, os.ModePerm)
	dbPath := filepath.Join(cacheDir, "database_test")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Cleanup function to remove the test database after tests
	cleanup := func() {
		db.Close()
		os.RemoveAll(cacheDir)
	}

	return db, cleanup
}

func TestCreateTable(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := cache.CreateTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Verify table creation
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='words';`
	row := db.QueryRow(query)
	var tableName string
	err = row.Scan(&tableName)
	if err != nil || tableName != "words" {
		t.Fatalf("Table 'words' was not created")
	}
}

func TestInsertIntoTable(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := cache.CreateTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	inputWord := "test"
	responseWord := "test_response"
	url := "http://example.com"
	text := "response text"

	err = cache.InsertIntoTable(db, inputWord, responseWord, url, text)
	if err != nil {
		t.Fatalf("Failed to insert into table: %v", err)
	}

	// Verify insertion
	query := `SELECT input_word, response_word, response_url, response_text FROM words WHERE input_word = ?`
	row := db.QueryRow(query, inputWord)
	var iw, rw, ru, rt string
	err = row.Scan(&iw, &rw, &ru, &rt)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted data: %v", err)
	}
	if iw != inputWord || rw != responseWord || ru != url || rt != text {
		t.Fatalf("Inserted data does not match: got (%s, %s, %s, %s)", iw, rw, ru, rt)
	}
}

func TestGetCache(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := cache.CreateTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	inputWord := "test"
	responseWord := "test_response"
	url := "http://example.com"
	text := "response text"

	err = cache.InsertIntoTable(db, inputWord, responseWord, url, text)
	if err != nil {
		t.Fatalf("Failed to insert into table: %v", err)
	}

	resURL, resText, err := cache.GetCache(db, inputWord, url)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}
	if resURL != url || resText != text {
		t.Fatalf("Cache data does not match: got (%s, %s)", resURL, resText)
	}
}

func TestDeleteWord(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := cache.CreateTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	inputWord := "test"
	responseWord := "test_response"
	url := "http://example.com"
	text := "response text"

	err = cache.InsertIntoTable(db, inputWord, responseWord, url, text)
	if err != nil {
		t.Fatalf("Failed to insert into table: %v", err)
	}

	deleted, data, err := cache.DeleteWord(db, inputWord)
	if err != nil {
		t.Fatalf("Failed to delete word: %v", err)
	}
	if !deleted {
		t.Fatalf("Word was not deleted")
	}
	if data[0] != inputWord || data[1] != url {
		t.Fatalf("Deleted data does not match: got (%s, %s)", data[0], data[1])
	}

	// Verify deletion
	resURL, resText, err := cache.GetCache(db, inputWord, url)
	if err != nil {
		t.Fatalf("Failed to get cache after deletion: %v", err)
	}
	if resURL != "" || resText != "" {
		t.Fatalf("Cache data still exists after deletion: got (%s, %s)", resURL, resText)
	}
}
