package def

import (
	"time"

	"github.com/google/uuid"
)

//ContextType コンテキストの種類
type ContextType string

//コンテキストの種類
const (
	Join    ContextType = "join"
	Leave   ContextType = "leave"
	Message ContextType = "message"
	Kick    ContextType = "kick"
	Close   ContextType = "close"
)

//Context 通信するときのデータ。JSONで送受信し、パースして使用する
type Context struct {
	Type        ContextType
	UUID        string
	DisplayName string
	Timestamp   string
	Body        string
}

// Normalize ゼロ値を、適切に埋める
func (ctx Context) Normalize() Context {
	if ctx.Type == "" {
		ctx.Type = Message
	}

	if ctx.UUID == "" {
		ctx.UUID = uuid.New().String()
	}

	if ctx.DisplayName == "" {
		ctx.DisplayName = ctx.UUID[:5]
	}

	if ctx.Timestamp == "" {
		ctx.Timestamp = time.Now().Format(time.RFC1123)
	}

	return ctx
}
