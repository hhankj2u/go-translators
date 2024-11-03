package dicts_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hhankj2u/translators/pkg/cache"
	"github.com/hhankj2u/translators/pkg/dicts"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = cache.CreateTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Cleanup function to close the test database after tests
	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	client := &http.Client{}
	resp, err := dicts.Fetch(server.URL, client)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestCacheRun(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	inputWord := "test"
	reqURL := "http://example.com"
	responseText := "<html><body><p>Hello, World!</p></body></html>"

	err := cache.InsertIntoTable(db, inputWord, inputWord, reqURL, responseText)
	if err != nil {
		t.Fatalf("Failed to insert into table: %v", err)
	}

	cached, soup, err := dicts.CacheRun(db, inputWord, reqURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !cached {
		t.Fatalf("Expected cached to be true, got false")
	}
	if soup == nil {
		t.Fatalf("Expected soup to be non-nil")
	}
}

func TestSave(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	inputWord := "test"
	responseWord := "test_response"
	responseURL := "http://example.com"
	responseText := "<html><body><p>Hello, World!</p></body></html>"

	dicts.Save(db, inputWord, responseWord, responseURL, responseText)

	// Verify insertion
	query := `SELECT input_word, response_word, response_url, response_text FROM words WHERE input_word = ?`
	row := db.QueryRow(query, inputWord)
	var iw, rw, ru, rt string
	err := row.Scan(&iw, &rw, &ru, &rt)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted data: %v", err)
	}
	if iw != inputWord || rw != responseWord || ru != responseURL || rt != responseText {
		t.Fatalf("Inserted data does not match: got (%s, %s, %s, %s)", iw, rw, ru, rt)
	}
}
