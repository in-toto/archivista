package message

import "github.com/in-toto/archivista/pkg/message"

type PredicateTypeMessenger struct{}

func (m PredicateTypeMessenger) Construct(gitoid string, payload []byte) message.Message {
	return message.Message{}
}
