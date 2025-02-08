package message

import "io"

type Session interface {
	HandleWrite() func(w io.Writer) bool
	ReadMessage() []byte
}
