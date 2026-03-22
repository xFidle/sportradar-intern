package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/xFidle/sportradar-intern/server/internal/setup"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	notifyCtx, notifyCancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer notifyCancel()

	<-notifyCtx.Done()

	log.Println("Shutting down gracefully, press Ctrl+C again to force.")
	notifyCancel() // Allow Ctrl+C to force shutdown

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	app := setup.NewApp()
	defer app.Close()

	done := make(chan bool, 1)

	go gracefulShutdown(app.Server, done)

	log.Printf("Server listening on: %s", app.Server.Addr)
	err := app.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
