package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiURL            = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	serverPort        = ":8080"
	dbTimeout         = 10 * time.Millisecond
	apiRequestTimeout = 200 * time.Millisecond
)

type Quote struct {
	Bid string `json:"bid"`
}

type Response struct {
	USDBRL Quote `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", quoteHandler)
	log.Println("Server is running on port", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), apiRequestTimeout)
	defer cancel()

	quote, err := fetchQuote(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)

	if err := saveQuoteToDB(quote.Bid); err != nil {
		log.Println("Failed to save quote to database:", err)
	}
}

func fetchQuote(ctx context.Context) (*Quote, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse.USDBRL, nil
}

func saveQuoteToDB(bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS quotes (id INTEGER PRIMARY KEY, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO quotes (bid) VALUES (?)`, bid); err != nil {
		return fmt.Errorf("failed to insert quote: %w", err)
	}

	return nil
}
