package test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"testing"
	"time"
)

func Test_Demo(t *testing.T) {
	// WebSocket 服务器 URL
	url := "ws://localhost:8081/api/v1/ws"

	// 指定子协议
	subprotocol := "chat"

	// 连接到 WebSocket 服务器，传入子协议
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial failed:", err)
	}
	defer conn.Close()
	//map[string][]string{
	//		"sec-websocket-protocol": "auth-token, auth-eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkX2F0IjoxNzQzNDczMzExLCJ1c2VyX2lkIjowLCJ1c2VyX3V1aWQiOiIwMUpNSzU4U1Y4UE1SOUhIR0NTWDAzQkM3NCJ9.uw_IehF7BoDOmeWBfIBUNUYmS582HBYN0mFz3ubmkpZnS3UNUyxOuP1y3EUK1Gx5V3yji0RcrZKCBcBhgkKrWw",
	//	}
	// 验证是否支持子协议
	if subprotocol != "" && conn.Subprotocol() != subprotocol {
		fmt.Printf("Expected subprotocol %s, but got %s\n", subprotocol, conn.Subprotocol())
		os.Exit(1)
	}

	// 向服务器发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket server!"))
	if err != nil {
		log.Println("Write failed:", err)
		return
	}

	// 监听服务器的响应
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read failed:", err)
				return
			}
			fmt.Printf("Received: %s\n", msg)
		}
	}()

	// 等待一段时间，模拟客户端持续运行
	time.Sleep(30 * time.Second)
}
