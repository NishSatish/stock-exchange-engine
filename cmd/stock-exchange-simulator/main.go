package main

import (
	"fmt"
	"log"
	"stock-exchange-simulator/api"
	"stock-exchange-simulator/pkg/app"
	"stock-exchange-simulator/pkg/db"
)

func main() {
	fmt.Println("Hello, Stock Exchange Simulator!")
	conn, err := db.InitPostgres()
	if err != nil {
		log.Fatalf("Cant connect: %v", err)
	}
	defer conn.Close()

	services := app.NewAppServices(conn)

	log.Println("All services initialized successfully.")

	router := api.SetupRouter(services)
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
