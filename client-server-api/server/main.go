package main

import (
	"net/http"

	"github.com/vitorpcruz/desafios-golang/client-server-api/server/db"
	"github.com/vitorpcruz/desafios-golang/client-server-api/server/repositories"
	"github.com/vitorpcruz/desafios-golang/client-server-api/server/services"
)

func main() {
	db := db.ConfigureDB()
	coinRepository := repositories.Init(db)
	coinService := services.Init(coinRepository)

	http.HandleFunc("/cotacao", coinService.HandleQuote)
	http.ListenAndServe(":8080", nil)
}
