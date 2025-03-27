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
	mux := handlers.NewJellingMux(ctx)
	http.ListenAndServe("localhost:30420", mux)
}
