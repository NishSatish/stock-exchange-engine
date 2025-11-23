package main

import (
	"fmt"
	"log"
	"runtime"
	"stock-exchange-simulator/api"
	"stock-exchange-simulator/pkg/app"
	"stock-exchange-simulator/pkg/db"
	"stock-exchange-simulator/pkg/libs"
	"stock-exchange-simulator/pkg/libs/taskqueue"
	"time"

	"github.com/hibiken/asynq"
)

func main() {
	fmt.Println("Hello, Stock Exchange Simulator!")
	conn, err := db.InitPostgres()
	if err != nil {
		log.Fatalf("Cant connect: %v", err)
	}
	defer conn.Close()

	// Init libs
	libsFactory := libs.NewLibsFactory()

	// Init app services
	services := app.NewAppServices(conn, libsFactory)

	// Start the asynq worker in the background
	go runTaskServer(services)

	log.Println("All services initialized successfully.")

	// Start the heap profiler in the background
	go logHeapStats()

	router := api.SetupRouter(services)
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

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

// runTaskServer starts the asynq worker server.
// decouple the run task server from libs to avoid circular dependencies.
func runTaskServer(app *app.AppServices) {
	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}
	taskServer := taskqueue.NewTaskServer(redisOpt)
	var mux = asynq.NewServeMux()

	// only shitty part about Asynq, you have to register all processors to event types in one place before you start the server
	mux = app.AsynqEventHandlerMap(mux)

	log.Println("Asynq worker server started...")
	err := taskServer.Server.Start(mux)
	if err != nil {
		panic("asynq faileedddd")
	}
}
