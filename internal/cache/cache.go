package cache

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbFileName string) *sql.DB {
	var DB *sql.DB
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "translators")
	os.MkdirAll(cacheDir, os.ModePerm)
	dbPath := filepath.Join(cacheDir, dbFileName)
	DB, _ = sql.Open("sqlite3", dbPath)
	// Create table if not exists
	CreateTable(DB)

	return DB
}

func CreateTable(con *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS words (
        input_word TEXT NOT NULL,
        response_word TEXT UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        response_url TEXT UNIQUE NOT NULL,
        response_text TEXT NOT NULL
    )`
	_, err := con.Exec(query)
	return err
}

func InsertIntoTable(con *sql.DB, inputWord, responseWord, url, text string) error {
	var currentDatetime = time.Now()
	query := `INSERT INTO words (input_word, response_word, created_at, response_url, response_text) VALUES (?, ?, ?, ?, ?)`
	_, err := con.Exec(query, inputWord, responseWord, currentDatetime, url, text)
	return err
}

func GetCache(con *sql.DB, word, requestURL string) (string, string, error) {
	query := `SELECT response_url, response_text, created_at FROM words WHERE response_url = ? OR response_word = ? OR input_word = ?`
	row := con.QueryRow(query, requestURL, word, word)

	var responseURL, responseText string
	var createdAt time.Time
	err := row.Scan(&responseURL, &responseText, &createdAt)
	if err == sql.ErrNoRows || createdAt.Before(time.Now().Add(-24*time.Hour)) {
		return "", "", nil
	}
	return responseURL, responseText, err
}

func GetResponseWords(con *sql.DB) ([][2]string, error) {
	query := `SELECT response_word, created_at FROM words`
	rows, err := con.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data [][2]string
	for rows.Next() {
		var responseWord string
		var createdAt time.Time
		err := rows.Scan(&responseWord, &createdAt)
		if err != nil {
			return nil, err
		}
		data = append(data, [2]string{responseWord, createdAt.String()})
	}
	return data, nil
}

func GetRandomWords(con *sql.DB) ([]string, error) {
	query := `SELECT response_word FROM words ORDER BY RANDOM() LIMIT 20`
	rows, err := con.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []string
	for rows.Next() {
		var responseWord string
		err := rows.Scan(&responseWord)
		if err != nil {
			return nil, err
		}
		data = append(data, responseWord)
	}
	return data, nil
}

func DeleteWord(con *sql.DB, word string) (bool, [2]string, error) {
	query := `SELECT input_word, response_url FROM words WHERE input_word = ? OR response_word = ?`
	row := con.QueryRow(query, word, word)

	var inputWord, responseURL string
	err := row.Scan(&inputWord, &responseURL)
	if err == sql.ErrNoRows {
		return false, [2]string{}, nil
	}
	if err != nil {
		return false, [2]string{}, err
	}

	query = `DELETE FROM words WHERE input_word = ? OR response_word = ?`
	_, err = con.Exec(query, word, word)
	if err != nil {
		return false, [2]string{}, err
	}

	return true, [2]string{inputWord, responseURL}, nil
}
