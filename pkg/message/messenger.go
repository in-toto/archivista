package message

type Message []byte

type Messenger interface {
	Construct(gitoid string, payload []byte) Message
}
