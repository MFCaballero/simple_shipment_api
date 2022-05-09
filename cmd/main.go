package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MFCaballero/simple_shipment_api.git/api"
	"github.com/MFCaballero/simple_shipment_api.git/db"
	"golang.org/x/sync/errgroup"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	dataSourceName := "postgres://postgres:secret@localhost/postgres?sslmode=disable"
	port := ":9081"

	ctx, cancel := context.WithCancel(context.Background())

	eg, egCtx := errgroup.WithContext(context.Background())
	eg.Go(createHttpServer(port, dataSourceName, ctx, &wg))

	go func() {
		<-egCtx.Done()
		cancel()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}()

	if err := eg.Wait(); err != nil {
		fmt.Printf("error in the server goroutines: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("everything closed successfully")

}

func createHttpServer(port, dataSourceName string, ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		store, err := db.NewStore(dataSourceName)
		if err != nil {
			return err
		}

		handler := api.NewHandler(store)

		server := &http.Server{
			Addr:         port,
			WriteTimeout: time.Minute * 15,
			ReadTimeout:  time.Minute * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      handler,
		}
		errChan := make(chan error)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- fmt.Errorf("error shutting down the server: %w", err)
			}
			fmt.Println("server closed")
			close(errChan)
			wg.Done()
		}()
		fmt.Printf("Server up, running on port %s", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("error starting the server: %w", err)
		}
		fmt.Println("the server is closing")
		err = <-errChan
		wg.Wait()
		return err
	}
}
