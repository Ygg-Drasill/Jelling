package main

import (
	"github.com/Ygg-Drasill/Jelling/backend/database"
	"github.com/Ygg-Drasill/Jelling/backend/handlers"
	"log"
	"net/http"
)

func main() {
	db := database.Open()
	err := database.Setup(db)
	if err != nil {
		log.Fatal(err)
	}
	ctx := handlers.NewContext(db)
	ctx.Logger.Info("database connection established")
	mux := handlers.NewJellingMux(ctx)
	address := "localhost:30420"
	ctx.Logger.Info("server started listening", "address", address)
	err = http.ListenAndServe("localhost:30420", mux)
	if err != nil {
		log.Fatal(err)
	}
}
