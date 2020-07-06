package chat

import (
	"fmt"
	"net"
	"time"

	"../../def"
)

//Observer 通知を受け取り、全ての接続に再配布する
type Observer struct {
	Senders []Sender
	Subject <-chan Notification
}

//WaitNotice 通知を待機する。一つ受け取ったら次を待機する
func (observer Observer) WaitNotice() {
	for {
		notice := <-observer.Subject

		ctx := notice.Context

		if ctx == nil {
			fmt.Printf("Context was nil\n")
			continue
		}

		normalized := ctx.Normalize()
		ctx = &normalized

		timestamp, error := time.Parse(time.RFC1123, ctx.Timestamp)

		if error != nil {
			fmt.Printf("Couldn't parse timestamp: %s\n", ctx.Timestamp)
			continue
		}

		formattedTime := timestamp.Format("15:04:05")

		switch ctx.Type {
		case def.Join:
			observer.Senders = appendSender(ctx.UUID, notice.Connection, observer.Senders)
			ctx.Body = fmt.Sprintf("%sが参加しました。", ctx.DisplayName)
		case def.Leave:
			observer.Senders = removeSender(ctx.UUID, observer.Senders)
			ctx.Body = fmt.Sprintf("%sが退出しました。", ctx.DisplayName)
		case def.Close:
			ctx.Body = "サーバーが閉鎖しました。"
		}
		fmt.Printf("[%s][%s] %s\n", ctx.DisplayName, formattedTime, ctx.Body)
		observer.SendAll(ctx)
	}
}

// SendAll  全ての接続に送信します
func (observer Observer) SendAll(context *def.Context) {
	for i := range observer.Senders {
		observer.Senders[i].SendMessage(context)
	}
}

func appendSender(id string, connection net.Conn, senders []Sender) []Sender {
	return append(senders, Sender{UUID: id, Connection: connection})
}

func removeSender(id string, senders []Sender) []Sender {
	find := -1
	for i := range senders {
		if senders[i].UUID == id {
			find = i
			break
		}
	}

	if find == -1 {
		return senders
	}

	return append(senders[:find], senders[find+1:]...)
}
