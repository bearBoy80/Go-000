package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, _ := errgroup.WithContext(ctx)
	//start htttp server
	g.Go(func() error {
		//if start error return err
		if err := appService(ctx, ":8080", "app1", nil); err != nil {
			cancel()
			return err
		}
		return nil
	})
	g.Go(func() error {
		//if start error return err
		if err := listenSignal(ctx, cancel); err != nil {
			cancel()
			return err
		}
		return nil
	})
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully shutdow ")
	}
}
func appService(ctx context.Context, addr string, srvName string, handler http.HandlerFunc) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World, %s\n", srvName)
	})
	srv := &http.Server{
		Addr:           addr,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		<-ctx.Done()
		shutdown(ctx, srv, srvName)
	}()
	fmt.Println(srvName + " server start ......... ")
	return srv.ListenAndServe()
}
func listenSignal(ctx context.Context, cancel context.CancelFunc) error {
	sigint := make(chan os.Signal, 1)
	fmt.Println("listenSignal start ......... ")
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigint:
		// We received an SIGTERM \SIGINT signal, shut down.
		fmt.Printf("handle signal: %d \n", sigint)
		cancel()
		return nil
	case <-ctx.Done():
		// returning not to leak the goroutine
		return nil
	}
}

// shutdown app server
func shutdown(ctx context.Context, srv *http.Server, name string) {
	fmt.Println("start " + name + "Server stop.....")
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("srv shutdown fail: %+v", err)
		return
	}
	fmt.Println("end " + name + "Server stop.....")
}
