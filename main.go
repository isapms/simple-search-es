package main

import (
	"log"
	"net/http"
	"os"
	"simple-search-es/cmd"
	"simple-search-es/internal/handlers/adverthdl"
	"simple-search-es/internal/repository/advertrepo"
	"simple-search-es/internal/services/advertsvc"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	elastic := cmd.ConfigEs()

	advertRepo := advertrepo.New(*elastic)
	advertSvc := advertsvc.New(advertRepo)
	advertHdl := adverthdl.New(advertSvc)

	router := mux.NewRouter()
	router.HandleFunc("/api/ads", advertHdl.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/ads/{id}", advertHdl.SearchOne).Methods(http.MethodGet)
	router.HandleFunc("/api/ads", advertHdl.Search).Methods(http.MethodGet)
	router.HandleFunc("/api/ads/{id}", advertHdl.Delete).Methods(http.MethodDelete)

	err := http.ListenAndServe(":"+os.Getenv("APP_SERVER_PORT"), router)
	if err != nil {
		log.Fatal(err)
	}
}
