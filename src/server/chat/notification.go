package chat

import (
	"net"

	"../../def"
)

//Notification チャネル同士でのやり取りに使用するデータ
type Notification struct {
	Connection net.Conn
	Context    *def.Context
}
