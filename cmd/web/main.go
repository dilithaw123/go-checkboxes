package main

import (
	"fmt"
	"go-checkboxes/internal/web"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/coder/websocket"
)

func main() {
	port := os.Getenv("PORT")
	boxesStr := os.Getenv("BOXES")
	domain := os.Getenv("DOMAIN")
	boxes, err := strconv.ParseUint(boxesStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := web.Dial("ws://bitsetserver:5050/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "Bye")
	server := web.NewServer(conn, boxes, domain)
	mux := http.NewServeMux()
	server.RouteHTTP(mux)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatal(err)
	}
}
