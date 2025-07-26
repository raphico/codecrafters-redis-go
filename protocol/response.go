package protocol

import (
	"fmt"
	"net"
)

type Response struct {
	conn net.Conn
}

func NewResponse(conn net.Conn) *Response {
	return &Response{conn}
}

func (r *Response) SendSimpleString(reply string) {
	fmt.Fprintf(r.conn, "+%s\r\n", reply)
}

func (r *Response) SendError(reply string) {
	fmt.Fprintf(r.conn, "-ERR %s\r\n", reply)
}

func (r *Response) SendNullBulkString() {
	fmt.Fprint(r.conn, "$-1\r\n")
}

func (r *Response) SendBulkString(reply string) {
	fmt.Fprintf(r.conn, "$%d\r\n%s\r\n", len(reply), reply)
}

func (r *Response) SendInteger(reply int) {
	fmt.Fprintf(r.conn, ":%d\r\n", reply)
}
