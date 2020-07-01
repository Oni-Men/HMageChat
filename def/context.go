package def

//ContextType コンテキストの種類
type ContextType int

//コンテキストの種類
const (
	Join ContextType = iota
	Leave
	Message
)

//Context 通信するときのデータ。JSONで送受信し、パースして使用する
type Context struct {
	Type        ContextType
	UUID        string
	DisplayName string
	Timestamp   string
	Body        string
}
