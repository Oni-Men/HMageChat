package chat

import (
	"encoding/json"
	"fmt"
	"net"

	"../../def"
)

//Sender 送信に必要な情報
type Sender struct {
	UUID       string
	Connection net.Conn
}

//SendMessage 接続先にJSON化されたコンテキストの情報を書き込む
func (sender Sender) SendMessage(context *def.Context) {
	buf, error := json.Marshal(context)

	if error != nil {
		fmt.Println("Couldn't marshal context")
	} else {
		_, error = sender.Connection.Write(buf)

		if error != nil {
			panic(error)
		}
	}
}
