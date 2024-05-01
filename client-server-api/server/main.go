package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Coin struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

var (
	QUOTE_API = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

func main() {
	http.HandleFunc("/cotacao", handleQuote)
	http.ListenAndServe(":8080", nil)
}

func handleQuote(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", QUOTE_API, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var coin map[string]Coin
	err = json.Unmarshal(body, &coin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coin["USDBRL"])
}
