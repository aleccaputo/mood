package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"mood/api"
	"mood/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listenAddr := flag.String("listen", ":8080", "Address to listen on")
	flag.Parse()

	postGresDatabase, err := db.NewPostgresDb()
	if err != nil {
		log.Fatal(err)
	}

	if err := postGresDatabase.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(*listenAddr, postGresDatabase)
	fmt.Println("server running on port: ", *listenAddr)
	log.Fatal(server.Start())

}
