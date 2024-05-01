package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
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
	coin, err := s.getUSDBRL()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		msgDefault := "Um erro ocorreu ao obter USDBRL: "

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CustomResponse(msgDefault + err.Error()))

		log.Printf(msgDefault, err.Error())
		return
	}

	if err := s.Save(coin); err != nil {
		msgDefault := "Não foi possível salvar a moeda USDBRL: "

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CustomResponse(msgDefault + err.Error()))

		log.Printf(msgDefault, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coin)
}

func (s *CoinService) getUSDBRL() (*models.Coin, error) {
	const QUOTE_API = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", QUOTE_API, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var coins map[string]models.Coin
	err = json.Unmarshal(body, &coins)

	if err != nil {
		return nil, err
	}

	coin, contains := coins["USDBRL"]

	if !contains {
		return nil, err
	}

	return &coin, nil
}

func (s *CoinService) Save(coin *models.Coin) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	select {
	case <-time.After(time.Millisecond * 10):
		if err := s.repo.Save(coin); err != nil {
			return err
		}
		return nil
	case <-ctx.Done():
		return errors.New("time exceeded")
	}
}
