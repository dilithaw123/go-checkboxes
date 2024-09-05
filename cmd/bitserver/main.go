package main

import (
	"fmt"
	"go-checkboxes/internal/bitset"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := os.Getenv("PORT")
	boxesStr := os.Getenv("BOXES")
	boxes, err := strconv.ParseUint(boxesStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	server := bitset.NewServer(boxes)
	mux := http.NewServeMux()
	server.ServeHTTP(mux)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatal(err)
	}
}
