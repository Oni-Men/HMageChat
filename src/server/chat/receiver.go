package chat

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"../../def"
)

// Receiver レシーバーの型。
type Receiver struct {
	UUID        string
	DisplayName string
	Connection  net.Conn
	Observer    chan<- Notification
}

// WaitMessage メッセージを待機する。一つ受け取ったら次を待機する
func (receiver Receiver) WaitMessage() {
	buf := make([]byte, 4*1024)

	for {
		len, error := receiver.Connection.Read(buf)
		if error != nil {
			receiver.Observer <- Notification{
				Context: &def.Context{
					Type:        def.Leave,
					UUID:        receiver.UUID,
					DisplayName: receiver.DisplayName,
					Timestamp:   time.Now().Format(time.RFC1123),
					Body:        "おやすー",
				},
			}
			return
		}

		context := new(def.Context)

		if err := json.Unmarshal(buf[:len], context); err != nil {
			fmt.Println("Couldn't unmarshal json.", err)
			continue
		}

		switch context.Type {
		case def.Join:
			receiver.UUID = context.UUID
			receiver.DisplayName = context.DisplayName
			receiver.Join(context)
		case def.Leave:
			receiver.Leave(context)
		case def.Message:
			receiver.Message(context)
		}
	}
}

// Join オブザーバーにJoinを通知する
func (receiver Receiver) Join(context *def.Context) {
	context.Type = def.Join
	receiver.Observer <- Notification{
		Connection: receiver.Connection,
		Context:    context,
	}
}

// Leave オブザーバーにLeaveを通知する
func (receiver Receiver) Leave(context *def.Context) {
	context.Type = def.Leave
	receiver.Observer <- Notification{
		Connection: receiver.Connection,
		Context:    context,
	}
}

// Message オブザーバーにMessageを通知する
func (receiver Receiver) Message(context *def.Context) {
	context.Type = def.Message
	receiver.Observer <- Notification{
		Context: context,
	}
}
