package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	serverURL            = "http://localhost:8080/cotacao"
	clientRequestTimeout = 300 * time.Millisecond
)

type Quote struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientRequestTimeout)
	defer cancel()

	quote, err := fetchQuote(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch quote: %v", err)
	}

	if err := saveQuoteToFile(quote.Bid); err != nil {
		log.Fatalf("Failed to save quote to file: %v", err)
	}
}

func fetchQuote(ctx context.Context) (*Quote, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
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

	var quote Quote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &quote, nil
}

func saveQuoteToFile(bid string) error {
	data := fmt.Sprintf("Dólar: %s", bid)
	if err := ioutil.WriteFile("cotacao.txt", []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
