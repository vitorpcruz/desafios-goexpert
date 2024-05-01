package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CoinBid struct {
	Bid string `json:"bid"`
}

func main() {
	getUSDBRL()
}

func getUSDBRL() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var coinBid CoinBid
	json.Unmarshal(body, &coinBid)

	byteBid, err := json.Marshal(coinBid.Bid)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	file.Write(byteBid)
	log.Println("Cotação salva.")
}
