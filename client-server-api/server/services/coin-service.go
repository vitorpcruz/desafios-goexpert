package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/vitorpcruz/desafios-golang/client-server-api/server/models"
	"github.com/vitorpcruz/desafios-golang/client-server-api/server/repositories"
)

type CoinServiceInterface interface {
	HandleQuote(w http.ResponseWriter, r *http.Request)
}

type CoinService struct {
	repo repositories.CoinRepositoryInterface
}

func Init(coinRepo repositories.CoinRepositoryInterface) *CoinService {
	return &CoinService{repo: coinRepo}
}

func (s *CoinService) HandleQuote(w http.ResponseWriter, r *http.Request) {
	const QUOTE_API = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

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
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var coins map[string]models.Coin
	err = json.Unmarshal(body, &coins)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	coin, contains := coins["USDBRL"]

	if !contains {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.repo.Save(&coin); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coin)
}

func (s *CoinService) Save(coin models.Coin) error {
	const timeout = time.Millisecond * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-time.After(time.Millisecond * 10):
		if err := s.Save(coin); err != nil {
			return err
		}
		return nil
	case <-ctx.Done():
		return errors.New("time exceeded")
	}
}
