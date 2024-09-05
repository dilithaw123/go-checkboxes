package web

import (
	"context"
	"encoding/binary"
	"errors"
	"go-checkboxes/internal/util"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/coder/websocket"
)

type WebServer struct {
	sync.Mutex
	Connections map[*websocket.Conn]struct{}
	internalRep []byte
	serverConn  *websocket.Conn
	templates   *template.Template
	boxes       uint64
	domain      string
}

func NewServer(sConn *websocket.Conn, boxes uint64, domain string) *WebServer {
	server := &WebServer{
		Connections: make(map[*websocket.Conn]struct{}),
		serverConn:  sConn,
		templates:   template.Must(template.ParseGlob("templates/*.html")),
		boxes:       boxes,
		domain:      domain,
	}
	go server.Listener()
	return server
}

func (s *WebServer) Listener() {
	for {
		_, bytes, err := s.serverConn.Read(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		s.internalRep = bytes
		s.Send(s.internalRep)
	}
}

func (s *WebServer) Add(conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()
	s.Connections[conn] = struct{}{}
}

func (s *WebServer) Remove(conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()
	delete(s.Connections, conn)
}

func (s *WebServer) Send(bytes []byte) {
	s.Lock()
	defer s.Unlock()
	for conn := range s.Connections {
		err := conn.Write(context.Background(), websocket.MessageBinary, bytes)
		if err != nil && errors.As(err, &websocket.CloseError{}) {
			delete(s.Connections, conn)
		}
	}
}

func (s *WebServer) Subscribe(w http.ResponseWriter, r *http.Request) {
	// Add domains here
	conn, err := websocket.Accept(
		w,
		r,
		&websocket.AcceptOptions{
			OriginPatterns:     []string{},
			InsecureSkipVerify: true,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	s.Add(conn)
	defer s.Remove(conn)
	ctx := r.Context()
	// try send initial state
	conn.Write(ctx, websocket.MessageBinary, s.internalRep)
	for {
		_, bytes, err := conn.Read(ctx)
		if err != nil && errors.As(err, &websocket.CloseError{}) {
			s.Remove(conn)
			return
		}
		if err != nil {
			return
		}
		str := string(bytes)
		log.Println("Received", str)
		if len(str) < 2 {
			return
		}
		pos, err := strconv.ParseUint(str[1:], 10, 64)
		if err != nil {
			return
		}
		bytes = util.EncodeSelection(pos, str[0] == '1')
		err = s.serverConn.Write(context.Background(), websocket.MessageBinary, bytes)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *WebServer) RouteHTTP(mux *http.ServeMux) {
	mux.HandleFunc("/subscribe", s.Subscribe)
	mux.HandleFunc("/setbit", func(w http.ResponseWriter, r *http.Request) {
		var pos int
		var err error
		if pos, err = strconv.Atoi(r.URL.Query().Get("pos")); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		bytes := make([]byte, 9)
		bytes[0] = 1
		binary.LittleEndian.PutUint64(bytes[1:], uint64(pos))
		err = s.serverConn.Write(context.Background(), websocket.MessageBinary, bytes)
		if err != nil {
			log.Fatal(err)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Kick rocks", http.StatusMethodNotAllowed)
		}
		type data struct {
			Num    int
			Domain string
		}
		err := s.templates.ExecuteTemplate(
			w,
			"index.html",
			data{
				Num:    int(s.boxes),
				Domain: s.domain,
			},
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})
}
