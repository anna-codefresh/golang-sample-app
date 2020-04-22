package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(CreateGreeting(name)))
}

func CreateGreeting(name string) string {
	if name == "" {
		name = "Guest"
	}
	return "Hello, " + name + "\n"
}

func main() {

	cli, err := client.NewEnvClient()
	if err != nil {
    		panic(err)
    }

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
    	if err != nil {
    		panic(err)
    	}

    for _, container := range containers {
    		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
    	}

	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
