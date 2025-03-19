package test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	upgrader := websocket.Upgrader{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 这个就是 websocket 的握手  搞升级的，或者说叫升级协议 初始化协议都行

		// conn 是一个 websocket 的连接
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// 升级WS协议失败
			w.Write([]byte(err.Error()))
			return
		}
		// 开启一个Goroutine 读取客户端发送的消息
		go func() {
			for {
				//  第一个参数是消息类型 （帧类型）文本帧 QpCode
				typ, message, err := conn.ReadMessage()
				switch typ {
				case websocket.TextMessage:
				case websocket.BinaryMessage:
				case websocket.CloseMessage:
					fmt.Println("close message")
					conn.Close() // 需要关闭WebSocket连接s
				case websocket.PingMessage:
				case websocket.PongMessage:
				}
				if err != nil {
					// 读取消息失败
					return
				}
				// 打印消息
				println(string(message))
			}
		}()

	})
}
