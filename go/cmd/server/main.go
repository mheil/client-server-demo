package main

import (
	"bufio"
	"client-server/pkg/errhandler"
	"client-server/pkg/msg"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

var mp = msg.NewMessagePrinter()

func main() {
	fmt.Println("server started")

	port := flag.Int("port", 7666, "port number to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	errhandler.PanicOnError(err)

	for {
		fmt.Printf("%s listening\n", ln.Addr())
		conn, err := ln.Accept()
		errhandler.PanicOnErrorWithMessage("listening for server failed", err)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("%s <-> %s closing failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
		}
	}()

	fmt.Printf("%s <-> %s opened\n", conn.LocalAddr(), conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			mp.PrintInOut(conn, "closed")
			break
		}

		if err != nil {
			mp.PrintInOut(conn, "failed %s", err)
			break
		}

		message = strings.TrimSpace(message)
		fmt.Printf("%s <-- %s '%s'\n", conn.LocalAddr(), conn.RemoteAddr(), message)
		if message == "exit" {
			mp.PrintOut(conn, "bye")
			_, err = fmt.Fprintln(conn, "bye")
			if err != nil {
				mp.PrintInOut(conn, "sending bye failed: %s", err)
			}
			mp.PrintInOut(conn, "closing")
			errhandler.CloseWithPanicOnError(conn)
			break
		}
		if message == "time" {
			timeMsg := time.Now().Format(time.RFC3339)
			mp.PrintOut(conn, timeMsg)
			_, err = fmt.Fprintln(conn, timeMsg)
			if err != nil {
				mp.PrintInOut(conn, "sending time failed: %s", err)
				mp.PrintInOut(conn, "closing")
				errhandler.CloseWithPanicOnError(conn)
				break
			}
		}
	}
	mp.PrintOut(conn, "finished")
}
