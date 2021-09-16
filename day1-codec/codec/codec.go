package codec

import "io"

type Header struct {
	ServiceMethod string
	Seq           int
	Error         string
}

type codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}
