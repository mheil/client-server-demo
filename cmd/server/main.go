package main

import (
	"bufio"
	"client-server/pkg/errhandler"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	fmt.Println("server started")

	port := flag.Int("port", 7666, "port number to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	errhandler.HandleError(err)

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
			fmt.Printf("%s<->%s closing failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
		}
	}()

	fmt.Printf("%s<->%s opened\n", conn.LocalAddr(), conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Printf("%s<->%s closed\n", conn.LocalAddr(), conn.RemoteAddr())
			break
		}

		if err != nil {
			fmt.Printf("%s<->%s failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
			break
		}

		message = strings.TrimSpace(message)
		fmt.Printf("%s<->%s received: '%s'\n", conn.LocalAddr(), conn.RemoteAddr(), message)
		if "exit" == message {
			fmt.Printf("%s<->%s exit.. sending bye\n", conn.LocalAddr(), conn.RemoteAddr())
			_, err = fmt.Fprintf(conn, "bye...\n")
			if err != nil {
				fmt.Printf("%s<->%s sending bye failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
			}
			fmt.Printf("%s<->%s closing\n", conn.LocalAddr(), conn.RemoteAddr())
			errhandler.CloseWithPanic(conn)
			break
		}
	}
	fmt.Printf("%s<->%s finished\n", conn.LocalAddr(), conn.RemoteAddr())
}
