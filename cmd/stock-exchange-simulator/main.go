package main

import (
	"fmt"
	"log"
	"runtime"
	"stock-exchange-simulator/api"
	"stock-exchange-simulator/pkg/app"
	"stock-exchange-simulator/pkg/db"
	"time"
)

// logHeapStats logs the current heap allocation every 10 seconds.
func logHeapStats() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		log.Printf("Heap Usage: %d KB", m.HeapAlloc/1024)
	}
}

func main() {
	fmt.Println("Hello, Stock Exchange Simulator!")
	conn, err := db.InitPostgres()
	if err != nil {
		log.Fatalf("Cant connect: %v", err)
	}
	defer conn.Close()

	services := app.NewAppServices(conn)

	log.Println("All services initialized successfully.")

	// Start the heap profiler in the background
	go logHeapStats()

	router := api.SetupRouter(services)
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
