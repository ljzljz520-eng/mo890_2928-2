package main

import (
	"flag"
	"log"
	"net/http"

	"community-memory-wall/internal/memorywall"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	store := memorywall.NewStore()
	server := &http.Server{Addr: *addr, Handler: memorywall.NewHTTPHandler(store)}
	log.Printf("community memory wall listening on %s", *addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
