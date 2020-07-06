package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"../def"
)

// Receiver レシーバーの型。
type Receiver struct {
	UUID        string
	DisplayName string
	Connection  net.Conn
}

// WaitMessage メッセージを待機する。一つ受け取ったら次を待機する
func (receiver Receiver) WaitMessage() {
	buf := make([]byte, 4*1024)

	for {
		len, err := receiver.Connection.Read(buf)

		if err != nil && err != io.EOF {
			panic(err)
		}
		context := new(def.Context)
		if err := json.Unmarshal(buf[:len], context); err != nil {
			fmt.Println("Couldn't unmarshal json.", err)
			continue
		}

		time, err := time.Parse(time.RFC1123, context.Timestamp)
		formattedTime := time.Format("15:04:05")

		if err != nil {
			fmt.Println("Couldn't parse timestamp:", time)
			continue
		}

		switch context.Type {
		case def.Join:
			fallthrough
		case def.Leave:
			fmt.Printf("%s\n", context.Body)
		case def.Message:
			fmt.Printf("[%s][%s] %s\n", context.UUID, formattedTime, context.Body)
		case def.Kick:
			receiver.Connection.Close()
			return
		case def.Close:
			receiver.Connection.Close()
			return
		}
	}
}
