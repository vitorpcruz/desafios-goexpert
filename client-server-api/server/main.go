package main

import (
	"log"
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
	log.Println("Server running at 8080")

	http.ListenAndServe(":8080", nil)
}
