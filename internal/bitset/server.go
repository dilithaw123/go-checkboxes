package bitset

import (
	"context"
	"errors"
	"go-checkboxes/internal/util"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

// websocket send

type BitSetServer struct {
	sync.Mutex
	Set         *BitSet
	Connections map[*websocket.Conn]struct{}
}

func (s *BitSetServer) Add(conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()
	s.Connections[conn] = struct{}{}
}

func (s *BitSetServer) Remove(conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()
	delete(s.Connections, conn)
}

func (s *BitSetServer) Send() {
	s.Lock()
	defer s.Unlock()
	bytes := s.Set.Bytes()
	for conn := range s.Connections {
		err := conn.Write(context.Background(), websocket.MessageBinary, bytes)
		if err != nil && errors.As(err, &websocket.CloseError{}) {
			delete(s.Connections, conn)
		}
	}
}

func (s *BitSetServer) Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	s.Add(conn)
	defer s.Remove(conn)
	ctx := r.Context()
	for {
		_, bytes, err := conn.Read(ctx)
		if err != nil && errors.As(err, &websocket.CloseError{}) {
			return
		}
		if err != nil {
			return
		}
		index, boolValue := util.DecodeSelection(bytes)
		if boolValue {
			s.Set.Set(index)
		} else {
			s.Set.Clear(index)
		}
		s.Send()
	}
}

func (s *BitSetServer) ServeHTTP(mux *http.ServeMux) {
	mux.HandleFunc("/", s.Subscribe)
}

func NewServer(size uint64) *BitSetServer {
	return &BitSetServer{
		Set:         NewBitSet(size),
		Connections: make(map[*websocket.Conn]struct{}),
	}
}
