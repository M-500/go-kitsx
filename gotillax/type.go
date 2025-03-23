package gotillax

import (
	"github.com/gorilla/websocket"
	"sync"
)

type WsHub struct {
	connMap map[string]*websocket.Conn
	mu      sync.Mutex
}
