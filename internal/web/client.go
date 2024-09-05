package web

import (
	"context"
	"encoding/binary"

	"github.com/coder/websocket"
)

func Dial(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{})
	conn.SetReadLimit(-1)
	return conn, err
}

func Send(conn *websocket.Conn, index int, value bool) error {
	bytes := make([]byte, 9)
	if value {
		bytes[0] = 1
	} else {
		bytes[0] = 0
	}
	binary.LittleEndian.AppendUint64(bytes[1:], uint64(index))
	return conn.Write(context.Background(), websocket.MessageBinary, bytes)
}

func Receive(conn *websocket.Conn) (index int, value bool, err error) {
	_, bytes, err := conn.Read(context.Background())
	if err != nil {
		return 0, false, err
	}
	if len(bytes) != 9 {
		return 0, false, nil
	}
	value = bytes[0] == 1
	index = int(binary.LittleEndian.Uint64(bytes[1:]))
	return index, value, nil
}
