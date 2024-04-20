package types

type MessageType string

const (
	MESSAGE_TYPE_TEXT   = MessageType("TEXT")
	MESSAGE_TYPE_OTHERS = MessageType("OTHERS")
)

type ParsedMessage struct {
	MessageType MessageType
	Text        string
	ID          string
	ReplyToken  string
}
