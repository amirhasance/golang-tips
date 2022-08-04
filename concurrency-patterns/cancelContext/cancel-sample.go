package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go startServer(ctx, &wg, 1)
	go startServer(ctx, &wg, 2)
	go startServer(ctx, &wg, 3)
	go longTimeRunningGoroutine(ctx, &wg)
	go longTimeRunningGoroutine(ctx, &wg)
	wg.Add(5)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancelFunc()
	wg.Wait()

}

func longTimeRunningGoroutine(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("shutting down Goroutines in ", time.Second)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		<-c.Done()
		fmt.Println("shut down goroutine ")
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// is doing sth
			time.Sleep(300 * time.Millisecond)
			fmt.Println("running default")

		}

	}
}

func startServer(ctx context.Context, wg *sync.WaitGroup, serverNumer int) {
	defer wg.Done()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Hello, World! %d ", serverNumer))
	})
	go func() {
		err := e.Start(fmt.Sprintf(":800%d", serverNumer))
		if err != nil {
			if err.Error() == "http: Server closed" {
				log.Error(err)
			} else {
				log.Fatal(err)
			}
		}
	}()

	// it waits for cancellation from parent Goroutine
	<-ctx.Done()
	log.Info("GraceFully shut down  server ", serverNumer)
	c, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(serverNumer))
	defer cancel()
	if err := e.Shutdown(c); err != nil {
		log.Fatal("server shutdown", err)
	}
	// this waits until context deadline is finished
	<-c.Done()
	log.Info(fmt.Sprintf("Server %d Gracefully shut down :)", serverNumer))

}
