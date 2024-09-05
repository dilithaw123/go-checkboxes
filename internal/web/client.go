package web

import (
	"context"

	"github.com/coder/websocket"
)

func Dial(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{})
	conn.SetReadLimit(-1)
	return conn, err
}
