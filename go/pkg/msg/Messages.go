package msg

import (
	"fmt"
	"net"
	"time"
)

type MessagePrinter struct {
	format string
}

func NewMessagePrinter() MessagePrinter {
	return MessagePrinter{
		format: time.RFC3339,
	}
}

func NewMessagePrinterWithFormat(format string) MessagePrinter {
	return MessagePrinter{
		format: format,
	}
}

func (m *MessagePrinter) PrintIn(conn net.Conn, message string, a ...interface{}) {
	msg := fmt.Sprintf(message, a...)
	fmt.Printf("%s %s <-- %s '%s'\n",
		time.Now().Format(m.format), conn.LocalAddr(), conn.RemoteAddr(), msg)
}

func (m *MessagePrinter) PrintOut(conn net.Conn, message string, a ...interface{}) {
	msg := fmt.Sprintf(message, a...)
	fmt.Printf("%s %s --> %s '%s'\n",
		time.Now().Format(m.format), conn.LocalAddr(), conn.RemoteAddr(), msg)
}

func (m *MessagePrinter) PrintInOut(conn net.Conn, message string, a ...interface{}) {
	msg := fmt.Sprintf(message, a...)
	fmt.Printf("%s %s <-> %s %s\n",
		time.Now().Format(m.format), conn.LocalAddr(), conn.RemoteAddr(), msg)
}
